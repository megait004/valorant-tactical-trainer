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

// Adapter Anthropic Messages:
//
//	POST https://api.anthropic.com/v1/messages
//	Header: x-api-key, anthropic-version: 2023-06-01
//	Body:   { model, max_tokens, system, messages: [{role:"user|assistant", content}] }

type anthropicReq struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResp struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Coach) callAnthropic(ctx context.Context, system string, messages []ChatMessage) (string, error) {
	msgs := make([]anthropicMessage, 0, len(messages))
	for _, m := range messages {
		role := m.Role
		if role != "assistant" {
			role = "user"
		}
		msgs = append(msgs, anthropicMessage{Role: role, Content: m.Content})
	}
	body, err := json.Marshal(anthropicReq{
		Model:     c.model,
		MaxTokens: 1500,
		System:    system,
		Messages:  msgs,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("err call Anthropic: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Anthropic status %d: %s", resp.StatusCode, truncate(string(raw), 300))
	}
	var parsed anthropicResp
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("err parse Anthropic response: %w", err)
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("Anthropic error: %s", parsed.Error.Message)
	}
	for _, block := range parsed.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}
	return "", errors.New("Anthropic không trả text block")
}
