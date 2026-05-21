// Package henrik là adapter gọi Henrik VALORANT API (api.henrikdev.xyz) —
// dùng làm fallback khi RGAPI dev key của user không có quyền VAL-MATCH-V1
// (case phổ biến với account dev mới).
//
// Tổ chức file:
//
//	client.go    — Client struct + factory (NewClient / NewClientForTest)
//	status.go    — Status DTO + Client.Status() (rate limit / consent guard)
//	fetch.go     — FetchMatchSnapshot (Wails entry) + buildMatchesURL
//	decode.go    — Parse JSON response → analysis.PlayerSnapshot
//	cache.go     — Disk cache theo URL hash, TTL từ settings
//	ratelimit.go — Local rate limiter (sliding window 1 phút)
package henrik

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// baseURL là endpoint Henrik public. Có thể bị override trong test (httptest
// server). Bị giữ ở package-level vì test pattern hiện tại setup/teardown
// global — refactor sang DI sẽ làm ở vòng sau.
var baseURL = "https://api.henrikdev.xyz/valorant"

// maxResponseBytes giới hạn body để tránh DoS từ Henrik (vd v3/matches của
// player chơi quá nhiều round).
const maxResponseBytes int64 = 64 * 1024 * 1024

// Client gọi Henrik API có cache local + rate limit.
type Client struct {
	httpClient *http.Client
	cacheDir   string
	now        func() time.Time

	mu       sync.Mutex
	requests []time.Time
}

// NewClient tạo client production với HTTP timeout 15s và cache dir trong
// UserCacheDir. Fallback về TempDir nếu OS không trả cache dir.
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

// NewClientForTest cho phép inject HTTP client + cache dir + clock — dùng
// trong unit test với httptest.NewServer + t.TempDir.
func NewClientForTest(httpClient *http.Client, cacheDir string, now func() time.Time) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	if now == nil {
		now = time.Now
	}
	return &Client{httpClient: httpClient, cacheDir: cacheDir, now: now}
}
