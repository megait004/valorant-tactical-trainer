package henrik

import (
	"strings"

	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
	"valorant-tactical-trainer/desktop/internal/infrastructure/env"
)

// Status là DTO trả về frontend mỗi khi UI mở Settings — cho user biết
// app có sẵn sàng fetch dữ liệu cá nhân chưa, rate limit là bao nhiêu, và
// bước tiếp theo cần làm gì.
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

// Status trả về trạng thái sẵn sàng fetch theo settings hiện tại. Logic
// guard gồm 3 tầng:
//
//  1. Consent — không bật consent thì block hoàn toàn (privacy first).
//  2. Riot ID — phải có name + tag để build URL.
//  3. API key — không bắt buộc (public tier 30req/min), có HDEV key thì lên
//     90req/min.
func (c *Client) Status(settings datasettings.DataSettings) Status {
	envHenrikKey := strings.TrimSpace(env.Load("HENRIK_API_KEY"))
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
