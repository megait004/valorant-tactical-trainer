package analysis

import "sort"

// calculateMetrics aggregate match list thành MetricSummary, đồng thời tìm ra
// "map yếu" (map có winrate thấp nhất). Có 2 pass:
//
//  1. Ưu tiên map có ≥ 2 trận để baseline ổn định.
//  2. Fallback chọn map WR thấp nhất trong sample bé (≥ 1 trận) — để UI không
//     hiển thị "N/A" cho user mới chỉ chơi vài trận. Confidence rule engine
//     sẽ tự xuống "low" trong trường hợp này.
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

	weakestMap, weakestMapWinRate, weakestMapSample := pickWeakestMap(mapStats)

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

// pickWeakestMap chọn map có winrate thấp nhất. Ưu tiên map ≥ 2 trận để
// baseline đáng tin, fallback về sample bé nếu không có map nào đủ 2 trận.
func pickWeakestMap(mapStats map[string]struct{ wins, games int }) (name string, winRate float64, sample int) {
	winRate = 1.0
	for mapName, stats := range mapStats {
		if stats.games < 2 {
			continue
		}
		wr := ratio(stats.wins, stats.games)
		if name == "" || wr < winRate {
			name = mapName
			winRate = wr
			sample = stats.games
		}
	}
	if name != "" {
		return name, winRate, sample
	}
	for mapName, stats := range mapStats {
		wr := ratio(stats.wins, stats.games)
		if name == "" || wr < winRate {
			name = mapName
			winRate = wr
			sample = stats.games
		}
	}
	return name, winRate, sample
}

// mostCommonRole tìm role xuất hiện nhiều nhất trong match history, fallback
// về PrimaryRole user khai báo nếu sample rỗng.
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
