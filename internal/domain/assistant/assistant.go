package assistant

import (
	"fmt"
	"strings"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

const defaultCooldown = 25 * time.Second

type Alert struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

type SessionState struct {
	Active       bool   `json:"active"`
	StartedAt    string `json:"startedAt"`
	RoundCount   int    `json:"roundCount"`
	TipsShown    int    `json:"tipsShown"`
	LastAlertAt  string `json:"lastAlertAt"`
	CurrentAlert *Alert `json:"currentAlert"`
	Message      string `json:"message"`
	QueueSize    int    `json:"queueSize"`
}

type TipResult struct {
	HasTip bool        `json:"hasTip"`
	Alert  Alert       `json:"alert"`
	State  SessionState `json:"state"`
}

type Engine struct {
	alerts      []Alert
	cursor      int
	lastShownAt time.Time
	cooldown    time.Duration
	state       SessionState
	report      analysis.Report
}

func NewEngine() *Engine {
	return &Engine{cooldown: defaultCooldown}
}

func BuildAlertsFromReport(report analysis.Report) []Alert {
	alerts := make([]Alert, 0, 8)

	if len(report.Recommendations) > 0 {
		rec := report.Recommendations[0]
		alerts = append(alerts, Alert{
			ID:       "opener-recommendation",
			Title:    "Mục tiêu phiên",
			Message:  truncate(rec.Title+" — "+rec.Drill, 160),
			Severity: "medium",
			Source:   rec.FindingID,
		})
	}

	for _, finding := range report.Findings {
		if alert, ok := alertFromFinding(finding, report); ok {
			alerts = append(alerts, alert)
		}
	}

	if len(report.PracticePlan) > 0 {
		task := report.PracticePlan[0]
		checklist := ""
		if len(task.Checklist) > 0 {
			checklist = task.Checklist[0]
		}
		alerts = append(alerts, Alert{
			ID:       "practice-focus",
			Title:    "Bài tập hôm nay",
			Message:  truncate(task.Focus+": "+checklist, 160),
			Severity: "medium",
			Source:   "practice-plan",
		})
	}

	if report.Metrics.WeakestMap != "" {
		alerts = append(alerts, Alert{
			ID:       "map-awareness",
			Title:    "Map cần chú ý",
			Message:  fmt.Sprintf("%s — chơi default đã luyện, tránh pick fight solo.", report.Metrics.WeakestMap),
			Severity: "low",
			Source:   "map-gap",
		})
	}

	if len(alerts) == 0 {
		alerts = append(alerts, Alert{
			ID:       "default-maintain",
			Title:    "Giữ form",
			Message:  "Utility trước peek, gọi trade trước contact chính.",
			Severity: "low",
			Source:   "default",
		})
	}

	return dedupeAlerts(alerts)
}

func alertFromFinding(finding analysis.Finding, report analysis.Report) (Alert, bool) {
	switch finding.ID {
	case "finding-first-death-non-duelist":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Round mới: smoke/info trước khi peek; nếu không phải Đối Đầu thì đừng entry sớm.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-low-headshot":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Pre-aim góc default, tap/burst — tránh wide swing khi crosshair chưa đặt đầu.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-low-survival-impact":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Chỉ nhận duel khi có trade path — sau death hỏi: có teammate trade trong 3 giây?",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-map-gap":
		mapName := report.Metrics.WeakestMap
		if mapName == "" {
			mapName = "map yếu"
		}
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  fmt.Sprintf("Ưu tiên protocol %s — default attack/defense, không hero pick mid.", mapName),
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-no-data":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Chưa có report thật — vào Settings bật consent và fetch Henrik trước khi luyện.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	default:
		return Alert{}, false
	}
}

func (e *Engine) Start(report analysis.Report) SessionState {
	e.report = report
	e.alerts = BuildAlertsFromReport(report)
	e.cursor = 0
	e.lastShownAt = time.Time{}
	now := time.Now().UTC().Format(time.RFC3339)
	e.state = SessionState{
		Active:     true,
		StartedAt:  now,
		RoundCount: 0,
		Message:    "Live Assistant đang bật. Chuyển overlay và vào trận.",
		QueueSize:  len(e.alerts),
	}
	first, _ := e.nextAlert(true)
	if first.ID != "" {
		e.applyAlert(first)
	}
	return e.state
}

func (e *Engine) Stop() SessionState {
	e.state.Active = false
	e.state.CurrentAlert = nil
	e.state.Message = "Đã tắt Live Assistant."
	return e.state
}

func (e *Engine) State() SessionState {
	return e.state
}

func (e *Engine) RequestTip() TipResult {
	if !e.state.Active {
		return TipResult{HasTip: false, State: e.state}
	}
	alert, ok := e.nextAlert(true)
	if !ok {
		return TipResult{HasTip: false, State: e.state}
	}
	e.applyAlert(alert)
	return TipResult{HasTip: true, Alert: alert, State: e.state}
}

func (e *Engine) MarkRoundStart() TipResult {
	if !e.state.Active {
		return TipResult{HasTip: false, State: e.state}
	}
	e.state.RoundCount++
	alert := e.roundAlert(e.state.RoundCount, e.report)
	e.applyAlert(alert)
	return TipResult{HasTip: true, Alert: alert, State: e.state}
}

func (e *Engine) PollAutoTip() TipResult {
	if !e.state.Active {
		return TipResult{HasTip: false, State: e.state}
	}
	if !e.cooldownPassed() {
		return TipResult{HasTip: false, State: e.state}
	}
	alert, ok := e.nextAlert(false)
	if !ok {
		return TipResult{HasTip: false, State: e.state}
	}
	e.applyAlert(alert)
	return TipResult{HasTip: true, Alert: alert, State: e.state}
}

func (e *Engine) roundAlert(round int, report analysis.Report) Alert {
	if round <= 3 {
		role := report.Metrics.PrimaryRoleObserved
		if role == "" {
			role = report.Player.PrimaryRole
		}
		if role != "duelist" {
			return Alert{
				ID:       fmt.Sprintf("round-%d-pistol", round),
				Title:    fmt.Sprintf("Round %d — pistol", round),
				Message:  "Pistol: giữ utility, đợi info — không peek đầu round một mình.",
				Severity: "medium",
				Source:   "round-pistol",
			}
		}
		return Alert{
			ID:       fmt.Sprintf("round-%d-pistol", round),
			Title:    fmt.Sprintf("Round %d — pistol", round),
			Message:  "Pistol: tạo space cho team, trade nhanh sau contact đầu.",
			Severity: "medium",
			Source:   "round-pistol",
		}
	}

	alert, ok := e.nextAlert(true)
	if ok {
		return alert
	}
	return Alert{
		ID:       fmt.Sprintf("round-%d-maintain", round),
		Title:    fmt.Sprintf("Round %d", round),
		Message:  "Giữ discipline: utility trước peek, call rotate sớm.",
		Severity: "low",
		Source:   "round-maintain",
	}
}

func (e *Engine) nextAlert(force bool) (Alert, bool) {
	if len(e.alerts) == 0 {
		return Alert{}, false
	}
	if !force && !e.cooldownPassed() {
		return Alert{}, false
	}
	alert := e.alerts[e.cursor%len(e.alerts)]
	e.cursor++
	return alert, true
}

func (e *Engine) applyAlert(alert Alert) {
	now := time.Now().UTC()
	e.lastShownAt = now
	e.state.LastAlertAt = now.Format(time.RFC3339)
	e.state.TipsShown++
	copy := alert
	e.state.CurrentAlert = &copy
	e.state.Message = "Đang hiển thị gợi ý."
}

func (e *Engine) cooldownPassed() bool {
	if e.lastShownAt.IsZero() {
		return true
	}
	return time.Since(e.lastShownAt) >= e.cooldown
}

func dedupeAlerts(alerts []Alert) []Alert {
	seen := map[string]bool{}
	out := make([]Alert, 0, len(alerts))
	for _, alert := range alerts {
		if seen[alert.ID] {
			continue
		}
		seen[alert.ID] = true
		out = append(out, alert)
	}
	return out
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max-3] + "..."
}
