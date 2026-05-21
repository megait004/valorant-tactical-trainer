package henrik

import (
	"fmt"
	"time"
)

// reserveRequest dùng sliding window 1 phút để giới hạn số request gửi sang
// Henrik. Mục tiêu là tránh user vô tình bấm fetch liên tục và bị Henrik
// block 429 → quota của tất cả user share key chung sẽ tệ.
//
// Lưu ý: limit này là local cho client instance, không persist sau khi app
// restart — Henrik server-side rate limit vẫn là final guard.
func (c *Client) reserveRequest(limit int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	windowStart := now.Add(-1 * time.Minute)
	kept := c.requests[:0]
	for _, item := range c.requests {
		if item.After(windowStart) {
			kept = append(kept, item)
		}
	}
	c.requests = kept

	if len(c.requests) >= limit {
		return fmt.Errorf("local rate limit %d/min đã đầy, thử lại sau", limit)
	}
	c.requests = append(c.requests, now)
	return nil
}
