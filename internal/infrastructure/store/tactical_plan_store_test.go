package store

import (
	"path/filepath"
	"testing"

	"valorant-tactical-trainer/desktop/internal/domain/tactical"
)

func TestTacticalPlanStoreSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tactical_plans.json")
	store := NewTacticalPlanStoreAt(path)

	plan := tactical.MapPlan{
		MapID: "ascent",
		Title: "A execute",
		Side:  "attack",
		Markers: []tactical.PlanMarker{
			{ID: "m1", Kind: "smoke", Label: "Main smoke", X: 42, Y: 55},
		},
	}

	saved, err := store.SavePlan(plan)
	if err != nil {
		t.Fatal(err)
	}
	if saved.MapID != "ascent" {
		t.Fatalf("expected ascent, got %s", saved.MapID)
	}

	loaded, ok, err := store.LoadPlan("ascent")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected plan to exist")
	}
	if len(loaded.Markers) != 1 || loaded.Markers[0].Label != "Main smoke" {
		t.Fatalf("unexpected markers %#v", loaded.Markers)
	}
}
