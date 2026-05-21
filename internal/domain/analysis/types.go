// Package analysis chứa rule engine phân tích match history Valorant của 1
// player và sinh ra Report cá nhân hoá (Findings, Recommendations, Practice
// Plan, Map/Agent breakdown). Toàn bộ code trong package là pure logic —
// không I/O, không HTTP — để dễ unit test deterministic.
//
// Tổ chức file:
//
//	types.go         — Toàn bộ struct DTO (PlayerSnapshot, Report, ...)
//	analyze.go       — Entry chính: AnalyzePlayer + AnalyzePlayerWithCoach hook
//	metrics.go       — Tính metrics tổng (KD, HS%, win rate, weakest map)
//	breakdown.go     — Tính map/agent breakdown
//	findings.go      — Sinh Findings + Recommendations từ pattern registry
//	practice_plan.go — Sinh giáo án luyện tập (PracticeTask) 4 ngày
//	pattern.go       — Type FindingPattern + registry helpers
//	pattern_catalog.go — Định nghĩa cụ thể từng pattern (first-death, HS, ...)
//	coach.go         — Coach interface (gọi LLM) + merge output
//	demo.go          — DemoSnapshot dùng khi user chưa fetch report thật
//	helpers.go       — Tiện ích nhỏ: ratio, round2, severity, confidence...
package analysis

// PlayerSnapshot là input gốc cho rule engine — match history + thông tin player.
type PlayerSnapshot struct {
	Name          string         `json:"name"`
	Tagline       string         `json:"tagline"`
	Region        string         `json:"region"`
	PrimaryRole   string         `json:"primaryRole"`
	RecentMatches []MatchSummary `json:"recentMatches"`
}

// MatchSummary là 1 trận trong match history. Một số trường (FirstBloods,
// FirstDeaths) có thể là 0 nếu adapter API không trả về dữ liệu đó.
type MatchSummary struct {
	ID              string  `json:"id"`
	Map             string  `json:"map"`
	Agent           string  `json:"agent"`
	Role            string  `json:"role"`
	Kills           int     `json:"kills"`
	Deaths          int     `json:"deaths"`
	Assists         int     `json:"assists"`
	RoundsPlayed    int     `json:"roundsPlayed"`
	FirstBloods     int     `json:"firstBloods"`
	FirstDeaths     int     `json:"firstDeaths"`
	HeadshotPercent float64 `json:"headshotPercent"`
	Won             bool    `json:"won"`
}

// MetricSummary là số liệu tổng đã được aggregate từ RecentMatches.
type MetricSummary struct {
	Matches             int     `json:"matches"`
	Rounds              int     `json:"rounds"`
	KD                  float64 `json:"kd"`
	KDA                 float64 `json:"kda"`
	HeadshotPercent     float64 `json:"headshotPercent"`
	FirstBloodRate      float64 `json:"firstBloodRate"`
	FirstDeathRate      float64 `json:"firstDeathRate"`
	WinRate             float64 `json:"winRate"`
	WeakestMap          string  `json:"weakestMap"`
	WeakestMapWinRate   float64 `json:"weakestMapWinRate"`
	WeakestMapSample    int     `json:"weakestMapSample"`
	PrimaryRoleObserved string  `json:"primaryRoleObserved"`
}

// Evidence gắn vào mỗi Finding để UI/coach có thể chỉ ra số liệu cụ thể.
type Evidence struct {
	MatchIDs           []string `json:"matchIds"`
	Map                string   `json:"map,omitempty"`
	Agent              string   `json:"agent,omitempty"`
	Metric             string   `json:"metric"`
	Value              float64  `json:"value"`
	SampleSize         int      `json:"sampleSize"`
	ComparisonBaseline float64  `json:"comparisonBaseline"`
}

// Finding là 1 vấn đề kỹ thuật rule engine phát hiện được. ID dùng để liên
// kết với Recommendation và PracticeTask cùng ID gốc (xem pattern.go).
type Finding struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Severity   string     `json:"severity"`
	Confidence string     `json:"confidence"`
	Detail     string     `json:"detail"`
	Evidence   []Evidence `json:"evidence"`
}

// Recommendation là lời khuyên cụ thể đính kèm từng Finding.
type Recommendation struct {
	ID        string `json:"id"`
	FindingID string `json:"findingId"`
	Title     string `json:"title"`
	Reason    string `json:"reason"`
	Drill     string `json:"drill"`
	Cadence   string `json:"cadence"`
}

// BreakdownRow là 1 dòng breakdown theo map hoặc agent.
type BreakdownRow struct {
	Name            string  `json:"name"`
	Matches         int     `json:"matches"`
	Rounds          int     `json:"rounds"`
	KD              float64 `json:"kd"`
	WinRate         float64 `json:"winRate"`
	HeadshotPercent float64 `json:"headshotPercent"`
}

// PracticeTask là 1 bài tập trong giáo án (1-4 ngày).
type PracticeTask struct {
	Day       int      `json:"day"`
	Focus     string   `json:"focus"`
	Map       string   `json:"map,omitempty"`
	Agent     string   `json:"agent,omitempty"`
	Duration  string   `json:"duration"`
	Checklist []string `json:"checklist"`
	Evidence  string   `json:"evidence"`
}

// Report là output đầy đủ của 1 lần AnalyzePlayer.
type Report struct {
	Player          PlayerSnapshot   `json:"player"`
	Metrics         MetricSummary    `json:"metrics"`
	MapBreakdown    []BreakdownRow   `json:"mapBreakdown"`
	AgentBreakdown  []BreakdownRow   `json:"agentBreakdown"`
	PracticePlan    []PracticeTask   `json:"practicePlan"`
	Findings        []Finding        `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
}
