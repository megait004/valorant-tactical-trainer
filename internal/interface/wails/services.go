// Package wailsiface chứa các service được Wails Bind ra cho frontend
// (window.go.wailsiface.*). Mỗi file = 1 service để dễ tìm:
//
//	services.go          — Services struct + factory NewServices() (file này)
//	stores.go            — Helper khởi tạo store local có fallback path
//	auth_service.go      — Riot Account-V1 login/logout, lấy PUUID + region
//	auth_sink.go         — Bridge giữa riot.Client và store (đồng bộ session)
//	settings_service.go  — Data & Consent (RiotID, region, Henrik key)
//	analysis_service.go  — Fetch + analyse match history, sinh Report
//	practice_service.go  — Practice progress + lịch sử session luyện tập
//	assistant_service.go — Live Assistant in-game tip engine
//	tactical_service.go  — Map catalog + Map planner state
//	chat_service.go      — Bot AI chat (LLM coach hỏi-đáp về report)
package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	"valorant-tactical-trainer/desktop/internal/domain/assistant"
	"valorant-tactical-trainer/desktop/internal/infrastructure/henrik"
	"valorant-tactical-trainer/desktop/internal/infrastructure/llm"
	"valorant-tactical-trainer/desktop/internal/infrastructure/riot"
)

// Services là tập hợp tất cả Wails service. main.go chỉ Bind các field này.
type Services struct {
	Settings  *SettingsService
	Analysis  *AnalysisService
	Practice  *PracticeService
	Assistant *AssistantService
	Tactical  *TacticalService
	Auth      *AuthService
	Chat      *ChatService
}

// NewServices wire toàn bộ dependency theo thứ tự:
//
//  1. Khởi tạo các store local (settings, report, practice, tactical) — có
//     fallback về relative path khi UserConfigDir() lỗi (dev mode).
//  2. Khởi tạo infrastructure clients (Henrik, Riot Match, LLM Coach).
//  3. Wire authSink để riot.Client tự đồng bộ session vào store khi
//     Login/Logout.
//  4. Đóng gói tất cả vào Services để main.go Bind 1 lần duy nhất.
func NewServices() Services {
	settingsStore := mustSettingsStore()
	reportStore := mustReportStore()
	practiceProgress := mustPracticeProgressStore()
	practiceSessions := mustPracticeSessionStore()
	tacticalStore := mustTacticalPlanStore()

	henrikClient := henrik.NewClient()
	riotMatchClient := riot.NewMatchClient()

	// LLM coach optional — nil khi .env không có LLM_API_KEY. Domain layer
	// (AnalyzePlayerWithCoach + ChatService) đã handle nil case.
	llmCoach := llm.NewCoachFromEnv()
	var coach analysis.Coach
	if llmCoach != nil {
		coach = llmCoach
	}

	riotClient := riot.NewClient()
	riotClient.SetSettingsSink(&authSink{settings: settingsStore, reportStore: reportStore})

	return Services{
		Settings:  &SettingsService{store: settingsStore, client: henrikClient},
		Analysis:  &AnalysisService{store: settingsStore, reportStore: reportStore, henrik: henrikClient, riot: riotMatchClient, coach: coach},
		Practice:  &PracticeService{store: practiceProgress, sessionStore: practiceSessions},
		Assistant: &AssistantService{engine: assistant.NewEngine(), reportStore: reportStore},
		Tactical:  &TacticalService{store: tacticalStore},
		Auth:      &AuthService{client: riotClient},
		Chat:      &ChatService{coach: llmCoach, reportStore: reportStore},
	}
}
