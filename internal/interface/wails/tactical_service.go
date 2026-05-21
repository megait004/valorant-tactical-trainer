package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/tactical"
	"valorant-tactical-trainer/desktop/internal/infrastructure/store"
)

// TacticalService expose Map Catalog + Map Plan (markers, lines, notes) cho UI.
type TacticalService struct {
	store *store.TacticalPlanStore
}

func (s *TacticalService) ListMaps() []tactical.MapCatalogEntry {
	return tactical.AllMaps()
}

func (s *TacticalService) LoadMapPlan(mapID string) (tactical.MapPlan, error) {
	plan, _, err := s.store.LoadPlan(mapID)
	return plan, err
}

func (s *TacticalService) SaveMapPlan(plan tactical.MapPlan) (tactical.MapPlan, error) {
	return s.store.SavePlan(plan)
}

func (s *TacticalService) DeleteMapPlan(mapID string) error {
	return s.store.DeletePlan(mapID)
}
