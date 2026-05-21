package riot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// cacheEnvelope wraps response body + timestamp để verify TTL khi đọc.
type cacheEnvelope struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Body      json.RawMessage `json:"body"`
}

// cachePath sinh đường dẫn cache file = hash(url) — vừa unique, vừa giấu
// PUUID/match ID khỏi tên file trên disk.
func (c *MatchClient) cachePath(url string) string {
	hash := sha256.Sum256([]byte(url))
	return filepath.Join(c.cacheDir, hex.EncodeToString(hash[:])+".json")
}

func (c *MatchClient) readCache(path string, ttl time.Duration) (cacheEnvelope, bool) {
	if ttl <= 0 {
		return cacheEnvelope{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cacheEnvelope{}, false
	}
	var envlp cacheEnvelope
	if err := json.Unmarshal(data, &envlp); err != nil {
		return cacheEnvelope{}, false
	}
	if envlp.FetchedAt.IsZero() || c.now().Sub(envlp.FetchedAt) > ttl {
		return cacheEnvelope{}, false
	}
	return envlp, true
}

func (c *MatchClient) writeCache(path string, envlp cacheEnvelope) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(envlp, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

// reserveRequest local rate limiter sliding window 1 phút.
// Riot dev key chính thức là 20req/sec, 100req/2min. Limit ở đây cẩn thận
// hơn (50/min) để tránh user vô tình spam.
func (c *MatchClient) reserveRequest() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	windowStart := now.Add(-1 * time.Minute)
	kept := c.requests[:0]
	for _, t := range c.requests {
		if t.After(windowStart) {
			kept = append(kept, t)
		}
	}
	c.requests = kept

	if len(c.requests) >= 50 {
		return errors.New("local rate limit 50/min đã đầy, thử lại sau")
	}
	c.requests = append(c.requests, now)
	return nil
}
