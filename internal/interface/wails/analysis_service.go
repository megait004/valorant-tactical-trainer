package wailsiface

import (
	"context"
	"fmt"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
	"valorant-tactical-trainer/desktop/internal/infrastructure/henrik"
	"valorant-tactical-trainer/desktop/internal/infrastructure/localstore"
	"valorant-tactical-trainer/desktop/internal/infrastructure/riot"
)

// AnalysisService chịu trách nhiệm fetch match history (Riot VAL-MATCH-V1
// hoặc Henrik fallback), gọi rule engine (analysis.AnalyzePlayer) và LLM coach
// để cá nhân hoá Recommendations.
type AnalysisService struct {
	store       *localstore.SettingsStore
	reportStore *localstore.ReportStore
	henrik      *henrik.Client
	riot        *riot.MatchClient
	coach       analysis.Coach
}

// LiveAnalysisResult là DTO trả về frontend mỗi lần fetch report.
type LiveAnalysisResult struct {
	Report    analysis.Report `json:"report"`
	Source    string          `json:"source"`
	Cached    bool            `json:"cached"`
	FetchedAt string          `json:"fetchedAt"`
	Message   string          `json:"message"`
}

// LastReportResult bọc LiveAnalysisResult cho endpoint GetLastReport (có flag
// hasReport để UI biết có nên hiện demo hay không).
type LastReportResult struct {
	HasReport bool               `json:"hasReport"`
	Result    LiveAnalysisResult `json:"result"`
}

// GenerateDemoReport sinh report từ DemoSnapshot — dùng khi user chưa fetch
// lần đầu, để UI có gì đó hiển thị.
func (s *AnalysisService) GenerateDemoReport() analysis.Report {
	return analysis.AnalyzePlayer(analysis.DemoSnapshot())
}

// AnalyzeSnapshot exposed cho debug/test — analyse 1 snapshot tuỳ ý.
func (s *AnalysisService) AnalyzeSnapshot(snapshot analysis.PlayerSnapshot) analysis.Report {
	return analysis.AnalyzePlayer(snapshot)
}

// GetLastReport load last_report.json (nếu có). Re-analyse phần thiếu để
// đảm bảo report luôn có MapBreakdown/AgentBreakdown/PracticePlan.
func (s *AnalysisService) GetLastReport() (LastReportResult, error) {
	stored, ok, err := s.reportStore.LoadLastReport()
	if err != nil {
		return LastReportResult{}, err
	}
	if !ok {
		return LastReportResult{}, nil
	}
	report := stored.Report
	if report.PracticePlan == nil || report.MapBreakdown == nil || report.AgentBreakdown == nil {
		report = analysis.AnalyzePlayer(report.Player)
	}
	return LastReportResult{HasReport: true, Result: LiveAnalysisResult{
		Report:    report,
		Source:    stored.Source,
		Cached:    stored.Cached,
		FetchedAt: stored.FetchedAt,
		Message:   stored.Message,
	}}, nil
}

// FetchLiveReport lấy match history theo thứ tự ưu tiên Riot → Henrik, sau
// đó analyse + cá nhân hoá bằng LLM coach (nếu có).
func (s *AnalysisService) FetchLiveReport() (LiveAnalysisResult, error) {
	settings, err := s.store.LoadDataSettings()
	if err != nil {
		return LiveAnalysisResult{}, err
	}

	snapshot, source, cached, fetchedAt, message, err := s.fetchSnapshot(context.Background(), settings)
	if err != nil {
		return LiveAnalysisResult{}, err
	}

	// Coach optional — nil thì AnalyzePlayerWithCoach tự dùng template.
	report := analysis.AnalyzePlayerWithCoach(context.Background(), snapshot, s.coach)

	live := LiveAnalysisResult{
		Report:    report,
		Source:    source,
		Cached:    cached,
		FetchedAt: fetchedAt,
		Message:   message,
	}
	if err := s.reportStore.SaveLastReport(localstore.StoredReport{
		Report:    live.Report,
		Source:    live.Source,
		Cached:    live.Cached,
		FetchedAt: live.FetchedAt,
		Message:   live.Message,
	}); err != nil {
		return LiveAnalysisResult{}, err
	}
	return live, nil
}

// fetchSnapshot ưu tiên VAL-MATCH-V1 (chính thống, fetch theo PUUID) → fallback
// Henrik public (theo name#tag) khi Riot trả 403/lỗi.
func (s *AnalysisService) fetchSnapshot(
	ctx context.Context,
	settings datasettings.DataSettings,
) (analysis.PlayerSnapshot, string, bool, string, string, error) {
	if s.riot != nil {
		if r, err := s.riot.FetchMatchSnapshot(ctx, settings); err == nil {
			return r.Snapshot, r.Source, r.Cached, r.FetchedAt, r.Message, nil
		} else {
			h, hErr := s.henrik.FetchMatchSnapshot(ctx, settings)
			if hErr != nil {
				return analysis.PlayerSnapshot{}, "", false, "", "",
					fmt.Errorf("Riot lỗi: %v · Henrik lỗi: %w", err, hErr)
			}
			return h.Snapshot, h.Source, h.Cached, h.FetchedAt,
				fmt.Sprintf("Riot bỏ qua (%v) → fallback Henrik. %s", err, h.Message), nil
		}
	}
	h, err := s.henrik.FetchMatchSnapshot(ctx, settings)
	if err != nil {
		return analysis.PlayerSnapshot{}, "", false, "", "", err
	}
	return h.Snapshot, h.Source, h.Cached, h.FetchedAt, h.Message, nil
}
