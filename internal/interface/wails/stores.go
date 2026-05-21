package wailsiface

import "valorant-tactical-trainer/desktop/internal/infrastructure/store"

// File này gom các helper khởi tạo store local. Mỗi helper try
// UserConfigDir() → fallback về relative path để `wails dev` chạy được
// kể cả khi OS không trả config dir.

func mustSettingsStore() *store.SettingsStore {
	s, err := store.NewSettingsStore()
	if err != nil {
		return store.NewSettingsStoreAt("settings.json")
	}
	return s
}

func mustReportStore() *store.ReportStore {
	s, err := store.NewReportStore()
	if err != nil {
		return store.NewReportStoreAt("last_report.json")
	}
	return s
}

func mustPracticeProgressStore() *store.PracticeProgressStore {
	s, err := store.NewPracticeProgressStore()
	if err != nil {
		return store.NewPracticeProgressStoreAt("practice_progress.json")
	}
	return s
}

func mustPracticeSessionStore() *store.PracticeSessionStore {
	s, err := store.NewPracticeSessionStore()
	if err != nil {
		return store.NewPracticeSessionStoreAt("practice_sessions.json")
	}
	return s
}

func mustTacticalPlanStore() *store.TacticalPlanStore {
	s, err := store.NewTacticalPlanStore()
	if err != nil {
		return store.NewTacticalPlanStoreAt("tactical_plans.json")
	}
	return s
}
