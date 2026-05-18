package localstore

import (
	"path/filepath"
	"testing"
)

func TestPracticeProgressStorePersistsState(t *testing.T) {
	store := NewPracticeProgressStoreAt(filepath.Join(t.TempDir(), "practice_progress.json"))

	initial, err := store.Load()
	if err != nil {
		t.Fatalf("expected missing progress to be safe, got %v", err)
	}
	if len(initial.Items) != 0 {
		t.Fatalf("expected empty progress, got %+v", initial.Items)
	}

	saved, err := store.Set("day-1-map-gap-0", true)
	if err != nil {
		t.Fatalf("expected set progress, got %v", err)
	}
	if !saved.Items["day-1-map-gap-0"] || saved.UpdatedAt == "" {
		t.Fatalf("unexpected saved progress: %+v", saved)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("expected load progress, got %v", err)
	}
	if !loaded.Items["day-1-map-gap-0"] {
		t.Fatalf("expected persisted item, got %+v", loaded.Items)
	}
}

func TestPracticeProgressStoreRejectsEmptyID(t *testing.T) {
	store := NewPracticeProgressStoreAt(filepath.Join(t.TempDir(), "practice_progress.json"))
	if _, err := store.Set(" ", true); err == nil {
		t.Fatal("expected err")
	}
}

func TestPracticeProgressStoreReset(t *testing.T) {
	store := NewPracticeProgressStoreAt(filepath.Join(t.TempDir(), "practice_progress.json"))
	if _, err := store.Set("day-1", true); err != nil {
		t.Fatalf("expected set progress, got %v", err)
	}
	reset, err := store.Reset()
	if err != nil {
		t.Fatalf("expected reset progress, got %v", err)
	}
	if len(reset.Items) != 0 || reset.UpdatedAt == "" {
		t.Fatalf("unexpected reset progress: %+v", reset)
	}
}
