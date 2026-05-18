package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Adapter Google Gemini (AI Studio):
//
//	POST https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent?key={apiKey}
//	Body: { system_instruction?, contents: [{ role:"user|model", parts:[{text}] }],
//	        generationConfig: { temperature, responseMimeType? } }
//
// Free tier rất chặt — có retry 1 lần với delay 2s khi gặp 429.

type geminiReq struct {
	SystemInstruction *geminiContent          `json:"system_instruction,omitempty"`
	Contents          []geminiContent         `json:"contents"`
	GenerationConfig  *geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	ResponseMIMEType string  `json:"responseMimeType,omitempty"`
}

type geminiResp struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (c *Coach) callGemini(ctx context.Context, system string, messages []ChatMessage, jsonMode bool) (string, error) {
	contents := make([]geminiContent, 0, len(messages))
	for _, m := range messages {
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, geminiContent{Role: role, Parts: []geminiPart{{Text: m.Content}}})
	}
	payload := geminiReq{
		Contents:         contents,
		GenerationConfig: &geminiGenerationConfig{Temperature: 0.4},
	}
	if system != "" {
		payload.SystemInstruction = &geminiContent{Parts: []geminiPart{{Text: system}}}
	}
	if jsonMode {
		payload.GenerationConfig.ResponseMIMEType = "application/json"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	const maxAttempts = 2
	var lastStatus int
	var lastBody string
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		text, status, raw, err := c.geminiOnce(ctx, body)
		if err == nil {
			return text, nil
		}
		lastStatus = status
		lastBody = raw
		if status != http.StatusTooManyRequests || attempt == maxAttempts {
			return "", err
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return "", fmt.Errorf("Gemini status %d: %s", lastStatus, truncate(lastBody, 300))
}

func (c *Coach) geminiOnce(ctx context.Context, body []byte) (text string, status int, raw string, err error) {
	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, "", fmt.Errorf("err call Gemini: %w", err)
	}
	defer resp.Body.Close()
	rawBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	raw = string(rawBytes)
	status = resp.StatusCode

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", status, raw, fmt.Errorf("Gemini 429 — hết quota free tier (model %s). Đợi 1 phút hoặc đổi LLM_MODEL trong .env (ví dụ: gemini-2.0-flash-lite)", c.model)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", status, raw, fmt.Errorf("Gemini status %d: %s", resp.StatusCode, truncate(raw, 300))
	}

	var parsed geminiResp
	if err := json.Unmarshal(rawBytes, &parsed); err != nil {
		return "", status, raw, fmt.Errorf("err parse Gemini response: %w", err)
	}
	if parsed.Error != nil {
		return "", status, raw, fmt.Errorf("Gemini error %s: %s", parsed.Error.Status, parsed.Error.Message)
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", status, raw, errors.New("Gemini không trả candidate")
	}
	return parsed.Candidates[0].Content.Parts[0].Text, status, raw, nil
}
