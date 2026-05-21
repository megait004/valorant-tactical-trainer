package analysis

// DemoSnapshot trả về 1 PlayerSnapshot mẫu để UI có data hiển thị trước khi
// user fetch report thật. Cũng dùng cho test deterministic.
//
// Lưu ý: ID match có prefix "demo-" để frontend phân biệt với match thật từ
// Henrik/Riot khi cần render badge "demo".
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
