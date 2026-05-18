package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Adapter OpenAI Chat Completions:
//
//	POST https://api.openai.com/v1/chat/completions
//	Header: Authorization: Bearer <api-key>
//	Body:   { model, messages: [{role:"system|user|assistant", content}], response_format? }

type openaiReq struct {
	Model          string                `json:"model"`
	Temperature    float64               `json:"temperature"`
	Messages       []openaiMessage       `json:"messages"`
	ResponseFormat *openaiResponseFormat `json:"response_format,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponseFormat struct {
	Type string `json:"type"`
}

type openaiResp struct {
	Choices []struct {
		Message openaiMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Coach) callOpenAI(ctx context.Context, system string, messages []ChatMessage, jsonMode bool) (string, error) {
	openaiMessages := []openaiMessage{{Role: "system", Content: system}}
	for _, m := range messages {
		role := m.Role
		if role != "assistant" {
			role = "user"
		}
		openaiMessages = append(openaiMessages, openaiMessage{Role: role, Content: m.Content})
	}
	payload := openaiReq{
		Model:       c.model,
		Temperature: 0.4,
		Messages:    openaiMessages,
	}
	if jsonMode {
		payload.ResponseFormat = &openaiResponseFormat{Type: "json_object"}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("err call OpenAI: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("OpenAI status %d: %s", resp.StatusCode, truncate(string(raw), 300))
	}
	var parsed openaiResp
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("err parse OpenAI response: %w", err)
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("OpenAI error: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("OpenAI không trả choices")
	}
	return parsed.Choices[0].Message.Content, nil
}
