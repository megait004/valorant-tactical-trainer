package assistant

import (
	"fmt"
	"strings"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

const defaultCooldown = 25 * time.Second

// Engine quản lý phiên Live Assistant — giữ queue alert, cooldown timer và
// session state. Không thread-safe (Wails service gọi tuần tự từ frontend).
type Engine struct {
	alerts      []Alert
	cursor      int
	lastShownAt time.Time
	cooldown    time.Duration
	state       SessionState
	report      analysis.Report
}

// NewEngine tạo engine mới với cooldown mặc định 25s.
func NewEngine() *Engine {
	return &Engine{cooldown: defaultCooldown}
}

// State trả về session state hiện tại — frontend poll khi mở panel.
func (e *Engine) State() SessionState {
	return e.state
}

// Start khởi tạo phiên mới với report được pass vào, build queue alert và
// hiển thị ngay alert đầu tiên (force, bypass cooldown).
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

// Stop tắt phiên, clear alert hiện tại.
func (e *Engine) Stop() SessionState {
	e.state.Active = false
	e.state.CurrentAlert = nil
	e.state.Message = "Đã tắt Live Assistant."
	return e.state
}

// RequestTip force lấy tip tiếp theo (bypass cooldown). Dùng khi user bấm
// "Nhắc tôi" hoặc phím tắt H.
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

// MarkRoundStart tăng round counter và hiển thị alert phù hợp round đó.
// Dùng khi user bấm "Round mới" hoặc phím tắt R.
func (e *Engine) MarkRoundStart() TipResult {
	if !e.state.Active {
		return TipResult{HasTip: false, State: e.state}
	}
	e.state.RoundCount++
	alert := e.roundAlert(e.state.RoundCount, e.report)
	e.applyAlert(alert)
	return TipResult{HasTip: true, Alert: alert, State: e.state}
}

// PollAutoTip lấy tip tiếp theo nếu cooldown đã qua. Frontend gọi định kỳ
// mỗi 8s để engine tự nhắc khi user đang focus game.
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

// roundAlert tuỳ chỉnh alert theo round number. Round 1-3 (pistol/eco sớm)
// có message riêng theo role; round sau dùng queue alert thông thường.
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

	if alert, ok := e.nextAlert(true); ok {
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

// nextAlert pop alert kế tiếp từ queue (circular). force = true bypass
// cooldown.
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

// applyAlert ghi alert vào state + reset timer.
func (e *Engine) applyAlert(alert Alert) {
	now := time.Now().UTC()
	e.lastShownAt = now
	e.state.LastAlertAt = now.Format(time.RFC3339)
	e.state.TipsShown++
	copy := alert
	e.state.CurrentAlert = &copy
	e.state.Message = "Đang hiển thị gợi ý."
}

// cooldownPassed = true nếu chưa từng hiển thị, hoặc đã qua cooldown.
func (e *Engine) cooldownPassed() bool {
	if e.lastShownAt.IsZero() {
		return true
	}
	return time.Since(e.lastShownAt) >= e.cooldown
}

// truncate cắt chuỗi để alert message không vượt quá độ dài giới hạn UI.
func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max-3] + "..."
}
