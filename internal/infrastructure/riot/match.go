package riot

import (
	"context"
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
	"valorant-tactical-trainer/desktop/internal/infrastructure/env"
)

// MatchClient gọi VAL-MATCH-V1 (Riot chính thức) để lấy match list +
// match detail theo PUUID. Dùng cùng RGAPI key trong .env như Client.
//
// Note: VAL-MATCH-V1 cần production-tier key. Dev key thường bị 403 → app
// fallback Henrik (xem AnalysisService.fetchSnapshot).
type MatchClient struct {
	httpClient *http.Client
	cacheDir   string
	now        func() time.Time

	mu       sync.Mutex
	requests []time.Time

	// baseHost cho phép inject host trong test (dù chưa có test riêng).
	baseHost func(shard string) string
}

// NewMatchClient tạo client production với HTTP timeout 15s và cache dir
// riêng (tránh trộn với Henrik cache).
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

// MatchSnapshotResult là kết quả 1 lần fetch — giống henrik.SnapshotResult
// nhưng giữ tên Riot riêng để code caller dễ phân biệt nguồn dữ liệu.
type MatchSnapshotResult struct {
	Snapshot  analysis.PlayerSnapshot `json:"snapshot"`
	Source    string                  `json:"source"`
	Cached    bool                    `json:"cached"`
	FetchedAt string                  `json:"fetchedAt"`
	Message   string                  `json:"message"`
}

// CanFetch validate các điều kiện cần thiết: consent, PUUID, API key.
func (c *MatchClient) CanFetch(settings datasettings.DataSettings) error {
	if !settings.ConsentPersonalData {
		return errors.New("Chưa có consent — login lại để bật consent")
	}
	if strings.TrimSpace(settings.PUUID) == "" {
		return errors.New("Thiếu PUUID — đăng nhập lại để lấy PUUID từ Riot")
	}
	if strings.TrimSpace(env.Load("RIOT_API_KEY")) == "" {
		return errors.New("Thiếu RIOT_API_KEY trong .env — paste key vào root project")
	}
	return nil
}

// FetchMatchSnapshot gọi VAL-MATCH-V1 theo pipeline:
//
//	GET /val/match/v1/matchlists/by-puuid/{puuid}
//	→ N matches mới nhất (lấy theo settings.MatchCount, max 20)
//	GET /val/match/v1/matches/{matchId} cho từng match
//	→ map về PlayerSnapshot.RecentMatches
//
// Có cache file theo URL hash (TTL = settings.CacheTTLMinutes).
func (c *MatchClient) FetchMatchSnapshot(ctx context.Context, settings datasettings.DataSettings) (MatchSnapshotResult, error) {
	if err := c.CanFetch(settings); err != nil {
		return MatchSnapshotResult{}, err
	}

	apiKey := env.Load("RIOT_API_KEY")
	host := c.baseHost(settings.Shard)
	puuid := strings.TrimSpace(settings.PUUID)
	cacheTTL := time.Duration(settings.CacheTTLMinutes) * time.Minute

	listURL := fmt.Sprintf("https://%s/val/match/v1/matchlists/by-puuid/%s", host, puuid)
	listBody, cachedList, err := c.loadOrFetch(ctx, listURL, apiKey, cacheTTL)
	if err != nil {
		return MatchSnapshotResult{}, err
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
		detailBody, detailCached, err := c.loadOrFetch(ctx, matchURL, apiKey, cacheTTL)
		if err != nil {
			return MatchSnapshotResult{}, err
		}
		if !detailCached {
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

// loadOrFetch try cache trước, fallback HTTP GET nếu miss/expire. Trả về body
// + flag cached để caller biết overall trạng thái snapshot.
func (c *MatchClient) loadOrFetch(ctx context.Context, url, apiKey string, ttl time.Duration) (body []byte, cached bool, err error) {
	path := c.cachePath(url)
	if envlp, ok := c.readCache(path, ttl); ok {
		return envlp.Body, true, nil
	}
	body, err = c.doGet(ctx, url, apiKey)
	if err != nil {
		return nil, false, err
	}
	c.writeCache(path, cacheEnvelope{FetchedAt: c.now(), Body: body})
	return body, false, nil
}

// doGet là HTTP GET helper với rate limit + status code handling.
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
		return nil, fmt.Errorf("Riot status %d: %s", resp.StatusCode, compactBody(body))
	}
}
