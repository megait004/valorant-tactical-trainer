package practice

import (
	"strings"
	"time"
)

type Session struct {
	ID              string `json:"id"`
	TaskID          string `json:"taskId"`
	Focus           string `json:"focus"`
	Map             string `json:"map"`
	Agent           string `json:"agent"`
	DurationSeconds int    `json:"durationSeconds"`
	StartedAt       string `json:"startedAt"`
	FinishedAt      string `json:"finishedAt"`
}

type SessionInput struct {
	TaskID          string `json:"taskId"`
	Focus           string `json:"focus"`
	Map             string `json:"map"`
	Agent           string `json:"agent"`
	DurationSeconds int    `json:"durationSeconds"`
	StartedAt       string `json:"startedAt"`
}

type SessionState struct {
	Sessions  []Session `json:"sessions"`
	UpdatedAt string    `json:"updatedAt"`
}

func NewSession(input SessionInput) Session {
	now := time.Now().Format(time.RFC3339)
	startedAt := strings.TrimSpace(input.StartedAt)
	if startedAt == "" {
		startedAt = now
	}
	if input.DurationSeconds < 0 {
		input.DurationSeconds = 0
	}
	return Session{
		ID:              "session-" + time.Now().Format("20060102150405.000000000"),
		TaskID:          strings.TrimSpace(input.TaskID),
		Focus:           strings.TrimSpace(input.Focus),
		Map:             strings.TrimSpace(input.Map),
		Agent:           strings.TrimSpace(input.Agent),
		DurationSeconds: input.DurationSeconds,
		StartedAt:       startedAt,
		FinishedAt:      now,
	}
}
