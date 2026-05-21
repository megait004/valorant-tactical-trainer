package henrik

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// cacheEnvelope là wrapper lưu thêm timestamp để verify TTL khi đọc.
type cacheEnvelope struct {
	FetchedAt time.Time       `json:"fetchedAt"`
	Body      json.RawMessage `json:"body"`
}

// cachePath sinh đường dẫn cache file = hash(url) — vừa unique, vừa giấu
// PUUID/tên player khỏi tên file trên disk.
func (c *Client) cachePath(requestURL string) string {
	hash := sha256.Sum256([]byte(requestURL))
	return filepath.Join(c.cacheDir, hex.EncodeToString(hash[:])+".json")
}

// readCache trả body cache nếu còn trong TTL. TTL ≤ 0 → bypass cache (force
// fetch). Lỗi đọc/parse file → coi như cache miss để không crash.
func (c *Client) readCache(path string, ttl time.Duration) (cacheEnvelope, bool) {
	if ttl <= 0 {
		return cacheEnvelope{}, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cacheEnvelope{}, false
	}
	var envelope cacheEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return cacheEnvelope{}, false
	}
	if envelope.FetchedAt.IsZero() || c.now().Sub(envelope.FetchedAt) > ttl {
		return cacheEnvelope{}, false
	}
	return envelope, true
}

// writeCache ghi body kèm timestamp. Best-effort — lỗi I/O bị nuốt vì cache
// chỉ là tối ưu, không ảnh hưởng kết quả fetch.
func (c *Client) writeCache(path string, envelope cacheEnvelope) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}
