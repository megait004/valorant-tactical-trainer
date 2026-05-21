package henrik

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
)

// File này giữ các struct mirror response của Henrik API + logic decode về
// analysis.PlayerSnapshot. Tách riêng để client.go không phải import
// encoding/json và mở rộng struct mới (vd thêm round data) chỉ đụng 1 file.

type matchesResponse struct {
	Data []matchData `json:"data"`
}

type matchData struct {
	Metadata matchMetadata `json:"metadata"`
	Players  matchPlayers  `json:"players"`
	Teams    matchTeams    `json:"teams"`
	Rounds   []matchRound  `json:"rounds"`
}

type matchMetadata struct {
	MatchID string `json:"matchid"`
	Map     string `json:"map"`
}

type matchPlayers struct {
	AllPlayers []matchPlayer `json:"all_players"`
}

type matchPlayer struct {
	Name      string      `json:"name"`
	Tag       string      `json:"tag"`
	Team      string      `json:"team"`
	Character string      `json:"character"`
	Stats     playerStats `json:"stats"`
}

type playerStats struct {
	Kills     int `json:"kills"`
	Deaths    int `json:"deaths"`
	Assists   int `json:"assists"`
	Headshots int `json:"headshots"`
	Bodyshots int `json:"bodyshots"`
	Legshots  int `json:"legshots"`
}

type matchTeams struct {
	Red  teamStats `json:"red"`
	Blue teamStats `json:"blue"`
}

type teamStats struct {
	HasWon    bool `json:"has_won"`
	RoundsWon int  `json:"rounds_won"`
}

type matchRound struct{}

// decodeMatches convert Henrik JSON body → PlayerSnapshot. Logic:
//   - Tìm player trùng tên+tag trong all_players (Henrik trả full team data).
//   - Tính rounds = len(rounds) → fallback team.RoundsWon × 2 → fallback 1.
//   - Tính headshot % từ headshots/bodyshots/legshots.
//   - PrimaryRole lấy từ trận đầu (UI dùng làm gợi ý — không bị aggregate).
func decodeMatches(settings datasettings.DataSettings, body []byte) (analysis.PlayerSnapshot, error) {
	var response matchesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return analysis.PlayerSnapshot{}, fmt.Errorf("err parsing Henrik matches: %w", err)
	}
	if len(response.Data) == 0 {
		return analysis.PlayerSnapshot{}, errors.New("Henrik không trả match history đủ để phân tích")
	}

	snapshot := analysis.PlayerSnapshot{
		Name:        settings.RiotName,
		Tagline:     settings.RiotTag,
		Region:      settings.Region,
		PrimaryRole: "unknown",
	}

	for index, match := range response.Data {
		player, ok := findPlayer(match.Players.AllPlayers, settings)
		if !ok {
			continue
		}

		matchID := match.Metadata.MatchID
		if matchID == "" {
			matchID = fmt.Sprintf("henrik-%d", index+1)
		}

		rounds := len(match.Rounds)
		if rounds == 0 {
			rounds = match.Teams.Red.RoundsWon + match.Teams.Blue.RoundsWon
		}
		if rounds == 0 {
			rounds = 1
		}

		snapshot.RecentMatches = append(snapshot.RecentMatches, analysis.MatchSummary{
			ID:              matchID,
			Map:             match.Metadata.Map,
			Agent:           player.Character,
			Role:            roleForAgent(player.Character),
			Kills:           player.Stats.Kills,
			Deaths:          player.Stats.Deaths,
			Assists:         player.Stats.Assists,
			RoundsPlayed:    rounds,
			FirstBloods:     0,
			FirstDeaths:     0,
			HeadshotPercent: headshotPercent(player.Stats),
			Won:             teamWon(match.Teams, player.Team),
		})
	}

	if len(snapshot.RecentMatches) == 0 {
		return analysis.PlayerSnapshot{}, errors.New("không tìm thấy player trong match history Henrik")
	}
	snapshot.PrimaryRole = snapshot.RecentMatches[0].Role
	return snapshot, nil
}

func findPlayer(players []matchPlayer, settings datasettings.DataSettings) (matchPlayer, bool) {
	name := strings.ToLower(settings.RiotName)
	tag := strings.ToLower(settings.RiotTag)
	for _, player := range players {
		if strings.ToLower(player.Name) == name && strings.ToLower(player.Tag) == tag {
			return player, true
		}
	}
	return matchPlayer{}, false
}

func teamWon(teams matchTeams, team string) bool {
	switch strings.ToLower(team) {
	case "red":
		return teams.Red.HasWon
	case "blue":
		return teams.Blue.HasWon
	default:
		return false
	}
}

func headshotPercent(stats playerStats) float64 {
	total := stats.Headshots + stats.Bodyshots + stats.Legshots
	if total == 0 {
		return 0
	}
	return float64(stats.Headshots) / float64(total) * 100
}

// roleForAgent map tên agent (Henrik trả) → vai trò. Giữ đồng bộ với
// riot.roleForAgent để output nhất quán giữa 2 nguồn.
func roleForAgent(agent string) string {
	switch strings.ToLower(agent) {
	case "jett", "raze", "reyna", "neon", "phoenix", "yoru", "iso", "waylay":
		return "duelist"
	case "omen", "brimstone", "viper", "astra", "harbor", "clove":
		return "controller"
	case "sova", "fade", "breach", "skye", "kayo", "kay/o", "gekko", "tejo":
		return "initiator"
	case "cypher", "killjoy", "sage", "chamber", "deadlock", "vyse":
		return "sentinel"
	default:
		return "unknown"
	}
}
