package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"valorant-tactical-trainer/desktop/internal/infrastructure/env"
)

// PlayerInfo là thông tin player sau khi đăng nhập thành công (Account-V1).
type PlayerInfo struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Region   string `json:"region"`
	Shard    string `json:"shard"`
}

// LoginResult là DTO trả cho Wails AuthService.Login.
type LoginResult struct {
	Success    bool        `json:"success"`
	Error      string      `json:"error,omitempty"`
	PlayerInfo *PlayerInfo `json:"playerInfo,omitempty"`
}

// SettingsSink là bridge để Login/Logout đồng bộ session vào persistent
// settings + xoá report cache cũ khi đổi nick.
//
// Tách interface ở đây để package riot không import infrastructure/store
// (tránh import cycle).
type SettingsSink interface {
	OnLogin(player PlayerInfo) error
	OnLogout() error
}

// Login xác thực Riot ID qua Account-V1.
//
// Tại sao Account-V1 mà không phải RSO/OAuth?
//
//	RGAPI dev key KHÔNG xác thực được password user (RSO chỉ available cho
//	production OAuth client). Cách hợp lệ duy nhất với dev key là gọi
//	`/riot/account/v1/accounts/by-riot-id/{name}/{tag}` — nếu key + Riot ID
//	hợp lệ thì coi như "đăng nhập".
//
// Pipeline:
//  1. Validate input + load RIOT_API_KEY từ .env.
//  2. Resolve routing/shard theo region (xem routing.go).
//  3. HTTP GET Account-V1, handle status code chi tiết (401/403/404/429).
//  4. Parse PUUID, save vào session, gọi sink.OnLogin nếu có.
func (c *Client) Login(riotID, tagLine, region string) (LoginResult, error) {
	riotID = strings.TrimSpace(riotID)
	tagLine = strings.TrimSpace(strings.TrimPrefix(tagLine, "#"))
	region = strings.ToLower(strings.TrimSpace(region))

	if riotID == "" || tagLine == "" {
		return LoginResult{Error: "Nhập Riot ID và tag (vd: giaphue#DATN)"}, nil
	}
	apiKey := env.Load("RIOT_API_KEY")
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
