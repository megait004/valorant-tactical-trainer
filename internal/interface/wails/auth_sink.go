package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/infrastructure/localstore"
	"valorant-tactical-trainer/desktop/internal/infrastructure/riot"
)

// authSink gắn riot.Client.Login/Logout với SettingsStore + ReportStore.
//
// Mục đích:
//   - Khi login: cập nhật RiotID + PUUID vào settings.json để
//     AnalysisService.FetchLiveReport biết đúng nick.
//   - Khi đổi nick (PUUID khác): xoá last_report.json để tránh stale cache
//     của nick cũ.
//   - Khi logout: clear PUUID/RiotID + xoá report cache.
type authSink struct {
	settings    *localstore.SettingsStore
	reportStore *localstore.ReportStore
}

func (s *authSink) OnLogin(player riot.PlayerInfo) error {
	current, err := s.settings.LoadDataSettings()
	if err != nil {
		return err
	}
	previousPUUID := current.PUUID

	current.ConsentPersonalData = true
	current.RiotName = player.GameName
	current.RiotTag = player.TagLine
	current.PUUID = player.PUUID
	current.Region = player.Region
	current.Shard = player.Shard
	// KHÔNG ép APIKeyHeader = X-Riot-Token: header đó dành riêng cho Riot client
	// (đã hard-code trong infrastructure/riot). Settings.APIKey/APIKeyHeader phục
	// vụ Henrik authenticated tier (HDEV- key). Nếu user không nhập key Henrik thì
	// Henrik vẫn fetch được public.
	if _, err := s.settings.SaveDataSettings(current); err != nil {
		return err
	}

	if previousPUUID != "" && previousPUUID != player.PUUID {
		_ = s.reportStore.Delete()
	}
	return nil
}

func (s *authSink) OnLogout() error {
	current, err := s.settings.LoadDataSettings()
	if err == nil {
		current.PUUID = ""
		current.RiotName = ""
		current.RiotTag = ""
		current.ConsentPersonalData = false
		_, _ = s.settings.SaveDataSettings(current)
	}
	_ = s.reportStore.Delete()
	return nil
}
