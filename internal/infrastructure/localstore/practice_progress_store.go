package localstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PracticeProgressState struct {
	Items     map[string]bool `json:"items"`
	UpdatedAt string          `json:"updatedAt"`
}

type PracticeProgressStore struct {
	path string
}

func NewPracticeProgressStore() (*PracticeProgressStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return &PracticeProgressStore{path: filepath.Join(configDir, "Valorant Tactical Trainer", "practice_progress.json")}, nil
}

func NewPracticeProgressStoreAt(path string) *PracticeProgressStore {
	return &PracticeProgressStore{path: path}
}

func (s *PracticeProgressStore) Load() (PracticeProgressState, error) {
	state := PracticeProgressState{Items: map[string]bool{}}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return PracticeProgressState{Items: map[string]bool{}}, err
	}
	if state.Items == nil {
		state.Items = map[string]bool{}
	}
	return state, nil
}

func (s *PracticeProgressStore) Set(itemID string, done bool) (PracticeProgressState, error) {
	state, err := s.Load()
	if err != nil {
		return state, err
	}
	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return state, errors.New("practice progress item id trống")
	}
	state.Items[itemID] = done
	state.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return state, err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return state, err
	}
	return state, os.WriteFile(s.path, data, 0o600)
}

func (s *PracticeProgressStore) Reset() (PracticeProgressState, error) {
	state := PracticeProgressState{Items: map[string]bool{}, UpdatedAt: time.Now().Format(time.RFC3339)}
	return s.Save(state)
}

func (s *PracticeProgressStore) Save(state PracticeProgressState) (PracticeProgressState, error) {
	if state.Items == nil {
		state.Items = map[string]bool{}
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
