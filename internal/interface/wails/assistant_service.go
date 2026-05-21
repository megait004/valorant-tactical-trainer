package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	"valorant-tactical-trainer/desktop/internal/domain/assistant"
	"valorant-tactical-trainer/desktop/internal/infrastructure/store"
)

// AssistantService expose Live Assistant (in-game tip) qua Wails. Engine giữ
// state in-memory; reportStore dùng để load context khi start session.
type AssistantService struct {
	engine      *assistant.Engine
	reportStore *store.ReportStore
}

func (s *AssistantService) GetSessionState() assistant.SessionState {
	return s.engine.State()
}

func (s *AssistantService) StartSession() (assistant.SessionState, error) {
	return s.engine.Start(s.loadReport()), nil
}

func (s *AssistantService) StopSession() assistant.SessionState {
	return s.engine.Stop()
}

func (s *AssistantService) RequestTip() assistant.TipResult {
	return s.engine.RequestTip()
}

func (s *AssistantService) MarkRoundStart() assistant.TipResult {
	return s.engine.MarkRoundStart()
}

func (s *AssistantService) PollAutoTip() assistant.TipResult {
	return s.engine.PollAutoTip()
}

// loadReport lấy report gần nhất làm context cho engine. Nếu chưa có thì rơi
// về DemoSnapshot — engine vẫn chạy được với rule heuristic mặc định.
func (s *AssistantService) loadReport() analysis.Report {
	stored, ok, err := s.reportStore.LoadLastReport()
	if err != nil || !ok {
		return analysis.AnalyzePlayer(analysis.DemoSnapshot())
	}
	report := stored.Report
	if report.PracticePlan == nil || report.MapBreakdown == nil || report.AgentBreakdown == nil {
		report = analysis.AnalyzePlayer(report.Player)
	}
	return report
}
