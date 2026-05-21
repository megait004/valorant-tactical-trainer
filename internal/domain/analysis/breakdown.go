package analysis

import "sort"

// calculateBreakdown group match theo key (map hoặc agent) và tổng hợp số
// liệu mỗi nhóm.
//
//   - keyFor:        hàm trích key từ MatchSummary (vd: m.Map hoặc m.Agent)
//   - weakestFirst:  sort tăng dần theo WinRate (true → map yếu lên đầu)
//     hoặc giảm dần theo số trận (false → agent dùng nhiều lên đầu).
func calculateBreakdown(matches []MatchSummary, keyFor func(MatchSummary) string, weakestFirst bool) []BreakdownRow {
	type aggregate struct {
		matches, rounds, kills, deaths, wins int
		weightedHS                           float64
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

// firstBreakdownName trả về tên row đầu tiên (sau sort) — dùng làm "weak map"
// / "main agent" trong PracticePlan placeholder. Rỗng nếu không có row.
func firstBreakdownName(rows []BreakdownRow) string {
	if len(rows) == 0 {
		return ""
	}
	return rows[0].Name
}
