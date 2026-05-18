package localstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"valorant-tactical-trainer/desktop/internal/domain/tactical"
)

type TacticalPlanStore struct {
	path string
}

func NewTacticalPlanStore() (*TacticalPlanStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return &TacticalPlanStore{path: filepath.Join(configDir, "Valorant Tactical Trainer", "tactical_plans.json")}, nil
}

func NewTacticalPlanStoreAt(path string) *TacticalPlanStore {
	return &TacticalPlanStore{path: path}
}

type tacticalPlansFile struct {
	Plans map[string]tactical.MapPlan `json:"plans"`
}

func (s *TacticalPlanStore) LoadPlan(mapID string) (tactical.MapPlan, bool, error) {
	file, err := s.readFile()
	if errors.Is(err, os.ErrNotExist) {
		return tactical.DefaultPlan(mapID), false, nil
	}
	if err != nil {
		return tactical.DefaultPlan(mapID), false, err
	}
	plan, ok := file.Plans[mapID]
	if !ok {
		return tactical.DefaultPlan(mapID), false, nil
	}
	return plan, true, nil
}

func (s *TacticalPlanStore) SavePlan(plan tactical.MapPlan) (tactical.MapPlan, error) {
	next := tactical.NormalizePlan(plan)
	file, err := s.readFile()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return next, err
	}
	if file.Plans == nil {
		file.Plans = map[string]tactical.MapPlan{}
	}
	file.Plans[next.MapID] = next
	return next, s.writeFile(file)
}

func (s *TacticalPlanStore) DeletePlan(mapID string) error {
	file, err := s.readFile()
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	delete(file.Plans, mapID)
	return s.writeFile(file)
}

func (s *TacticalPlanStore) readFile() (tacticalPlansFile, error) {
	value := tacticalPlansFile{Plans: map[string]tactical.MapPlan{}}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return tacticalPlansFile{Plans: map[string]tactical.MapPlan{}}, err
	}
	if value.Plans == nil {
		value.Plans = map[string]tactical.MapPlan{}
	}
	return value, nil
}

func (s *TacticalPlanStore) writeFile(file tacticalPlansFile) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
