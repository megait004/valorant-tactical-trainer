package wailsiface

import (
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
	"valorant-tactical-trainer/desktop/internal/infrastructure/henrik"
	"valorant-tactical-trainer/desktop/internal/infrastructure/store"
)

// SettingsService expose Data & Consent settings + Henrik API status cho UI.
type SettingsService struct {
	store  *store.SettingsStore
	client *henrik.Client
}

func (s *SettingsService) GetDataSettings() (datasettings.DataSettings, error) {
	return s.store.LoadDataSettings()
}

func (s *SettingsService) SaveDataSettings(value datasettings.DataSettings) (datasettings.DataSettings, error) {
	return s.store.SaveDataSettings(value)
}

func (s *SettingsService) GetAPIStatus() (henrik.Status, error) {
	value, err := s.store.LoadDataSettings()
	if err != nil {
		return henrik.Status{}, err
	}
	return s.client.Status(value), nil
}
