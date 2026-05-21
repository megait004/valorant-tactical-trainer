package riot

import (
	"encoding/json"
	"strings"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

// File này giữ struct mirror response VAL-MATCH-V1 + logic decode → MatchSummary.
// Tách riêng để match.go chỉ tập trung HTTP/pipeline, decode đụng đâu chỉ sửa
// 1 file.

type matchlistResponse struct {
	PUUID   string                `json:"puuid"`
	History []matchlistHistoryRow `json:"history"`
}

type matchlistHistoryRow struct {
	MatchID         string `json:"matchId"`
	GameStartMillis int64  `json:"gameStartTimeMillis"`
	QueueID         string `json:"queueId"`
}

type matchDetail struct {
	MatchInfo    matchDetailInfo     `json:"matchInfo"`
	Players      []matchDetailPlayer `json:"players"`
	Teams        []matchDetailTeam   `json:"teams"`
	RoundResults []json.RawMessage   `json:"roundResults"`
}

type matchDetailInfo struct {
	MatchID  string `json:"matchId"`
	MapID    string `json:"mapId"`
	IsRanked bool   `json:"isRanked"`
	QueueID  string `json:"queueId"`
}

type matchDetailPlayer struct {
	PUUID           string           `json:"puuid"`
	GameName        string           `json:"gameName"`
	TagLine         string           `json:"tagLine"`
	TeamID          string           `json:"teamId"`
	CharacterID     string           `json:"characterId"`
	Stats           matchDetailStats `json:"stats"`
	CompetitiveTier int              `json:"competitiveTier"`
}

type matchDetailStats struct {
	Score          int `json:"score"`
	RoundsPlayed   int `json:"roundsPlayed"`
	Kills          int `json:"kills"`
	Deaths         int `json:"deaths"`
	Assists        int `json:"assists"`
	PlaytimeMillis int `json:"playtimeMillis"`
	AbilityCasts   any `json:"abilityCasts"`
}

type matchDetailTeam struct {
	TeamID       string `json:"teamId"`
	Won          bool   `json:"won"`
	RoundsPlayed int    `json:"roundsPlayed"`
	RoundsWon    int    `json:"roundsWon"`
}

// decodeMatchDetail convert VAL-MATCH-V1 detail → analysis.MatchSummary.
// Logic:
//   - Tìm player có PUUID khớp.
//   - Tính rounds = stats.RoundsPlayed → fallback roundResults length →
//     fallback max teams.RoundsPlayed → fallback 1.
//   - Xác định Won bằng cách match TeamID với teams[].Won.
//   - Map agent CharacterID UUID → tên (catalog.go) → role (roleForAgent).
//   - HeadshotPercent = 0 vì VAL-MATCH-V1 không trả body/head/leg trực tiếp
//     (cần parse từ roundResults, sẽ làm khi cần độ chính xác cao hơn).
func decodeMatchDetail(body []byte, puuid string) (analysis.MatchSummary, bool) {
	var detail matchDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return analysis.MatchSummary{}, false
	}

	var me *matchDetailPlayer
	puuidLower := strings.ToLower(strings.TrimSpace(puuid))
	for i := range detail.Players {
		if strings.ToLower(detail.Players[i].PUUID) == puuidLower {
			me = &detail.Players[i]
			break
		}
	}
	if me == nil {
		return analysis.MatchSummary{}, false
	}

	rounds := me.Stats.RoundsPlayed
	if rounds == 0 {
		rounds = len(detail.RoundResults)
	}
	if rounds == 0 {
		for _, t := range detail.Teams {
			if t.RoundsPlayed > rounds {
				rounds = t.RoundsPlayed
			}
		}
	}
	if rounds == 0 {
		rounds = 1
	}

	won := false
	for _, t := range detail.Teams {
		if t.TeamID == me.TeamID && t.Won {
			won = true
			break
		}
	}

	agent := agentNameFromUUID(me.CharacterID)
	mapName := mapNameFromUUID(detail.MatchInfo.MapID)

	return analysis.MatchSummary{
		ID:              detail.MatchInfo.MatchID,
		Map:             mapName,
		Agent:           agent,
		Role:            roleForAgent(agent),
		Kills:           me.Stats.Kills,
		Deaths:          me.Stats.Deaths,
		Assists:         me.Stats.Assists,
		RoundsPlayed:    rounds,
		FirstBloods:     0,
		FirstDeaths:     0,
		HeadshotPercent: 0,
		Won:             won,
	}, true
}

// roleForAgent giữ đồng bộ với henrik.roleForAgent để output nhất quán.
// Khi cập nhật agent mới, sửa cả 2 nơi.
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
