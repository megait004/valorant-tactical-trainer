package wailsiface

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"valorant-tactical-trainer/internal/infrastructure/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SettingsService struct {
	ctx   context.Context
	store *storage.Store
}

func NewSettingsService(store *storage.Store) *SettingsService {
	return &SettingsService{store: store}
}

func (service *SettingsService) Startup(ctx context.Context) {
	service.ctx = ctx
}

type ResetResult struct {
	Message string `json:"message"`
}

type SettingsDTO struct {
	APIKeyConfigured    bool   `json:"apiKeyConfigured"`
	Language            string `json:"language"`
	DataPath            string `json:"dataPath"`
	CacheEntries        int    `json:"cacheEntries"`
	ExpiredCacheEntries int    `json:"expiredCacheEntries"`
	Players             int    `json:"players"`
	Matches             int    `json:"matches"`
	RankSnapshots       int    `json:"rankSnapshots"`
	Reports             int    `json:"reports"`
	Message             string `json:"message"`
}

type SaveSettingsInput struct {
	APIKey string `json:"apiKey"`
}

type SaveLanguageInput struct {
	Language string `json:"language"`
}

type ClearCacheResult struct {
	Cleared int    `json:"cleared"`
	Message string `json:"message"`
}

type ExportDataResult struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (service *SettingsService) GetSettings() (SettingsDTO, error) {
	return service.settingsDTO("settings loaded")
}

func (service *SettingsService) SaveSettings(input SaveSettingsInput) (SettingsDTO, error) {
	apiKey := strings.TrimSpace(input.APIKey)
	if apiKey == "" {
		if err := service.store.DeleteSetting(context.Background(), "api_key"); err != nil {
			return SettingsDTO{}, err
		}
		return service.settingsDTO("api key cleared")
	}

	if err := service.store.SaveSetting(context.Background(), "api_key", apiKey); err != nil {
		return SettingsDTO{}, err
	}

	return service.settingsDTO("api key saved locally")
}

func (service *SettingsService) SaveLanguage(input SaveLanguageInput) (SettingsDTO, error) {
	language := normalizeLanguage(input.Language)
	if err := service.store.SaveSetting(context.Background(), "language", language); err != nil {
		return SettingsDTO{}, err
	}

	return service.settingsDTO("language saved")
}

func (service *SettingsService) ClearExpiredCache() (ClearCacheResult, error) {
	cleared, err := service.store.ClearExpiredAPICache(context.Background())
	if err != nil {
		return ClearCacheResult{}, err
	}

	return ClearCacheResult{Cleared: cleared, Message: "expired cache cleared"}, nil
}

func (service *SettingsService) ExportLocalData() (ExportDataResult, error) {
	if service.ctx == nil {
		return ExportDataResult{}, fmt.Errorf("ExportUnavailable: desktop context not ready")
	}

	path, err := runtime.SaveFileDialog(service.ctx, runtime.SaveDialogOptions{
		Title:           "Export Valorant Tactical Trainer data",
		DefaultFilename: fmt.Sprintf("valorant-tactical-trainer-export-%s.json", time.Now().Format("20060102-150405")),
		Filters: []runtime.FileFilter{{
			DisplayName: "JSON files (*.json)",
			Pattern:     "*.json",
		}},
		CanCreateDirectories: true,
	})
	if err != nil {
		return ExportDataResult{}, fmt.Errorf("err showing export dialog: %w", err)
	}
	if path == "" {
		return ExportDataResult{Message: "export cancelled"}, nil
	}

	data, err := service.store.ExportJSON(context.Background())
	if err != nil {
		return ExportDataResult{}, err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return ExportDataResult{}, fmt.Errorf("write export file: %w", err)
	}

	return ExportDataResult{Path: path, Message: "local data exported"}, nil
}

func (service *SettingsService) ResetAllData() (ResetResult, error) {
	if service.ctx != nil {
		choice, err := runtime.MessageDialog(service.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Reset local data",
			Message:       "Reset all local Valorant Tactical Trainer data on this machine? This deletes player, consent, cache, matches, reports, and recommendations.",
			Buttons:       []string{"Reset", "Cancel"},
			DefaultButton: "Cancel",
			CancelButton:  "Cancel",
		})
		if err != nil {
			return ResetResult{}, fmt.Errorf("err showing reset dialog: %w", err)
		}
		if choice != "Reset" {
			return ResetResult{Message: "reset cancelled"}, nil
		}
	}

	if err := service.store.ResetAll(context.Background()); err != nil {
		return ResetResult{}, fmt.Errorf("err reset data: %w", err)
	}

	return ResetResult{Message: "local data reset"}, nil
}

func (service *SettingsService) settingsDTO(message string) (SettingsDTO, error) {
	ctx := context.Background()
	_, hasAPIKey, err := service.store.Setting(ctx, "api_key")
	if err != nil {
		return SettingsDTO{}, err
	}
	language, hasLanguage, err := service.store.Setting(ctx, "language")
	if err != nil {
		return SettingsDTO{}, err
	}
	if !hasLanguage {
		language = "en"
	}
	language = normalizeLanguage(language)
	stats, err := service.store.Stats(ctx)
	if err != nil {
		return SettingsDTO{}, err
	}

	return SettingsDTO{
		APIKeyConfigured:    hasAPIKey,
		Language:            language,
		DataPath:            service.store.Path(),
		CacheEntries:        stats.CacheEntries,
		ExpiredCacheEntries: stats.ExpiredCacheEntries,
		Players:             stats.Players,
		Matches:             stats.Matches,
		RankSnapshots:       stats.RankSnapshots,
		Reports:             stats.Reports,
		Message:             message,
	}, nil
}

func normalizeLanguage(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "vi":
		return "vi"
	default:
		return "en"
	}
}
