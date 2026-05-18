package riot

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
)

// MatchClient gọi VAL-MATCH-V1 (Riot chính thức) để lấy match history + match details
// theo PUUID. Dùng cùng RGAPI key đã load từ .env.
type MatchClient struct {
	httpClient *http.Client
	cacheDir   string
	now        func() time.Time
	mu         sync.Mutex
	requests   []time.Time
	baseHost   func(shard string) string
}

func NewMatchClient() *MatchClient {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return &MatchClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		cacheDir:   filepath.Join(cacheDir, "Valorant Tactical Trainer", "riot-val"),
		now:        time.Now,
		baseHost:   defaultMatchHost,
	}
}

// MatchSnapshotResult là kết quả fetch match list + decode → PlayerSnapshot
type MatchSnapshotResult struct {
	Snapshot  analysis.PlayerSnapshot `json:"snapshot"`
	Source    string                  `json:"source"`
	Cached    bool                    `json:"cached"`
	FetchedAt string                  `json:"fetchedAt"`
	Message   string                  `json:"message"`
}

// CanFetch kiểm tra điều kiện cần thiết
func (c *MatchClient) CanFetch(settings datasettings.DataSettings) error {
	if !settings.ConsentPersonalData {
		return errors.New("Chưa có consent — login lại để bật consent")
	}
	if strings.TrimSpace(settings.PUUID) == "" {
		return errors.New("Thiếu PUUID — đăng nhập lại để lấy PUUID từ Riot")
	}
	if strings.TrimSpace(loadAPIKey()) == "" {
		return errors.New("Thiếu RIOT_API_KEY trong .env — paste key vào root project")
	}
	return nil
}

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

// FetchMatchSnapshot gọi VAL-MATCH-V1:
//
//	GET /val/match/v1/matchlists/by-puuid/{puuid}
//	→ N matches mới nhất (lấy theo MatchCount)
//	GET /val/match/v1/matches/{matchId} cho từng match
//	→ map về PlayerSnapshot.RecentMatches
//
// Có cache file theo URL hash (TTL = settings.CacheTTLMinutes).
func (c *MatchClient) FetchMatchSnapshot(ctx context.Context, settings datasettings.DataSettings) (MatchSnapshotResult, error) {
	if err := c.CanFetch(settings); err != nil {
		return MatchSnapshotResult{}, err
	}

	apiKey := loadAPIKey()
	host := c.baseHost(settings.Shard)
	puuid := strings.TrimSpace(settings.PUUID)

	listURL := fmt.Sprintf("https://%s/val/match/v1/matchlists/by-puuid/%s", host, puuid)
	listCachePath := c.cachePath(listURL)
	cacheTTL := time.Duration(settings.CacheTTLMinutes) * time.Minute

	var listBody []byte
	cachedList := false
	if env, ok := c.readCache(listCachePath, cacheTTL); ok {
		listBody = env.Body
		cachedList = true
	} else {
		body, err := c.doGet(ctx, listURL, apiKey)
		if err != nil {
			return MatchSnapshotResult{}, err
		}
		listBody = body
		c.writeCache(listCachePath, cacheEnvelope{FetchedAt: c.now(), Body: body})
	}

	var list matchlistResponse
	if err := json.Unmarshal(listBody, &list); err != nil {
		return MatchSnapshotResult{}, fmt.Errorf("err parse matchlist: %w", err)
	}
	if len(list.History) == 0 {
		return MatchSnapshotResult{}, errors.New("Riot trả 0 match — chưa có trận nào trong VAL-MATCH-V1 cho PUUID này")
	}

	limit := settings.MatchCount
	if limit <= 0 {
		limit = 5
	}
	if limit > len(list.History) {
		limit = len(list.History)
	}

	snapshot := analysis.PlayerSnapshot{
		Name:        settings.RiotName,
		Tagline:     settings.RiotTag,
		Region:      settings.Region,
		PrimaryRole: "unknown",
	}
	allCached := cachedList
	for _, row := range list.History[:limit] {
		matchURL := fmt.Sprintf("https://%s/val/match/v1/matches/%s", host, row.MatchID)
		matchCachePath := c.cachePath(matchURL)

		var detailBody []byte
		if env, ok := c.readCache(matchCachePath, cacheTTL); ok {
			detailBody = env.Body
		} else {
			body, err := c.doGet(ctx, matchURL, apiKey)
			if err != nil {
				return MatchSnapshotResult{}, err
			}
			detailBody = body
			c.writeCache(matchCachePath, cacheEnvelope{FetchedAt: c.now(), Body: body})
			allCached = false
		}

		summary, ok := decodeMatchDetail(detailBody, puuid)
		if !ok {
			continue
		}
		snapshot.RecentMatches = append(snapshot.RecentMatches, summary)
	}

	if len(snapshot.RecentMatches) == 0 {
		return MatchSnapshotResult{}, errors.New("Không decode được match nào — kiểm tra RGAPI key có quyền VAL-MATCH-V1 không")
	}
	snapshot.PrimaryRole = snapshot.RecentMatches[0].Role

	msg := "Fetch Riot VAL-MATCH-V1 OK."
	if allCached {
		msg = "Dùng cache Riot local."
	}
	return MatchSnapshotResult{
		Snapshot:  snapshot,
		Source:    listURL,
		Cached:    allCached,
		FetchedAt: c.now().Format(time.RFC3339),
		Message:   msg,
	}, nil
}

func (c *MatchClient) doGet(ctx context.Context, url, apiKey string) ([]byte, error) {
	if err := c.reserveRequest(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Riot-Token", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ValorantTacticalTrainer/0.1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("err call Riot: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<20))

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.StatusUnauthorized:
		return nil, errors.New("Riot 401 — RGAPI key sai hoặc hết hạn")
	case http.StatusForbidden:
		return nil, errors.New("Riot 403 — RGAPI dev key không có quyền VAL-MATCH-V1 (cần production key) → app sẽ fallback Henrik")
	case http.StatusNotFound:
		return nil, errors.New("Riot 404 — match/PUUID không tồn tại trên shard này")
	case http.StatusTooManyRequests:
		return nil, errors.New("Riot 429 rate limit — thử lại sau")
	default:
		return nil, fmt.Errorf("Riot status %d: %s", resp.StatusCode, compactBodyMatch(body))
	}
}

func (c *MatchClient) reserveRequest() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	windowStart := now.Add(-1 * time.Minute)
	kept := c.requests[:0]
	for _, t := range c.requests {
		if t.After(windowStart) {
			kept = append(kept, t)
		}
	}
	c.requests = kept

	// Riot dev key: 20req/sec, 100req/2min. Cứ giới hạn 50/phút cho an toàn.
	if len(c.requests) >= 50 {
		return errors.New("local rate limit 50/min đã đầy, thử lại sau")
	}
	c.requests = append(c.requests, now)
	return nil
}

type cacheEnvelope struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Body      json.RawMessage `json:"body"`
}

func (c *MatchClient) cachePath(url string) string {
	hash := sha256.Sum256([]byte(url))
	return filepath.Join(c.cacheDir, hex.EncodeToString(hash[:])+".json")
}

func (c *MatchClient) readCache(path string, ttl time.Duration) (cacheEnvelope, bool) {
	if ttl <= 0 {
		return cacheEnvelope{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cacheEnvelope{}, false
	}
	var env cacheEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return cacheEnvelope{}, false
	}
	if env.FetchedAt.IsZero() || c.now().Sub(env.FetchedAt) > ttl {
		return cacheEnvelope{}, false
	}
	return env, true
}

func (c *MatchClient) writeCache(path string, env cacheEnvelope) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

// defaultMatchHost map shard → VAL platform host
// (xem developer.riotgames.com/apis#val-match-v1)
func defaultMatchHost(shard string) string {
	switch strings.ToLower(strings.TrimSpace(shard)) {
	case "na", "br", "latam":
		return "na.api.riotgames.com"
	case "eu":
		return "eu.api.riotgames.com"
	case "kr":
		return "kr.api.riotgames.com"
	case "ap":
		return "ap.api.riotgames.com"
	default:
		return "ap.api.riotgames.com"
	}
}

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
		HeadshotPercent: 0, // VAL-MATCH-V1 không trả body/leg/headshot trực tiếp; cần parse từ roundResults nếu cần.
		Won:             won,
	}, true
}

func compactBodyMatch(body []byte) string {
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		return "empty body"
	}
	msg = strings.Join(strings.Fields(msg), " ")
	if len(msg) > 300 {
		return msg[:300] + "..."
	}
	return msg
}

// roleForAgent giữ giống henrik client để output nhất quán.
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
