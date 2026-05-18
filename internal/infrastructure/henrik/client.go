package henrik

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
	"valorant-tactical-trainer/desktop/internal/infrastructure/riot"
)

var baseURL = "https://api.henrikdev.xyz/valorant"

const maxResponseBytes int64 = 64 * 1024 * 1024

type Status struct {
	BaseURL              string `json:"baseURL"`
	ConsentGranted       bool   `json:"consentGranted"`
	CanFetchPersonalData bool   `json:"canFetchPersonalData"`
	RateLimitPerMinute   int    `json:"rateLimitPerMinute"`
	CacheTTLMinutes      int    `json:"cacheTTLMinutes"`
	SafeMode             bool   `json:"safeMode"`
	Message              string `json:"message"`
	NextStep             string `json:"nextStep"`
}

type Client struct {
	httpClient *http.Client
	cacheDir   string
	now        func() time.Time
	mu         sync.Mutex
	requests   []time.Time
}

func NewClient() *Client {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		cacheDir:   filepath.Join(cacheDir, "Valorant Tactical Trainer", "henrik"),
		now:        time.Now,
	}
}

func NewClientForTest(httpClient *http.Client, cacheDir string, now func() time.Time) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	if now == nil {
		now = time.Now
	}
	return &Client{httpClient: httpClient, cacheDir: cacheDir, now: now}
}

func (c *Client) Status(settings datasettings.DataSettings) Status {
	envHenrikKey := strings.TrimSpace(riot.LoadEnvKey("HENRIK_API_KEY"))
	hasHenrikKey := envHenrikKey != "" || (settings.APIKey != "" && !strings.HasPrefix(settings.APIKey, "RGAPI-"))

	limit := 30
	if settings.RateLimitTier == "enhanced" || hasHenrikKey {
		limit = 90
	}

	status := Status{
		BaseURL:            baseURL,
		ConsentGranted:     settings.ConsentPersonalData,
		RateLimitPerMinute: limit,
		CacheTTLMinutes:    settings.CacheTTLMinutes,
		SafeMode:           true,
	}

	if !settings.ConsentPersonalData {
		status.Message = "Chưa có consent nên app không fetch dữ liệu cá nhân."
		status.NextStep = "Bật consent và nhập Riot ID trước khi kết nối API thật."
		return status
	}

	if settings.RiotName == "" || settings.RiotTag == "" {
		status.Message = "Đã có consent nhưng thiếu Riot ID hoặc tag."
		status.NextStep = "Nhập Riot name và tag để chuẩn bị fetch match history."
		return status
	}

	status.CanFetchPersonalData = true
	if hasHenrikKey {
		status.Message = "Đã có HDEV key, sẵn sàng fetch Henrik authenticated (90req/min)."
		status.NextStep = "Bấm 'Fetch report' để lấy match history."
	} else {
		status.Message = "Sẵn sàng fetch Henrik public (30req/min). Paste HENRIK_API_KEY vào .env để nâng lên 90req/min."
		status.NextStep = "Bấm 'Fetch report' để lấy match history."
	}
	return status
}

type SnapshotResult struct {
	Snapshot  analysis.PlayerSnapshot `json:"snapshot"`
	Source    string                  `json:"source"`
	Cached    bool                    `json:"cached"`
	FetchedAt string                  `json:"fetchedAt"`
	Message   string                  `json:"message"`
}

type cacheEnvelope struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Body      json.RawMessage `json:"body"`
}

func (c *Client) FetchMatchSnapshot(ctx context.Context, settings datasettings.DataSettings) (SnapshotResult, error) {
	status := c.Status(settings)
	if !status.CanFetchPersonalData {
		return SnapshotResult{}, errors.New(status.Message)
	}

	requestURL := buildMatchesURL(settings)
	cachePath := c.cachePath(requestURL)
	if cached, ok := c.readCache(cachePath, time.Duration(settings.CacheTTLMinutes)*time.Minute); ok {
		snapshot, err := decodeMatches(settings, cached.Body)
		if err != nil {
			return SnapshotResult{}, err
		}
		return SnapshotResult{Snapshot: snapshot, Source: requestURL, Cached: true, FetchedAt: cached.FetchedAt.Format(time.RFC3339), Message: "Dùng cache Henrik local."}, nil
	}

	if err := c.reserveRequest(status.RateLimitPerMinute); err != nil {
		return SnapshotResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return SnapshotResult{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ValorantTacticalTrainer/0.1")
	applyAPIKey(req, settings)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("err fetching Henrik matches: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBytes+1))
	if err != nil {
		return SnapshotResult{}, fmt.Errorf("err reading Henrik response: %w", err)
	}
	if int64(len(body)) > maxResponseBytes {
		return SnapshotResult{}, fmt.Errorf("Henrik response quá lớn trên %dMB, giảm số match hoặc dùng endpoint nhẹ hơn", maxResponseBytes/(1024*1024))
	}

	if res.StatusCode == http.StatusTooManyRequests {
		return SnapshotResult{}, errors.New("Henrik rate limit 429, thử lại sau hoặc giảm thao tác fetch")
	}
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return SnapshotResult{}, fmt.Errorf("Henrik auth status %d: kiểm tra API key, header hoặc format key. Response: %s", res.StatusCode, compactBody(body))
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return SnapshotResult{}, fmt.Errorf("Henrik API trả status %d: %s", res.StatusCode, compactBody(body))
	}

	snapshot, err := decodeMatches(settings, body)
	if err != nil {
		return SnapshotResult{}, err
	}
	fetchedAt := c.now()
	c.writeCache(cachePath, cacheEnvelope{FetchedAt: fetchedAt, Body: body})

	return SnapshotResult{Snapshot: snapshot, Source: requestURL, Cached: false, FetchedAt: fetchedAt.Format(time.RFC3339), Message: "Fetch Henrik thành công qua Go core."}, nil
}

func applyAPIKey(req *http.Request, settings datasettings.DataSettings) {
	// Ưu tiên HENRIK_API_KEY (HDEV-) trong .env. Nếu không có thì thử settings.APIKey
	// (user paste qua UI). Tuyệt đối không gửi RGAPI- (key Riot) sang Henrik — nó sẽ 401.
	key := strings.TrimSpace(riot.LoadEnvKey("HENRIK_API_KEY"))
	if key == "" {
		key = strings.TrimSpace(settings.APIKey)
	}
	if key == "" || strings.HasPrefix(key, "RGAPI-") {
		return
	}
	header := settings.APIKeyHeader
	if header == "" || strings.EqualFold(header, "X-Riot-Token") {
		header = "Authorization"
	}
	req.Header.Set(header, key)
}

func compactBody(body []byte) string {
	message := strings.TrimSpace(string(body))
	if message == "" {
		return "empty body"
	}
	message = strings.Join(strings.Fields(message), " ")
	if len(message) > 500 {
		return message[:500] + "..."
	}
	return message
}

func buildMatchesURL(settings datasettings.DataSettings) string {
	region := url.PathEscape(strings.ToLower(strings.TrimSpace(settings.Region)))
	name := url.PathEscape(strings.TrimSpace(settings.RiotName))
	tag := url.PathEscape(strings.TrimSpace(settings.RiotTag))
	matchCount := settings.MatchCount
	if matchCount <= 0 {
		matchCount = 5
	}
	if matchCount > 20 {
		matchCount = 20
	}
	return fmt.Sprintf("%s/v3/matches/%s/%s/%s?size=%d", baseURL, region, name, tag, matchCount)
}

func (c *Client) reserveRequest(limit int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	windowStart := now.Add(-1 * time.Minute)
	kept := c.requests[:0]
	for _, item := range c.requests {
		if item.After(windowStart) {
			kept = append(kept, item)
		}
	}
	c.requests = kept

	if len(c.requests) >= limit {
		return fmt.Errorf("local rate limit %d/min đã đầy, thử lại sau", limit)
	}
	c.requests = append(c.requests, now)
	return nil
}

func (c *Client) cachePath(requestURL string) string {
	hash := sha256.Sum256([]byte(requestURL))
	return filepath.Join(c.cacheDir, hex.EncodeToString(hash[:])+".json")
}

func (c *Client) readCache(path string, ttl time.Duration) (cacheEnvelope, bool) {
	if ttl <= 0 {
		return cacheEnvelope{}, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cacheEnvelope{}, false
	}

	var envelope cacheEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return cacheEnvelope{}, false
	}
	if envelope.FetchedAt.IsZero() || c.now().Sub(envelope.FetchedAt) > ttl {
		return cacheEnvelope{}, false
	}
	return envelope, true
}

func (c *Client) writeCache(path string, envelope cacheEnvelope) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

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
