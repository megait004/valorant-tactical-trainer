package analysis

import (
	"fmt"
	"math"
	"sort"
	"time"

	matchdomain "valorant-tactical-trainer/internal/domain/match"
)

type Report struct {
	ID              int64            `json:"id"`
	PlayerPUUID     string           `json:"playerPuuid"`
	GeneratedAt     time.Time        `json:"generatedAt"`
	MatchCount      int              `json:"matchCount"`
	AverageKDA      float64          `json:"averageKda"`
	HeadshotPercent float64          `json:"headshotPercent"`
	AverageDamage   float64          `json:"averageDamage"`
	TopAgent        string           `json:"topAgent"`
	TopMap          string           `json:"topMap"`
	Summary         string           `json:"summary"`
	Findings        []Finding        `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Finding struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Confidence  float64  `json:"confidence"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
}

type Recommendation struct {
	Title    string   `json:"title"`
	Drill    string   `json:"drill"`
	Priority string   `json:"priority"`
	Reason   string   `json:"reason"`
	Evidence []string `json:"evidence"`
	Status   string   `json:"status"`
}

func GenerateReport(playerPUUID string, matches []matchdomain.Summary) Report {
	report := Report{
		PlayerPUUID: playerPUUID,
		GeneratedAt: time.Now().UTC(),
		MatchCount:  len(matches),
		Summary:     "Need more match data before analysis is reliable.",
	}

	if len(matches) == 0 {
		report.Findings = append(report.Findings, Finding{
			Type:        "sample-size",
			Severity:    "medium",
			Confidence:  0.9,
			Title:       "No stored matches yet",
			Description: "Refresh match history before generating tactical conclusions.",
			Evidence:    []string{"0 stored matches"},
		})
		report.Recommendations = append(report.Recommendations, Recommendation{
			Title:    "Import recent matches",
			Drill:    "Lookup player, refresh 10 recent matches, then regenerate report.",
			Priority: "high",
			Reason:   "The analysis engine needs match evidence.",
			Evidence: []string{"match_count=0"},
			Status:   "new",
		})
		return report
	}

	agentCount := map[string]int{}
	mapCount := map[string]int{}
	totalKills := 0
	totalDeaths := 0
	totalAssists := 0
	totalHeadshots := 0
	totalBodyshots := 0
	totalLegshots := 0
	totalDamage := 0

	for _, match := range matches {
		totalKills += match.Kills
		totalDeaths += match.Deaths
		totalAssists += match.Assists
		totalHeadshots += match.Headshots
		totalBodyshots += match.Bodyshots
		totalLegshots += match.Legshots
		totalDamage += match.DamageMade

		if match.Agent != "" {
			agentCount[match.Agent]++
		}
		if match.MapName != "" {
			mapCount[match.MapName]++
		}
	}

	report.AverageKDA = round2(float64(totalKills+totalAssists) / float64(max(totalDeaths, 1)))
	report.HeadshotPercent = round2(float64(totalHeadshots) * 100 / float64(max(totalHeadshots+totalBodyshots+totalLegshots, 1)))
	report.AverageDamage = round2(float64(totalDamage) / float64(len(matches)))
	report.TopAgent = topKey(agentCount)
	report.TopMap = topKey(mapCount)
	report.Summary = fmt.Sprintf("Analyzed %d matches. KDA %.2f, HS %.1f%%, average damage %.0f.", report.MatchCount, report.AverageKDA, report.HeadshotPercent, report.AverageDamage)

	if report.HeadshotPercent < 18 && totalHeadshots+totalBodyshots+totalLegshots >= 80 {
		finding := Finding{
			Type:        "aim-discipline",
			Severity:    severityFromThreshold(report.HeadshotPercent, 12, 18),
			Confidence:  confidenceBySample(len(matches)),
			Title:       "Headshot rate is below target",
			Description: "Recent matches suggest crosshair placement or first-bullet discipline needs work.",
			Evidence:    []string{fmt.Sprintf("headshot_percent=%.1f", report.HeadshotPercent), fmt.Sprintf("shots=%d", totalHeadshots+totalBodyshots+totalLegshots)},
		}
		report.Findings = append(report.Findings, finding)
		report.Recommendations = append(report.Recommendations, Recommendation{
			Title:    "10-minute crosshair discipline block",
			Drill:    "Run 2 deathmatch games focusing only on head-level crosshair and burst reset. Review first deaths after each game.",
			Priority: finding.Severity,
			Reason:   finding.Title,
			Evidence: finding.Evidence,
			Status:   "new",
		})
	}

	if report.AverageKDA < 1.1 && len(matches) >= 3 {
		finding := Finding{
			Type:        "survivability",
			Severity:    severityFromThreshold(report.AverageKDA, 0.8, 1.1),
			Confidence:  confidenceBySample(len(matches)),
			Title:       "Low KDA across recent matches",
			Description: "The player may be taking low-value duels or missing trade timing.",
			Evidence:    []string{fmt.Sprintf("average_kda=%.2f", report.AverageKDA), fmt.Sprintf("matches=%d", len(matches))},
		}
		report.Findings = append(report.Findings, finding)
		report.Recommendations = append(report.Recommendations, Recommendation{
			Title:    "Trade spacing review",
			Drill:    "Review 5 deaths. Label each as isolated duel, late trade, utility gap, or timing error. Queue with a duo and call trade positions each round.",
			Priority: finding.Severity,
			Reason:   finding.Title,
			Evidence: finding.Evidence,
			Status:   "new",
		})
	}

	if report.TopMap != "" && mapCount[report.TopMap] >= max(2, len(matches)/2) {
		finding := Finding{
			Type:        "map-sample",
			Severity:    "low",
			Confidence:  0.6,
			Title:       fmt.Sprintf("Most data comes from %s", report.TopMap),
			Description: "Map-specific conclusions should be separated from global player conclusions.",
			Evidence:    []string{fmt.Sprintf("%s_matches=%d", report.TopMap, mapCount[report.TopMap])},
		}
		report.Findings = append(report.Findings, finding)
		report.Recommendations = append(report.Recommendations, Recommendation{
			Title:    fmt.Sprintf("Build %s round notes", report.TopMap),
			Drill:    "Write 3 attack defaults, 3 defensive setups, and 2 retake utility rules for this map.",
			Priority: "medium",
			Reason:   finding.Title,
			Evidence: finding.Evidence,
			Status:   "new",
		})
	}

	if len(report.Findings) == 0 {
		report.Findings = append(report.Findings, Finding{
			Type:        "baseline",
			Severity:    "low",
			Confidence:  confidenceBySample(len(matches)),
			Title:       "No critical leak detected in MVP rules",
			Description: "The current rule set did not find a major issue. Import more matches or expand rules for economy/round analysis.",
			Evidence:    []string{fmt.Sprintf("matches=%d", len(matches))},
		})
		report.Recommendations = append(report.Recommendations, Recommendation{
			Title:    "Expand review sample",
			Drill:    "Import at least 10 more matches, then compare agent/map-specific trends.",
			Priority: "low",
			Reason:   "No high-confidence leak found yet.",
			Evidence: []string{fmt.Sprintf("matches=%d", len(matches))},
			Status:   "new",
		})
	}

	sort.Slice(report.Findings, func(i, j int) bool {
		return severityRank(report.Findings[i].Severity) > severityRank(report.Findings[j].Severity)
	})
	sort.Slice(report.Recommendations, func(i, j int) bool {
		return severityRank(report.Recommendations[i].Priority) > severityRank(report.Recommendations[j].Priority)
	})

	return report
}

func topKey(values map[string]int) string {
	top := ""
	topCount := 0
	for key, count := range values {
		if count > topCount {
			top = key
			topCount = count
		}
	}
	return top
}

func confidenceBySample(count int) float64 {
	return round2(math.Min(0.95, 0.45+float64(count)*0.08))
}

func severityFromThreshold(value float64, critical float64, warning float64) string {
	if value <= critical {
		return "high"
	}
	if value <= warning {
		return "medium"
	}
	return "low"
}

func severityRank(value string) int {
	switch value {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
