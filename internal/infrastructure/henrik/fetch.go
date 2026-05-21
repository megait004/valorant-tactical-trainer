package henrik

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
	"valorant-tactical-trainer/desktop/internal/infrastructure/env"
)

// SnapshotResult là kết quả 1 lần fetch — Snapshot + metadata (cached, fetched-at,
// source URL) để UI hiển thị badge "from cache" / "fetched at HH:MM".
type SnapshotResult struct {
	Snapshot  analysis.PlayerSnapshot `json:"snapshot"`
	Source    string                  `json:"source"`
	Cached    bool                    `json:"cached"`
	FetchedAt string                  `json:"fetchedAt"`
	Message   string                  `json:"message"`
}

// FetchMatchSnapshot là Wails entry — load match history từ Henrik (qua cache
// nếu có), decode sang PlayerSnapshot domain layer dùng được.
//
// Pipeline:
//  1. Status check (consent + Riot ID).
//  2. Cache lookup theo URL hash (TTL từ settings.CacheTTLMinutes).
//  3. Local rate limit reservation.
//  4. HTTP GET v3/matches/{region}/{name}/{tag}?size=N với Authorization header.
//  5. Handle status code (429, 401/403, non-2xx).
//  6. Decode body → Snapshot, ghi cache.
func (c *Client) FetchMatchSnapshot(ctx context.Context, settings datasettings.DataSettings) (SnapshotResult, error) {
	status := c.Status(settings)
	if !status.CanFetchPersonalData {
		return SnapshotResult{}, errors.New(status.Message)
	}

	requestURL := buildMatchesURL(settings)
	cachePath := c.cachePath(requestURL)
	cacheTTL := time.Duration(settings.CacheTTLMinutes) * time.Minute
	if cached, ok := c.readCache(cachePath, cacheTTL); ok {
		snapshot, err := decodeMatches(settings, cached.Body)
		if err != nil {
			return SnapshotResult{}, err
		}
		return SnapshotResult{
			Snapshot:  snapshot,
			Source:    requestURL,
			Cached:    true,
			FetchedAt: cached.FetchedAt.Format(time.RFC3339),
			Message:   "Dùng cache Henrik local.",
		}, nil
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

	return SnapshotResult{
		Snapshot:  snapshot,
		Source:    requestURL,
		Cached:    false,
		FetchedAt: fetchedAt.Format(time.RFC3339),
		Message:   "Fetch Henrik thành công qua Go core.",
	}, nil
}

// buildMatchesURL build URL v3/matches/{region}/{name}/{tag}?size=N.
// matchCount kẹp trong khoảng 1..20 để tránh response quá lớn.
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

// applyAPIKey gắn header authorization cho Henrik authenticated tier.
// Ưu tiên HENRIK_API_KEY (HDEV-) trong .env. Fallback settings.APIKey (user
// paste qua UI). KHÔNG gửi RGAPI- (key Riot) sang Henrik — nó sẽ 401.
func applyAPIKey(req *http.Request, settings datasettings.DataSettings) {
	key := strings.TrimSpace(env.Load("HENRIK_API_KEY"))
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

// compactBody rút gọn response body để in vào error message — tránh log đầy
// HTML page lớn từ Henrik error.
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
