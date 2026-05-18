package riot

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const envKey = "RIOT_API_KEY"

var (
	envOnce   sync.Once
	envLoaded map[string]string
)

// LoadEnvKey trả về giá trị của một env key theo thứ tự ưu tiên:
//  1. biến môi trường thực sự (os.Getenv)
//  2. file .env ở working dir / parent / cạnh executable
//
// Trả về "" nếu không tìm thấy. Dùng chung cho cả RIOT_API_KEY và
// HENRIK_API_KEY (Henrik client gọi qua đây để load HDEV key).
func LoadEnvKey(name string) string {
	if v := strings.TrimSpace(os.Getenv(name)); v != "" {
		return v
	}
	envOnce.Do(func() {
		envLoaded = readDotenvFromKnownLocations()
	})
	return strings.TrimSpace(envLoaded[name])
}

// loadAPIKey trả về RGAPI key cho Riot Account-V1.
func loadAPIKey() string {
	return LoadEnvKey(envKey)
}

func readDotenvFromKnownLocations() map[string]string {
	result := map[string]string{}

	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(cwd, ".env"))
		// Cũng thử lùi 1 cấp (nếu chạy từ subfolder build/)
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

// parseDotenv đọc file .env đơn giản theo format KEY=VALUE.
// Hỗ trợ comment (#), value có quote ("..." hoặc '...'), bỏ qua dòng trống.
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
		// hỗ trợ optional `export ` prefix
		line = strings.TrimPrefix(line, "export ")

		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])

		// strip inline comment khi value không quoted
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
