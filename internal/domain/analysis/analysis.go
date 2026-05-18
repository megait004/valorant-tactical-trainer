package analysis

import (
	"math"
	"sort"
)

type PlayerSnapshot struct {
	Name          string         `json:"name"`
	Tagline       string         `json:"tagline"`
	Region        string         `json:"region"`
	PrimaryRole   string         `json:"primaryRole"`
	RecentMatches []MatchSummary `json:"recentMatches"`
}

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

type Evidence struct {
	MatchIDs           []string `json:"matchIds"`
	Map                string   `json:"map,omitempty"`
	Agent              string   `json:"agent,omitempty"`
	Metric             string   `json:"metric"`
	Value              float64  `json:"value"`
	SampleSize         int      `json:"sampleSize"`
	ComparisonBaseline float64  `json:"comparisonBaseline"`
}

type Finding struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Severity   string     `json:"severity"`
	Confidence string     `json:"confidence"`
	Detail     string     `json:"detail"`
	Evidence   []Evidence `json:"evidence"`
}

type Recommendation struct {
	ID        string `json:"id"`
	FindingID string `json:"findingId"`
	Title     string `json:"title"`
	Reason    string `json:"reason"`
	Drill     string `json:"drill"`
	Cadence   string `json:"cadence"`
}

type BreakdownRow struct {
	Name            string  `json:"name"`
	Matches         int     `json:"matches"`
	Rounds          int     `json:"rounds"`
	KD              float64 `json:"kd"`
	WinRate         float64 `json:"winRate"`
	HeadshotPercent float64 `json:"headshotPercent"`
}

type PracticeTask struct {
	Day       int      `json:"day"`
	Focus     string   `json:"focus"`
	Map       string   `json:"map,omitempty"`
	Agent     string   `json:"agent,omitempty"`
	Duration  string   `json:"duration"`
	Checklist []string `json:"checklist"`
	Evidence  string   `json:"evidence"`
}

type Report struct {
	Player          PlayerSnapshot   `json:"player"`
	Metrics         MetricSummary    `json:"metrics"`
	MapBreakdown    []BreakdownRow   `json:"mapBreakdown"`
	AgentBreakdown  []BreakdownRow   `json:"agentBreakdown"`
	PracticePlan    []PracticeTask   `json:"practicePlan"`
	Findings        []Finding        `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
}

func AnalyzePlayer(snapshot PlayerSnapshot) Report {
	metrics := calculateMetrics(snapshot)
	findings := generateFindings(snapshot, metrics)
	recommendations := generateRecommendations(findings, metrics)
	mapBreakdown := calculateBreakdown(snapshot.RecentMatches, func(match MatchSummary) string {
		return match.Map
	}, true)
	agentBreakdown := calculateBreakdown(snapshot.RecentMatches, func(match MatchSummary) string {
		return match.Agent
	}, false)

	return Report{
		Player:          snapshot,
		Metrics:         metrics,
		MapBreakdown:    mapBreakdown,
		AgentBreakdown:  agentBreakdown,
		PracticePlan:    generatePracticePlan(metrics, mapBreakdown, agentBreakdown, findings),
		Findings:        findings,
		Recommendations: recommendations,
	}
}

func DemoSnapshot() PlayerSnapshot {
	return PlayerSnapshot{
		Name:        "giaphue",
		Tagline:     "DATN",
		Region:      "ap",
		PrimaryRole: "controller",
		RecentMatches: []MatchSummary{
			{ID: "demo-001", Map: "Ascent", Agent: "Omen", Role: "controller", Kills: 14, Deaths: 18, Assists: 9, RoundsPlayed: 22, FirstBloods: 1, FirstDeaths: 5, HeadshotPercent: 16.4, Won: false},
			{ID: "demo-002", Map: "Haven", Agent: "Omen", Role: "controller", Kills: 18, Deaths: 16, Assists: 7, RoundsPlayed: 21, FirstBloods: 2, FirstDeaths: 4, HeadshotPercent: 19.1, Won: true},
			{ID: "demo-003", Map: "Ascent", Agent: "Brimstone", Role: "controller", Kills: 10, Deaths: 17, Assists: 12, RoundsPlayed: 19, FirstBloods: 0, FirstDeaths: 5, HeadshotPercent: 14.2, Won: false},
			{ID: "demo-004", Map: "Bind", Agent: "Cypher", Role: "sentinel", Kills: 21, Deaths: 14, Assists: 5, RoundsPlayed: 23, FirstBloods: 3, FirstDeaths: 2, HeadshotPercent: 22.8, Won: true},
			{ID: "demo-005", Map: "Ascent", Agent: "Omen", Role: "controller", Kills: 12, Deaths: 19, Assists: 6, RoundsPlayed: 20, FirstBloods: 1, FirstDeaths: 6, HeadshotPercent: 15.6, Won: false},
		},
	}
}

func calculateMetrics(snapshot PlayerSnapshot) MetricSummary {
	if len(snapshot.RecentMatches) == 0 {
		return MetricSummary{PrimaryRoleObserved: snapshot.PrimaryRole}
	}

	var kills, deaths, assists, rounds, firstBloods, firstDeaths, wins int
	var weightedHS float64
	roles := map[string]int{}
	mapStats := map[string]struct{ wins, games int }{}

	for _, match := range snapshot.RecentMatches {
		kills += match.Kills
		deaths += match.Deaths
		assists += match.Assists
		rounds += match.RoundsPlayed
		firstBloods += match.FirstBloods
		firstDeaths += match.FirstDeaths
		weightedHS += match.HeadshotPercent * float64(max(match.RoundsPlayed, 1))
		roles[match.Role]++
		stats := mapStats[match.Map]
		stats.games++
		if match.Won {
			wins++
			stats.wins++
		}
		mapStats[match.Map] = stats
	}

	weakestMap := ""
	weakestMapWinRate := 1.0
	weakestMapSample := 0
	// Lượt 1: ưu tiên map có ≥ 2 trận để baseline ổn định.
	for mapName, stats := range mapStats {
		if stats.games < 2 {
			continue
		}
		winRate := ratio(stats.wins, stats.games)
		if weakestMap == "" || winRate < weakestMapWinRate {
			weakestMap = mapName
			weakestMapWinRate = winRate
			weakestMapSample = stats.games
		}
	}
	// Lượt 2 (fallback): nếu không map nào đủ 2 trận (sample bé), vẫn chọn map
	// có WR thấp nhất trong các map đã chơi để UI không hiển thị "N/A". Sample
	// = 1 sẽ làm confidence rule engine xuống "low" — Finding/Recommendation
	// vẫn ưu tiên cảnh báo nhẹ thay vì biến mất.
	if weakestMap == "" {
		for mapName, stats := range mapStats {
			winRate := ratio(stats.wins, stats.games)
			if weakestMap == "" || winRate < weakestMapWinRate {
				weakestMap = mapName
				weakestMapWinRate = winRate
				weakestMapSample = stats.games
			}
		}
	}

	return MetricSummary{
		Matches:             len(snapshot.RecentMatches),
		Rounds:              rounds,
		KD:                  round2(float64(kills) / float64(max(deaths, 1))),
		KDA:                 round2(float64(kills+assists) / float64(max(deaths, 1))),
		HeadshotPercent:     round2(weightedHS / float64(max(rounds, 1))),
		FirstBloodRate:      round2(float64(firstBloods) / float64(max(rounds, 1))),
		FirstDeathRate:      round2(float64(firstDeaths) / float64(max(rounds, 1))),
		WinRate:             round2(ratio(wins, len(snapshot.RecentMatches))),
		WeakestMap:          weakestMap,
		WeakestMapWinRate:   round2(weakestMapWinRate),
		WeakestMapSample:    weakestMapSample,
		PrimaryRoleObserved: mostCommonRole(roles, snapshot.PrimaryRole),
	}
}

func calculateBreakdown(matches []MatchSummary, keyFor func(MatchSummary) string, weakestFirst bool) []BreakdownRow {
	type aggregate struct {
		matches    int
		rounds     int
		kills      int
		deaths     int
		wins       int
		weightedHS float64
	}

	statsByName := map[string]aggregate{}
	for _, match := range matches {
		name := keyFor(match)
		if name == "" {
			name = "Unknown"
		}
		stats := statsByName[name]
		stats.matches++
		stats.rounds += match.RoundsPlayed
		stats.kills += match.Kills
		stats.deaths += match.Deaths
		stats.weightedHS += match.HeadshotPercent * float64(max(match.RoundsPlayed, 1))
		if match.Won {
			stats.wins++
		}
		statsByName[name] = stats
	}

	rows := make([]BreakdownRow, 0, len(statsByName))
	for name, stats := range statsByName {
		rows = append(rows, BreakdownRow{
			Name:            name,
			Matches:         stats.matches,
			Rounds:          stats.rounds,
			KD:              round2(float64(stats.kills) / float64(max(stats.deaths, 1))),
			WinRate:         round2(ratio(stats.wins, stats.matches)),
			HeadshotPercent: round2(stats.weightedHS / float64(max(stats.rounds, 1))),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if weakestFirst && rows[i].WinRate != rows[j].WinRate {
			return rows[i].WinRate < rows[j].WinRate
		}
		if rows[i].Matches != rows[j].Matches {
			return rows[i].Matches > rows[j].Matches
		}
		return rows[i].Name < rows[j].Name
	})

	return rows
}

func generateFindings(snapshot PlayerSnapshot, metrics MetricSummary) []Finding {
	if metrics.Matches == 0 {
		return []Finding{{
			ID:         "finding-no-data",
			Title:      "Chưa có dữ liệu trận để phân tích",
			Severity:   "low",
			Confidence: "high",
			Detail:     "Cần import match history hoặc VOD trước khi sinh giáo án cá nhân.",
			Evidence: []Evidence{{
				Metric:     "matches",
				Value:      0,
				SampleSize: 0,
			}},
		}}
	}

	findings := make([]Finding, 0, len(findingPatterns))
	for _, pattern := range findingPatterns {
		evidence, ok := pattern.Match(snapshot, metrics)
		if !ok {
			continue
		}
		findings = append(findings, Finding{
			ID:         pattern.ID,
			Title:      pattern.DefaultTitle,
			Severity:   pattern.SeverityFn(metrics),
			Confidence: pattern.ConfidenceFn(metrics),
			Detail:     pattern.DefaultDetail(metrics),
			Evidence:   []Evidence{evidence},
		})
	}
	return findings
}

func generateRecommendations(findings []Finding, metrics MetricSummary) []Recommendation {
	recs := make([]Recommendation, 0, len(findings))
	for _, f := range findings {
		if f.ID == "finding-no-data" {
			recs = append(recs, Recommendation{
				ID:        "rec-import-data",
				FindingID: f.ID,
				Title:     "Import dữ liệu trước khi coach",
				Reason:    "Không có evidence thì app không nên sinh lời khuyên giả.",
				Drill:     "Login Riot và bấm Fetch report để có match history thật.",
				Cadence:   "Một lần setup",
			})
			continue
		}
		pattern := FindingPatternByID(f.ID)
		if pattern == nil || pattern.DefaultRecommendation == nil {
			continue
		}
		recs = append(recs, pattern.DefaultRecommendation(metrics))
	}
	return recs
}

func generatePracticePlan(metrics MetricSummary, mapBreakdown []BreakdownRow, agentBreakdown []BreakdownRow, findings []Finding) []PracticeTask {
	if metrics.Matches == 0 {
		return []PracticeTask{{
			Day:       1,
			Focus:     "Setup dữ liệu",
			Duration:  "10 phút",
			Checklist: []string{"Login Riot", "Bấm Fetch report", "Kiểm tra report có evidence"},
			Evidence:  "Chưa có match history để coach.",
		}}
	}

	weakMap := firstBreakdownName(mapBreakdown)
	mainAgent := firstBreakdownName(agentBreakdown)
	plan := []PracticeTask{}

	for _, f := range findings {
		pattern := FindingPatternByID(f.ID)
		if pattern == nil || pattern.DefaultPracticeTask == nil {
			continue
		}
		plan = append(plan, pattern.DefaultPracticeTask(metrics, weakMap, mainAgent, len(plan)+1))
	}

	if len(plan) == 0 {
		plan = append(plan, PracticeTask{
			Day:      1,
			Focus:    "Maintain form",
			Map:      weakMap,
			Agent:    mainAgent,
			Duration: "20 phút",
			Checklist: []string{
				"Warmup aim 10 phút",
				"Review 1 win và 1 loss gần nhất",
				"Ghi 1 adjustment cho map đang gặp nhiều nhất",
			},
			Evidence: "Không có finding vượt ngưỡng MVP trong sample hiện tại.",
		})
	}

	if len(plan) > 4 {
		return plan[:4]
	}
	return plan
}

func firstBreakdownName(rows []BreakdownRow) string {
	if len(rows) == 0 {
		return ""
	}
	return rows[0].Name
}

func collectMatchIDs(matches []MatchSummary) []string {
	ids := make([]string, 0, len(matches))
	for _, match := range matches {
		ids = append(ids, match.ID)
	}
	return ids
}

func matchesForMap(matches []MatchSummary, mapName string) []string {
	ids := []string{}
	for _, match := range matches {
		if match.Map == mapName {
			ids = append(ids, match.ID)
		}
	}
	return ids
}

func mostCommonRole(roles map[string]int, fallback string) string {
	type roleCount struct {
		role  string
		count int
	}
	values := make([]roleCount, 0, len(roles))
	for role, count := range roles {
		values = append(values, roleCount{role: role, count: count})
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].count == values[j].count {
			return values[i].role < values[j].role
		}
		return values[i].count > values[j].count
	})
	if len(values) == 0 {
		return fallback
	}
	return values[0].role
}

func roleDisplayName(role string) string {
	switch role {
	case "duelist":
		return "Đối Đầu (Duelist)"
	case "initiator":
		return "Khởi Tranh (Initiator)"
	case "controller":
		return "Kiểm Soát (Controller)"
	case "sentinel":
		return "Hộ Vệ (Sentinel)"
	default:
		return role
	}
}

func confidence(sampleSize int, target int) string {
	if sampleSize >= target {
		return "high"
	}
	if sampleSize >= target/2 {
		return "medium"
	}
	return "low"
}

func severity(value float64, highThreshold float64) string {
	if value >= highThreshold {
		return "high"
	}
	return "medium"
}

func ratio(value int, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(value) / float64(total)
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
