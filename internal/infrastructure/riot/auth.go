package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// PlayerInfo là thông tin người chơi sau khi đăng nhập thành công
type PlayerInfo struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Region   string `json:"region"`
	Shard    string `json:"shard"`
}

// LoginResult trả về sau khi gọi Login()
type LoginResult struct {
	Success    bool        `json:"success"`
	Error      string      `json:"error,omitempty"`
	PlayerInfo *PlayerInfo `json:"playerInfo,omitempty"`
}

// session lưu trạng thái auth trong memory
type session struct {
	apiKey     string
	playerInfo *PlayerInfo
}

// SettingsSink là interface để Login/Logout đồng bộ Riot ID + PUUID vào
// settings.json và xóa report cache cũ khi đổi nick. Tách ra interface để
// tránh import cycle với infrastructure/localstore.
type SettingsSink interface {
	OnLogin(player PlayerInfo) error
	OnLogout() error
}

// Client là Riot API client cho Account-V1 auth
type Client struct {
	mu         sync.Mutex
	httpClient *http.Client
	sess       *session
	sink       SettingsSink
}

// NewClient tạo Riot client mới
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 12 * time.Second},
	}
}

// SetSettingsSink đăng ký sink để Login/Logout cập nhật persistent state.
// Gọi 1 lần khi khởi tạo Services.
func (c *Client) SetSettingsSink(sink SettingsSink) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sink = sink
}

// IsLoggedIn kiểm tra đã đăng nhập chưa
func (c *Client) IsLoggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sess != nil && c.sess.playerInfo != nil
}

// GetPlayerInfo trả về thông tin người chơi hiện tại
func (c *Client) GetPlayerInfo() *PlayerInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sess == nil {
		return nil
	}
	return c.sess.playerInfo
}

// GetAPIKey trả về API key đang dùng (rỗng nếu chưa login)
func (c *Client) GetAPIKey() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sess == nil {
		return ""
	}
	return c.sess.apiKey
}

// Logout xóa session hiện tại
func (c *Client) Logout() {
	c.mu.Lock()
	sink := c.sink
	c.sess = nil
	c.mu.Unlock()
	if sink != nil {
		_ = sink.OnLogout()
	}
}

// Login xác thực Riot ID qua Account-V1 (Riot API chính thức).
// API key được load từ biến môi trường RIOT_API_KEY hoặc file `.env` ở root project,
// không yêu cầu user nhập trên UI.
//
// Đây là cách duy nhất hợp lệ với một dev API key — RGAPI key KHÔNG xác thực
// được password của user (cần OAuth client production để làm RSO thật).
func (c *Client) Login(riotID, tagLine, region string) (LoginResult, error) {
	riotID = strings.TrimSpace(riotID)
	tagLine = strings.TrimSpace(strings.TrimPrefix(tagLine, "#"))
	region = strings.ToLower(strings.TrimSpace(region))

	if riotID == "" || tagLine == "" {
		return LoginResult{Error: "Nhập Riot ID và tag (vd: giaphue#DATN)"}, nil
	}
	apiKey := loadAPIKey()
	if apiKey == "" {
		return LoginResult{Error: "Chưa có RIOT_API_KEY — paste key vào file .env ở root project"}, nil
	}
	if !strings.HasPrefix(apiKey, "RGAPI-") {
		return LoginResult{Error: "RIOT_API_KEY không hợp lệ — phải bắt đầu bằng 'RGAPI-'"}, nil
	}
	if region == "" {
		region = "ap"
	}

	routing, shard := resolveRoutingShard(region)

	endpoint := fmt.Sprintf(
		"https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		routing,
		url.PathEscape(riotID),
		url.PathEscape(tagLine),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return LoginResult{Error: fmt.Sprintf("err build request: %v", err)}, nil
	}
	req.Header.Set("X-Riot-Token", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ValorantTacticalTrainer/0.1")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LoginResult{Error: fmt.Sprintf("Không kết nối được Riot API: %v", err)}, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

	switch resp.StatusCode {
	case http.StatusOK:
		// fallthrough
	case http.StatusUnauthorized:
		return LoginResult{Error: "API key không hợp lệ hoặc đã hết hạn (dev key hết hạn sau 24h, lấy lại tại developer.riotgames.com)"}, nil
	case http.StatusForbidden:
		return LoginResult{Error: "API key bị chặn — kiểm tra account dev hoặc xin lại key"}, nil
	case http.StatusNotFound:
		return LoginResult{Error: fmt.Sprintf("Không tìm thấy Riot ID '%s#%s' ở region '%s'. Kiểm tra chính tả và region.", riotID, tagLine, region)}, nil
	case http.StatusTooManyRequests:
		return LoginResult{Error: "Riot API rate limit (429), thử lại sau ít giây"}, nil
	default:
		return LoginResult{Error: fmt.Sprintf("Riot API trả status %d: %s", resp.StatusCode, compactBody(body))}, nil
	}

	var account struct {
		PUUID    string `json:"puuid"`
		GameName string `json:"gameName"`
		TagLine  string `json:"tagLine"`
	}
	if err := json.Unmarshal(body, &account); err != nil {
		return LoginResult{Error: fmt.Sprintf("err parse account response: %v", err)}, nil
	}
	if account.PUUID == "" {
		return LoginResult{Error: "Riot API trả response không có PUUID"}, nil
	}

	player := &PlayerInfo{
		PUUID:    account.PUUID,
		GameName: account.GameName,
		TagLine:  account.TagLine,
		Region:   region,
		Shard:    shard,
	}

	c.mu.Lock()
	c.sess = &session{apiKey: apiKey, playerInfo: player}
	sink := c.sink
	c.mu.Unlock()

	if sink != nil {
		if err := sink.OnLogin(*player); err != nil {
			return LoginResult{Error: fmt.Sprintf("Login OK nhưng lưu settings lỗi: %v", err)}, nil
		}
	}

	return LoginResult{Success: true, PlayerInfo: player}, nil
}

func compactBody(body []byte) string {
	message := strings.TrimSpace(string(body))
	if message == "" {
		return "empty body"
	}
	message = strings.Join(strings.Fields(message), " ")
	if len(message) > 300 {
		return message[:300] + "..."
	}
	return message
}

// resolveRoutingShard map region (ap, eu, na, kr, br, latam) → routing cluster
// (americas/europe/asia) cho Account-V1 và shard cho Valorant content.
func resolveRoutingShard(region string) (routing, shard string) {
	switch region {
	case "na", "br", "latam":
		return "americas", "na"
	case "eu", "euw", "eune", "tr", "ru":
		return "europe", "eu"
	case "kr":
		return "asia", "kr"
	case "ap", "jp", "oce":
		return "asia", "ap"
	default:
		return "asia", "ap"
	}
}
