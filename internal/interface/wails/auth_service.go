package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/infrastructure/riot"
)

// AuthService expose Riot auth methods cho Wails frontend.
//
// Cơ chế: do RGAPI dev key KHÔNG cho phép xác thực password user (đó là RSO/OAuth
// production-only), app dùng Account-V1 để verify Riot ID + key. User nhập
// gameName#tagLine + region + key, app gọi /riot/account/v1/accounts/by-riot-id
// → nếu key hợp lệ và Riot ID tồn tại thì coi như "đăng nhập".
type AuthService struct {
	client *riot.Client
}

// RiotPlayerInfo là DTO trả về frontend
type RiotPlayerInfo struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
	Region   string `json:"region"`
	Shard    string `json:"shard"`
}

// RiotLoginResult là DTO trả về frontend sau Login()
type RiotLoginResult struct {
	Success    bool            `json:"success"`
	Error      string          `json:"error,omitempty"`
	PlayerInfo *RiotPlayerInfo `json:"playerInfo,omitempty"`
}

// Login xác thực Riot ID qua Account-V1.
// API key được Go core đọc từ env / file .env, frontend không cần nhập.
func (s *AuthService) Login(riotID, tagLine, region string) (RiotLoginResult, error) {
	result, err := s.client.Login(riotID, tagLine, region)
	if err != nil {
		return RiotLoginResult{Error: err.Error()}, nil
	}
	return mapLoginResult(result), nil
}

// GetPlayerInfo trả về thông tin người chơi đang đăng nhập
func (s *AuthService) GetPlayerInfo() *RiotPlayerInfo {
	info := s.client.GetPlayerInfo()
	if info == nil {
		return nil
	}
	return &RiotPlayerInfo{
		PUUID:    info.PUUID,
		GameName: info.GameName,
		TagLine:  info.TagLine,
		Region:   info.Region,
		Shard:    info.Shard,
	}
}

// IsLoggedIn kiểm tra trạng thái đăng nhập
func (s *AuthService) IsLoggedIn() bool {
	return s.client.IsLoggedIn()
}

// Logout xóa session hiện tại
func (s *AuthService) Logout() {
	s.client.Logout()
}

// mapLoginResult convert riot.LoginResult → RiotLoginResult
func mapLoginResult(r riot.LoginResult) RiotLoginResult {
	res := RiotLoginResult{
		Success: r.Success,
		Error:   r.Error,
	}
	if r.PlayerInfo != nil {
		res.PlayerInfo = &RiotPlayerInfo{
			PUUID:    r.PlayerInfo.PUUID,
			GameName: r.PlayerInfo.GameName,
			TagLine:  r.PlayerInfo.TagLine,
			Region:   r.PlayerInfo.Region,
			Shard:    r.PlayerInfo.Shard,
		}
	}
	return res
}
