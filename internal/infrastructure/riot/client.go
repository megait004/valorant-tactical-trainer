// Package riot là adapter cho Riot Games API chính thức — gồm 2 client tách
// rời nhưng dùng chung helper:
//
//   - Client      (account.go): Riot Account-V1 cho login flow.
//   - MatchClient (match.go):   VAL-MATCH-V1 cho match history.
//
// Cấu hình:
//
//   - RGAPI dev key đọc từ env `RIOT_API_KEY` (thường ở file .env root).
//   - Không bao giờ ghi key vào settings.json — header X-Riot-Token chỉ
//     được set bởi package này.
//
// Tổ chức file:
//
//	client.go    — Client struct (Account-V1) + session state + helpers
//	account.go   — Login() + DTOs (PlayerInfo, LoginResult, SettingsSink)
//	routing.go   — resolveRoutingShard + defaultMatchHost (region map)
//	catalog.go   — UUID → name lookup (agent, map)
//	match.go     — MatchClient + FetchMatchSnapshot
//	match_decode.go — Parse VAL-MATCH-V1 detail → analysis.MatchSummary
//	match_cache.go  — Disk cache + sliding window rate limit
package riot

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// session lưu state auth in-memory cho Client. Không persist — restart app
// thì user phải login lại (đã có guard ở UI dùng IsLoggedIn).
type session struct {
	apiKey     string
	playerInfo *PlayerInfo
}

// Client là Riot Account-V1 client (login flow). Match history dùng
// MatchClient riêng để tách concern (mỗi client có rate limit / cache riêng).
type Client struct {
	mu         sync.Mutex
	httpClient *http.Client
	sess       *session
	sink       SettingsSink
}

// NewClient tạo client với HTTP timeout 12s.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 12 * time.Second},
	}
}

// SetSettingsSink đăng ký bridge để Login/Logout đồng bộ session vào store.
// Gọi 1 lần khi khởi tạo Services. Tách interface để tránh import cycle với
// infrastructure/store.
func (c *Client) SetSettingsSink(sink SettingsSink) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sink = sink
}

// IsLoggedIn kiểm tra session có hợp lệ chưa.
func (c *Client) IsLoggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sess != nil && c.sess.playerInfo != nil
}

// GetPlayerInfo trả về player info từ session hiện tại (nil nếu chưa login).
func (c *Client) GetPlayerInfo() *PlayerInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sess == nil {
		return nil
	}
	return c.sess.playerInfo
}

// GetAPIKey trả về API key đang dùng (rỗng nếu chưa login). Hiện chưa có
// caller nào dùng — giữ làm escape hatch cho code mở rộng (vd debug panel).
func (c *Client) GetAPIKey() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sess == nil {
		return ""
	}
	return c.sess.apiKey
}

// Logout clear session + notify sink (sink sẽ xoá PUUID/RiotID khỏi
// settings.json + xoá last_report.json để tránh stale cache).
func (c *Client) Logout() {
	c.mu.Lock()
	sink := c.sink
	c.sess = nil
	c.mu.Unlock()
	if sink != nil {
		_ = sink.OnLogout()
	}
}

// compactBody rút gọn body lỗi để in vào error message (tránh log full HTML).
func compactBody(body []byte) string {
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
