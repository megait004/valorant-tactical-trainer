package henrik

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	datasettings "valorant-tactical-trainer/desktop/internal/domain/settings"
)

func TestStatusBlocksPersonalFetchWithoutConsent(t *testing.T) {
	status := NewClient().Status(datasettings.DefaultDataSettings())
	if status.CanFetchPersonalData {
		t.Fatal("expected personal fetch to be blocked without consent")
	}
	if !status.SafeMode {
		t.Fatal("expected safe mode")
	}
}

func TestStatusAllowsFetchWhenConsentRiotIDAndAPIKeyExist(t *testing.T) {
	settings := datasettings.DefaultDataSettings()
	settings.ConsentPersonalData = true
	settings.RiotName = "giaphue"
	settings.RiotTag = "DATN"
	settings.APIKey = "key-123"
	settings.RateLimitTier = "enhanced"

	status := NewClient().Status(settings)
	if !status.CanFetchPersonalData {
		t.Fatal("expected personal fetch to be allowed")
	}
	if status.RateLimitPerMinute != 90 {
		t.Fatalf("expected enhanced limit, got %d", status.RateLimitPerMinute)
	}
}

func TestFetchMatchSnapshotParsesAndCachesHenrikResponse(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Header.Get("X-API-Key") != "key-123" {
			t.Fatalf("expected api key header, got %q", r.Header.Get("X-API-Key"))
		}
		if r.URL.Query().Get("size") != "3" {
			t.Fatalf("expected match count query, got %q", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [{
				"metadata": {"matchid": "m-1", "map": "Ascent"},
				"players": {"all_players": [{"name": "giaphue", "tag": "DATN", "team": "Blue", "character": "Omen", "stats": {"kills": 18, "deaths": 12, "assists": 7, "headshots": 10, "bodyshots": 30, "legshots": 10}}]},
				"teams": {"red": {"has_won": false, "rounds_won": 9}, "blue": {"has_won": true, "rounds_won": 13}},
				"rounds": [{}, {}, {}]
			}]
		}`))
	}))
	defer server.Close()

	oldBaseURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = oldBaseURL }()

	settings := datasettings.DefaultDataSettings()
	settings.ConsentPersonalData = true
	settings.RiotName = "giaphue"
	settings.RiotTag = "DATN"
	settings.APIKey = "key-123"
	settings.APIKeyHeader = "X-API-Key"
	settings.MatchCount = 3
	settings.CacheTTLMinutes = 30

	client := NewClientForTest(server.Client(), t.TempDir(), func() time.Time { return time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC) })
	first, err := client.FetchMatchSnapshot(context.Background(), settings)
	if err != nil {
		t.Fatalf("expected fetch snapshot, got %v", err)
	}
	if first.Cached {
		t.Fatal("expected first response from network")
	}
	if first.Snapshot.RecentMatches[0].Agent != "Omen" || first.Snapshot.RecentMatches[0].Role != "controller" {
		t.Fatalf("unexpected snapshot: %+v", first.Snapshot.RecentMatches[0])
	}

	second, err := client.FetchMatchSnapshot(context.Background(), settings)
	if err != nil {
		t.Fatalf("expected cached snapshot, got %v", err)
	}
	if !second.Cached {
		t.Fatal("expected cached response")
	}
	if requests != 1 {
		t.Fatalf("expected one network request, got %d", requests)
	}
}

func TestFetchMatchSnapshotRequiresConsent(t *testing.T) {
	_, err := NewClientForTest(nil, t.TempDir(), nil).FetchMatchSnapshot(context.Background(), datasettings.DefaultDataSettings())
	if err == nil || !strings.Contains(err.Error(), "consent") {
		t.Fatalf("expected consent error, got %v", err)
	}
}
