package valorantapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLookupAccount(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/account/Player/VN2" {
			t.Fatalf("unexpected path %s", request.URL.Path)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"status": 200,
			"data": {
				"puuid": "p1",
				"region": "ap",
				"account_level": 123,
				"name": "Player",
				"tag": "VN2",
				"card": {"small": "small.png", "large": "large.png"},
				"last_update": "today"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	account, err := client.LookupAccount(context.Background(), "Player", "#VN2")
	if err != nil {
		t.Fatalf("lookup account: %v", err)
	}
	if account.PUUID != "p1" || account.Region != "ap" || account.AccountLevel != 123 {
		t.Fatalf("unexpected account: %+v", account)
	}
}

func TestLookupAccountRateLimited(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	_, err := client.LookupAccount(context.Background(), "Player", "VN2")
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestMatchesByPUUID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/by-puuid/matches/ap/p1" {
			t.Fatalf("unexpected path %s", request.URL.Path)
		}
		if request.URL.Query().Get("size") != "10" {
			t.Fatalf("unexpected size %s", request.URL.Query().Get("size"))
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"status": 200,
			"data": [{
				"metadata": {
					"map": "Ascent",
					"game_length": 1800,
					"game_start": 1710000000,
					"rounds_played": 22,
					"mode": "Competitive",
					"queue": "competitive",
					"season_id": "s1",
					"matchid": "m1",
					"region": "ap",
					"cluster": "sg"
				},
				"players": {"all_players": [{
					"puuid": "p1",
					"team": "Blue",
					"character": "Sova",
					"stats": {"kills": 20, "deaths": 10, "assists": 5, "headshots": 12, "bodyshots": 30, "legshots": 4},
					"damage_made": 4200
				}]}
			}]
		}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	matches, raw, err := client.MatchesByPUUID(context.Background(), "p1", "ap", "10")
	if err != nil {
		t.Fatalf("matches by puuid: %v", err)
	}
	if raw == "" {
		t.Fatal("expected raw payload")
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	match := matches[0]
	if match.MatchID != "m1" || match.MapName != "Ascent" || match.Agent != "Sova" || match.Kills != 20 {
		t.Fatalf("unexpected match: %+v", match)
	}
}

func TestMMRByPUUID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/by-puuid/mmr/ap/p1" {
			t.Fatalf("unexpected path %s", request.URL.Path)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"status": 200,
			"data": {
				"currenttier": 18,
				"currenttierpatched": "Diamond 1",
				"ranking_in_tier": 67,
				"mmr_change_to_last_game": 14,
				"elo": 1467,
				"season_id": "s1"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	snapshot, raw, err := client.MMRByPUUID(context.Background(), "p1", "ap")
	if err != nil {
		t.Fatalf("mmr by puuid: %v", err)
	}
	if raw == "" {
		t.Fatal("expected raw payload")
	}
	if snapshot.TierName != "Diamond 1" || snapshot.RankingInTier != 67 || snapshot.Elo != 1467 {
		t.Fatalf("unexpected snapshot: %+v", snapshot)
	}
}

func TestRateLimiterWaitsBetweenRequests(t *testing.T) {
	t.Parallel()

	limiter := NewRateLimiter(20 * time.Millisecond)
	client := NewClient(WithRateLimiter(limiter))
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{
			"status": 200,
			"data": {
				"puuid": "p1",
				"region": "ap",
				"account_level": 123,
				"name": "Player",
				"tag": "VN2",
				"card": {"small": "small.png", "large": "large.png"},
				"last_update": "today"
			}
		}`))
	}))
	defer server.Close()

	client.baseURL = server.URL
	started := time.Now()
	if _, err := client.LookupAccount(context.Background(), "Player", "VN2"); err != nil {
		t.Fatalf("first lookup: %v", err)
	}
	if _, err := client.LookupAccount(context.Background(), "Player", "VN2"); err != nil {
		t.Fatalf("second lookup: %v", err)
	}

	if requests != 2 {
		t.Fatalf("expected 2 requests, got %d", requests)
	}
	if elapsed := time.Since(started); elapsed < 20*time.Millisecond {
		t.Fatalf("expected throttle delay, got %s", elapsed)
	}
}

func TestRateLimiterReservesConcurrentSlots(t *testing.T) {
	t.Parallel()

	limiter := NewRateLimiter(10 * time.Millisecond)
	ctx := context.Background()
	started := time.Now()
	done := make(chan error, 3)

	for range 3 {
		go func() {
			done <- limiter.Wait(ctx)
		}()
	}

	for range 3 {
		if err := <-done; err != nil {
			t.Fatalf("wait: %v", err)
		}
	}

	if elapsed := time.Since(started); elapsed < 20*time.Millisecond {
		t.Fatalf("expected reserved slots, got %s", elapsed)
	}
}
