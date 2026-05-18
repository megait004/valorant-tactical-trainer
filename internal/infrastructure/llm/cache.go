package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

// Disk cache cho LLM response để tránh gọi lại cùng prompt trong TTL.
// Key = SHA256(provider + model + payload). 2 loại cache tách biệt:
//
//   - SuggestRecommendations → []analysis.Recommendation
//   - SuggestFullReport      → analysis.CoachOutput
//
// File nằm trong %CacheDir%/Valorant Tactical Trainer/llm/.

// --- Recommendations cache ---

type recommendationsCacheEnvelope struct {
	FetchedAt time.Time                 `json:"fetchedAt"`
	Recs      []analysis.Recommendation `json:"recs"`
}

func (c *Coach) cacheKeyRecommendations(payload promptPayload) string {
	body, _ := json.Marshal(struct {
		Provider string        `json:"p"`
		Model    string        `json:"m"`
		Payload  promptPayload `json:"d"`
	}{c.provider, c.model, payload})
	h := sha256.Sum256(body)
	return filepath.Join(c.cacheDir, hex.EncodeToString(h[:])+".json")
}

func (c *Coach) readCacheRecommendations(path string) ([]analysis.Recommendation, bool) {
	if c.cacheTTL <= 0 {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var env recommendationsCacheEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, false
	}
	if env.FetchedAt.IsZero() || c.now().Sub(env.FetchedAt) > c.cacheTTL {
		return nil, false
	}
	return env.Recs, true
}

func (c *Coach) writeCacheRecommendations(path string, recs []analysis.Recommendation) {
	if c.cacheTTL <= 0 {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(recommendationsCacheEnvelope{FetchedAt: c.now(), Recs: recs}, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

// --- Full report cache ---

type fullReportCacheEnvelope struct {
	FetchedAt time.Time            `json:"fetchedAt"`
	Output    analysis.CoachOutput `json:"output"`
}

func (c *Coach) cacheKeyFullReport(input analysis.CoachInput) string {
	body, _ := json.Marshal(struct {
		Kind     string              `json:"k"`
		Provider string              `json:"p"`
		Model    string              `json:"m"`
		Input    analysis.CoachInput `json:"i"`
	}{"full-report", c.provider, c.model, input})
	h := sha256.Sum256(body)
	return filepath.Join(c.cacheDir, "full-"+hex.EncodeToString(h[:])+".json")
}

func (c *Coach) readCacheFullReport(path string) (analysis.CoachOutput, bool) {
	if c.cacheTTL <= 0 {
		return analysis.CoachOutput{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return analysis.CoachOutput{}, false
	}
	var env fullReportCacheEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return analysis.CoachOutput{}, false
	}
	if env.FetchedAt.IsZero() || c.now().Sub(env.FetchedAt) > c.cacheTTL {
		return analysis.CoachOutput{}, false
	}
	return env.Output, true
}

func (c *Coach) writeCacheFullReport(path string, out analysis.CoachOutput) {
	if c.cacheTTL <= 0 {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(fullReportCacheEnvelope{FetchedAt: c.now(), Output: out}, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}
