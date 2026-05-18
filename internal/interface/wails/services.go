// Package wailsiface chứa các service được Wails Bind ra cho frontend
// (window.go.wailsiface.*). Mỗi file = 1 service để dễ tìm:
//
//	auth_service.go      — Riot Account-V1 login/logout, lấy PUUID + region.
//	settings_service.go  — Data & Consent (RiotID, region, Henrik key).
//	analysis_service.go  — Fetch + analyse match history, sinh Report.
//	practice_service.go  — Practice progress + lịch sử session luyện tập.
//	assistant_service.go — Live Assistant in-game tip engine.
//	tactical_service.go  — Map catalog + Map planner state.
//	chat_service.go      — Bot AI chat (LLM coach hỏi-đáp về report).
//	services.go          — Factory NewServices() + auth sink (file này).
package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	"valorant-tactical-trainer/desktop/internal/domain/assistant"
	"valorant-tactical-trainer/desktop/internal/infrastructure/henrik"
	"valorant-tactical-trainer/desktop/internal/infrastructure/llm"
	"valorant-tactical-trainer/desktop/internal/infrastructure/localstore"
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

// NewServices wire toàn bộ dependency: store local → infrastructure clients →
// Wails service. Lỗi khởi tạo store fallback về relative path để dev mode (vd
// `wails dev`) vẫn chạy được nếu chưa có config dir hệ điều hành.
func NewServices() Services {
	store := mustSettingsStore()
	reportStore := mustReportStore()
	practiceStore := mustPracticeProgressStore()
	sessionStore := mustPracticeSessionStore()
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
	riotClient.SetSettingsSink(&authSink{settings: store, reportStore: reportStore})

	return Services{
		Settings:  &SettingsService{store: store, client: henrikClient},
		Analysis:  &AnalysisService{store: store, reportStore: reportStore, henrik: henrikClient, riot: riotMatchClient, coach: coach},
		Practice:  &PracticeService{store: practiceStore, sessionStore: sessionStore},
		Assistant: &AssistantService{engine: assistant.NewEngine(), reportStore: reportStore},
		Tactical:  &TacticalService{store: tacticalStore},
		Auth:      &AuthService{client: riotClient},
		Chat:      &ChatService{coach: llmCoach, reportStore: reportStore},
	}
}

// --- Store factories with fallback ---

func mustSettingsStore() *localstore.SettingsStore {
	s, err := localstore.NewSettingsStore()
	if err != nil {
		return localstore.NewSettingsStoreAt("settings.json")
	}
	return s
}

func mustReportStore() *localstore.ReportStore {
	s, err := localstore.NewReportStore()
	if err != nil {
		return localstore.NewReportStoreAt("last_report.json")
	}
	return s
}

func mustPracticeProgressStore() *localstore.PracticeProgressStore {
	s, err := localstore.NewPracticeProgressStore()
	if err != nil {
		return localstore.NewPracticeProgressStoreAt("practice_progress.json")
	}
	return s
}

func mustPracticeSessionStore() *localstore.PracticeSessionStore {
	s, err := localstore.NewPracticeSessionStore()
	if err != nil {
		return localstore.NewPracticeSessionStoreAt("practice_sessions.json")
	}
	return s
}

func mustTacticalPlanStore() *localstore.TacticalPlanStore {
	s, err := localstore.NewTacticalPlanStore()
	if err != nil {
		return localstore.NewTacticalPlanStoreAt("tactical_plans.json")
	}
	return s
}
