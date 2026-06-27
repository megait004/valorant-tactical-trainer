// Package llm chứa adapter gọi LLM (hiện implement Google Gemini) để cá nhân
// hoá Recommendations và bot AI chat. Nếu .env không có LLM_API_KEY thì
// factory trả về nil — domain layer sẽ tự fallback template.
//
// Cấu trúc file:
//
//	coach.go           — Coach struct, factory (NewCoachFromEnv), Complete dispatcher.
//	prompts.go         — System prompts (defaultSystemPrompt, fullReportSystemPrompt).
//	suggest.go         — SuggestRecommendations + SuggestFullReport (analysis.Coach impl).
//	gemini_provider.go — Adapter Google Gemini generateContent + retry 429.
//	parser.go          — Sanitize JSON LLM trả về (trim, drop entry rỗng).
//	cache.go           — Disk cache TTL theo provider+model+input.
//	payload.go         — promptPayload gọn (snapshot + findings) để giảm token.
//
package llm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"valorant-tactical-trainer/desktop/internal/infrastructure/env"
)

const (
	envAPIKey       = "LLM_API_KEY"
	envModel        = "LLM_MODEL"
	envCacheMinutes = "LLM_CACHE_MINUTES"

	providerGemini = "gemini"

	// gemini-2.5-flash-lite: free tier rộng (15 RPM, 1000 RPD), được mở cho hầu
	// hết account. Nếu account user không có quyền, set LLM_MODEL trong .env.
	defaultGeminiModel  = "gemini-2.5-flash-lite"
	defaultCacheMinutes = 60
	defaultTimeout      = 30 * time.Second
)

// Coach là implementation của analysis.Coach gọi LLM ngoài.
type Coach struct {
	httpClient   *http.Client
	provider     string
	apiKey       string
	model        string
	cacheDir     string
	cacheTTL     time.Duration
	now          func() time.Time
	maxFindings  int
	systemPrompt string
}

// Option cho phép override cấu hình mặc định khi khởi tạo Coach (cho test).
type Option func(*Coach)

func WithHTTPClient(c *http.Client) Option { return func(co *Coach) { co.httpClient = c } }
func WithCacheDir(dir string) Option       { return func(co *Coach) { co.cacheDir = dir } }
func WithNow(now func() time.Time) Option  { return func(co *Coach) { co.now = now } }

// NewCoachFromEnv đọc env (LLM_API_KEY, LLM_MODEL...) qua env.Load để cũng hỗ
// trợ file .env. Trả nil nếu thiếu key. App hiện chỉ dùng Google Gemini.
func NewCoachFromEnv(opts ...Option) *Coach {
	apiKey := strings.TrimSpace(env.Load(envAPIKey))
	if apiKey == "" {
		return nil
	}
	model := pickModel(env.Load(envModel))
	cacheTTL := parseCacheTTL(env.Load(envCacheMinutes))

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}

	c := &Coach{
		httpClient:   &http.Client{Timeout: defaultTimeout},
		provider:     providerGemini,
		apiKey:       apiKey,
		model:        model,
		cacheDir:     filepath.Join(cacheDir, "Valorant Tactical Trainer", "llm"),
		cacheTTL:     cacheTTL,
		now:          time.Now,
		maxFindings:  6,
		systemPrompt: defaultSystemPrompt,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func pickModel(override string) string {
	if v := strings.TrimSpace(override); v != "" {
		return v
	}
	return defaultGeminiModel
}

func parseCacheTTL(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Duration(defaultCacheMinutes) * time.Minute
	}
	var minutes int
	if _, err := fmt.Sscanf(value, "%d", &minutes); err == nil && minutes >= 0 {
		return time.Duration(minutes) * time.Minute
	}
	return time.Duration(defaultCacheMinutes) * time.Minute
}

// ChatMessage là 1 turn trong cuộc hội thoại với role "user" hoặc "assistant".
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Complete là entry point gọi Gemini. Trả plain text reply.
// jsonMode = true thì yêu cầu Gemini trả JSON.
func (c *Coach) Complete(ctx context.Context, system string, messages []ChatMessage, jsonMode bool) (string, error) {
	if c == nil {
		return "", errors.New("coach chưa được khởi tạo (thiếu LLM_API_KEY trong .env)")
	}
	if len(messages) == 0 {
		return "", errors.New("messages rỗng")
	}
	return c.callGemini(ctx, system, messages, jsonMode)
}

// truncate cắt chuỗi quá dài cho error message đỡ ngợp.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
