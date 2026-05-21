package store

import (
	"path/filepath"
	"testing"

	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
)

func TestSettingsStorePersistsNormalizedDataSettings(t *testing.T) {
	store := NewSettingsStoreAt(filepath.Join(t.TempDir(), "settings.json"))

	saved, err := store.SaveDataSettings(datasettings.DataSettings{
		ConsentPersonalData: true,
		RiotName:            " giaphue ",
		RiotTag:             " #DATN ",
		APIKey:              " key-123 ",
		APIKeyHeader:        "X-API-Key",
		Region:              "AP",
		RateLimitTier:       "enhanced",
		MatchCount:          25,
		CacheTTLMinutes:     999,
	})
	if err != nil {
		t.Fatalf("expected save settings, got %v", err)
	}
	if saved.CacheTTLMinutes != 240 {
		t.Fatalf("expected capped cache ttl, got %d", saved.CacheTTLMinutes)
	}
	if saved.MatchCount != 20 {
		t.Fatalf("expected capped match count, got %d", saved.MatchCount)
	}

	loaded, err := store.LoadDataSettings()
	if err != nil {
		t.Fatalf("expected load settings, got %v", err)
	}
	if loaded.RiotName != "giaphue" || loaded.RiotTag != "DATN" || loaded.Region != "ap" {
		t.Fatalf("unexpected normalized settings: %+v", loaded)
	}
	if loaded.APIKey != "key-123" || loaded.APIKeyHeader != "X-API-Key" {
		t.Fatalf("unexpected api key settings: %+v", loaded)
	}
	if loaded.MatchCount != 20 {
		t.Fatalf("unexpected match count: %+v", loaded)
	}
}
