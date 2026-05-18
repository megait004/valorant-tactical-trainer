package localstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
)

type SettingsStore struct {
	path string
}

func NewSettingsStore() (*SettingsStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return &SettingsStore{path: filepath.Join(configDir, "Valorant Tactical Trainer", "settings.json")}, nil
}

func NewSettingsStoreAt(path string) *SettingsStore {
	return &SettingsStore{path: path}
}

func (s *SettingsStore) LoadDataSettings() (datasettings.DataSettings, error) {
	value := datasettings.DefaultDataSettings()

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return value, nil
	}
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return datasettings.DefaultDataSettings(), err
	}

	return normalize(value), nil
}

func (s *SettingsStore) SaveDataSettings(value datasettings.DataSettings) (datasettings.DataSettings, error) {
	next := normalize(value)
	next.LastUpdatedAt = time.Now().Format(time.RFC3339)

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return next, err
	}

	data, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return next, err
	}

	return next, os.WriteFile(s.path, data, 0o600)
}

func normalize(value datasettings.DataSettings) datasettings.DataSettings {
	value.RiotName = strings.TrimSpace(value.RiotName)
	value.RiotTag = strings.TrimPrefix(strings.TrimSpace(value.RiotTag), "#")
	value.PUUID = strings.TrimSpace(value.PUUID)
	value.APIKey = strings.TrimSpace(value.APIKey)
	value.Region = strings.ToLower(strings.TrimSpace(value.Region))
	value.Shard = strings.ToLower(strings.TrimSpace(value.Shard))
	value.APIKeyHeader = strings.TrimSpace(value.APIKeyHeader)
	value.RateLimitTier = strings.ToLower(strings.TrimSpace(value.RateLimitTier))

	// RGAPI key (Riot) không thuộc settings này — Riot client đọc từ .env. Nếu lỡ
	// có RGAPI key trong settings.APIKey (do bug cũ) thì xóa, tránh gửi nhầm sang
	// Henrik và bị 401.
	if strings.HasPrefix(value.APIKey, "RGAPI-") {
		value.APIKey = ""
	}

	if value.Region == "" {
		value.Region = "ap"
	}
	// X-Riot-Token là header Riot, không phải Henrik. Settings.APIKeyHeader chỉ dùng
	// cho Henrik authenticated tier — nên reset về default Authorization để Henrik không
	// bị nhận key sai và trả 401.
	if value.APIKeyHeader != "X-API-Key" {
		value.APIKeyHeader = "Authorization"
	}
	if value.RateLimitTier != "enhanced" {
		value.RateLimitTier = "basic"
	}
	if value.MatchCount <= 0 {
		value.MatchCount = 5
	}
	if value.MatchCount > 20 {
		value.MatchCount = 20
	}
	if value.CacheTTLMinutes <= 0 {
		value.CacheTTLMinutes = 30
	}
	if value.CacheTTLMinutes > 240 {
		value.CacheTTLMinutes = 240
	}

	return value
}
