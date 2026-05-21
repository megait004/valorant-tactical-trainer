// Package env cung cấp helper đọc biến môi trường có fallback xuống file .env
// ở root project. Dùng chung cho mọi adapter cần API key (riot, henrik, llm).
//
// Cách dùng:
//
//	apiKey := env.Load("RIOT_API_KEY")
//
// Thứ tự ưu tiên:
//  1. os.Getenv(name) — env thực sự của process
//  2. .env file ở working directory
//  3. .env file ở thư mục cha (khi chạy từ subfolder build/)
//  4. .env file cạnh executable (binary release)
//
// Format .env hỗ trợ: KEY=VALUE, comment bằng `#`, value có quote
// ("..." hoặc '...'), optional `export ` prefix.
package env

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	dotenvOnce sync.Once
	dotenvData map[string]string
)

// Load trả về giá trị của 1 env key, fallback xuống file .env.
// Trả về "" nếu không tìm thấy.
func Load(name string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	dotenvOnce.Do(func() {
		dotenvData = readDotenvFromKnownLocations()
	})
	return strings.TrimSpace(dotenvData[name])
}

func readDotenvFromKnownLocations() map[string]string {
	result := map[string]string{}

	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, ".env"))
		candidates = append(candidates, filepath.Join(filepath.Dir(cwd), ".env"))
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), ".env"))
	}

	seen := map[string]struct{}{}
	for _, path := range candidates {
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}

		values, ok := parseDotenv(path)
		if !ok {
			continue
		}
		for k, v := range values {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}

func parseDotenv(path string) (map[string]string, bool) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])

		if !isQuoted(value) {
			if hash := strings.IndexByte(value, '#'); hash >= 0 {
				value = strings.TrimSpace(value[:hash])
			}
		}
		value = unquote(value)

		if key != "" {
			values[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, false
	}
	return values, true
}

func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	first, last := s[0], s[len(s)-1]
	return (first == '"' && last == '"') || (first == '\'' && last == '\'')
}

func unquote(s string) string {
	if isQuoted(s) {
		return s[1 : len(s)-1]
	}
	return s
}
