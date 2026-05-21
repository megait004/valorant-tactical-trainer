package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/practice"
)

type PracticeSessionStore struct {
	path string
}

func NewPracticeSessionStore() (*PracticeSessionStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &PracticeSessionStore{path: filepath.Join(configDir, "Valorant Tactical Trainer", "practice_sessions.json")}, nil
}

func NewPracticeSessionStoreAt(path string) *PracticeSessionStore {
	return &PracticeSessionStore{path: path}
}

func (s *PracticeSessionStore) Load() (practice.SessionState, error) {
	state := practice.SessionState{Sessions: []practice.Session{}}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return practice.SessionState{Sessions: []practice.Session{}}, err
	}
	if state.Sessions == nil {
		state.Sessions = []practice.Session{}
	}
	return state, nil
}

func (s *PracticeSessionStore) Add(input practice.SessionInput) (practice.SessionState, error) {
	state, err := s.Load()
	if err != nil {
		return state, err
	}
	state.Sessions = append([]practice.Session{practice.NewSession(input)}, state.Sessions...)
	return s.Save(state)
}

func (s *PracticeSessionStore) Save(state practice.SessionState) (practice.SessionState, error) {
	if state.Sessions == nil {
		state.Sessions = []practice.Session{}
	}
	if len(state.Sessions) > 100 {
		state.Sessions = state.Sessions[:100]
	}
	if state.UpdatedAt == "" {
		state.UpdatedAt = time.Now().Format(time.RFC3339)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return state, err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return state, err
	}
	return state, os.WriteFile(s.path, data, 0o600)
}
