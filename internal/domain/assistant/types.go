// Package assistant chứa "Live Assistant" — engine sinh tip in-game dựa trên
// Report của user (từ domain/analysis). Engine giữ state in-memory (cursor +
// cooldown timer), không persistence.
//
// Tổ chức file:
//
//	types.go   — DTO (Alert, SessionState, TipResult)
//	engine.go  — Engine struct + State / Start / Stop / Request / Poll
//	alerts.go  — BuildAlertsFromReport + alertFromFinding (mapping Finding → tip)
package assistant

// Alert là 1 gợi ý hiển thị cho user (overlay hoặc panel). Source dùng để
// trace lại finding/rule sinh ra alert này.
type Alert struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

// SessionState là state hiện tại của 1 phiên Live Assistant. Frontend dùng
// để render trạng thái + alert đang hiển thị.
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

// TipResult là kết quả mỗi lần xin tip. HasTip = false khi cooldown chưa qua
// hoặc engine chưa active — UI dùng flag này để biết có cần re-render alert.
type TipResult struct {
	HasTip bool         `json:"hasTip"`
	Alert  Alert        `json:"alert"`
	State  SessionState `json:"state"`
}
