package store

import (
	"path/filepath"
	"testing"

	"valorant-tactical-trainer/desktop/internal/domain/practice"
)

func TestPracticeSessionStoreAddsSession(t *testing.T) {
	store := NewPracticeSessionStoreAt(filepath.Join(t.TempDir(), "sessions.json"))
	state, err := store.Add(practice.SessionInput{TaskID: "t1", Focus: "Aim", DurationSeconds: 42})
	if err != nil {
		t.Fatalf("Add err: %v", err)
	}
	if len(state.Sessions) != 1 || state.Sessions[0].TaskID != "t1" || state.Sessions[0].DurationSeconds != 42 {
		t.Fatalf("unexpected state: %#v", state)
	}
}
