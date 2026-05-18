package llm

import "valorant-tactical-trainer/desktop/internal/domain/analysis"

// promptPayload là input gọn cho LLM, chỉ giữ field hữu ích để giảm token.
// Dùng cho SuggestRecommendations (legacy). SuggestFullReport dùng trực tiếp
// analysis.CoachInput vì đã có cấu trúc gọn ở domain layer.
type promptPayload struct {
	Player        promptPlayer    `json:"player"`
	RecentMatches []recentMatch   `json:"recentMatches"`
	Findings      []findingForLLM `json:"findings"`
}

type promptPlayer struct {
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Region      string `json:"region"`
	PrimaryRole string `json:"primaryRole"`
}

type recentMatch struct {
	Map             string  `json:"map"`
	Agent           string  `json:"agent"`
	Role            string  `json:"role"`
	Kills           int     `json:"kills"`
	Deaths          int     `json:"deaths"`
	Assists         int     `json:"assists"`
	RoundsPlayed    int     `json:"roundsPlayed"`
	HeadshotPercent float64 `json:"headshotPercent"`
	Won             bool    `json:"won"`
}

type findingForLLM struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	Detail     string `json:"detail"`
}

func buildPromptPayload(snapshot analysis.PlayerSnapshot, findings []analysis.Finding) promptPayload {
	const maxMatches = 8

	p := promptPayload{
		Player: promptPlayer{
			Name:        snapshot.Name,
			Tagline:     snapshot.Tagline,
			Region:      snapshot.Region,
			PrimaryRole: snapshot.PrimaryRole,
		},
	}

	matches := snapshot.RecentMatches
	if len(matches) > maxMatches {
		matches = matches[:maxMatches]
	}
	for _, m := range matches {
		p.RecentMatches = append(p.RecentMatches, recentMatch{
			Map:             m.Map,
			Agent:           m.Agent,
			Role:            m.Role,
			Kills:           m.Kills,
			Deaths:          m.Deaths,
			Assists:         m.Assists,
			RoundsPlayed:    m.RoundsPlayed,
			HeadshotPercent: m.HeadshotPercent,
			Won:             m.Won,
		})
	}

	for _, f := range findings {
		p.Findings = append(p.Findings, findingForLLM{
			ID:         f.ID,
			Title:      f.Title,
			Severity:   f.Severity,
			Confidence: f.Confidence,
			Detail:     f.Detail,
		})
	}
	return p
}
