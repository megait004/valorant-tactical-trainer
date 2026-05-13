package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	analysisdomain "valorant-tactical-trainer/internal/domain/analysis"
	matchdomain "valorant-tactical-trainer/internal/domain/match"
	"valorant-tactical-trainer/internal/domain/player"
	"valorant-tactical-trainer/internal/domain/rank"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()

	store, err := OpenPath(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	})

	return store
}

func TestSavePlayerWithConsentAndCurrentPlayer(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()

	account := player.Account{PUUID: "p1", Name: "Player", Tag: "VN2", Region: "ap", AccountLevel: 42}
	consent := player.Consent{PlayerPUUID: "p1", Name: "Player", Tag: "VN2", Region: "ap", Provider: "provider", ConsentVersion: player.ConsentVersion, ConsentedAt: time.Now().UTC()}

	if err := store.SavePlayerWithConsent(ctx, account, consent); err != nil {
		t.Fatalf("save player: %v", err)
	}

	current, ok, err := store.CurrentPlayer(ctx)
	if err != nil {
		t.Fatalf("current player: %v", err)
	}
	if !ok || current.PUUID != "p1" || current.AccountLevel != 42 {
		t.Fatalf("unexpected current player: ok=%v player=%+v", ok, current)
	}
}

func TestSaveMatchesDedupesByMatchAndPlayer(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()
	seedPlayer(t, store)

	matches := []matchdomain.Summary{{MatchID: "m1", PlayerPUUID: "p1", MapName: "Ascent", Kills: 10}}
	if _, err := store.SaveMatches(ctx, matches); err != nil {
		t.Fatalf("save matches: %v", err)
	}
	matches[0].Kills = 20
	if _, err := store.SaveMatches(ctx, matches); err != nil {
		t.Fatalf("save matches update: %v", err)
	}

	stored, err := store.MatchesForPlayer(ctx, "p1")
	if err != nil {
		t.Fatalf("list matches: %v", err)
	}
	if len(stored) != 1 || stored[0].Kills != 20 {
		t.Fatalf("unexpected stored matches: %+v", stored)
	}
}

func TestAPICacheExpires(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()

	if err := store.SaveAPICache(ctx, "k1", "endpoint", "payload", -time.Minute); err != nil {
		t.Fatalf("save cache: %v", err)
	}
	if _, ok, err := store.APICache(ctx, "k1"); err != nil || ok {
		t.Fatalf("expected expired cache miss, ok=%v err=%v", ok, err)
	}

	if err := store.SaveAPICache(ctx, "k1", "endpoint", "payload", time.Minute); err != nil {
		t.Fatalf("save cache: %v", err)
	}
	payload, ok, err := store.APICache(ctx, "k1")
	if err != nil || !ok || payload != "payload" {
		t.Fatalf("expected cache hit, payload=%q ok=%v err=%v", payload, ok, err)
	}
}

func TestSettingsStatsAndClearExpiredCache(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()
	seedPlayer(t, store)

	if store.Path() == "" {
		t.Fatal("expected db path")
	}
	if err := store.SaveSetting(ctx, "api_key", "secret"); err != nil {
		t.Fatalf("save setting: %v", err)
	}
	if err := store.SaveAPICache(ctx, "expired", "endpoint", "payload", -time.Minute); err != nil {
		t.Fatalf("save expired cache: %v", err)
	}
	if err := store.SaveAPICache(ctx, "fresh", "endpoint", "payload", time.Minute); err != nil {
		t.Fatalf("save fresh cache: %v", err)
	}

	stats, err := store.Stats(ctx)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Players != 1 || stats.CacheEntries != 2 || stats.ExpiredCacheEntries != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	cleared, err := store.ClearExpiredAPICache(ctx)
	if err != nil {
		t.Fatalf("clear expired cache: %v", err)
	}
	if cleared != 1 {
		t.Fatalf("expected 1 cleared cache entry, got %d", cleared)
	}
	if err := store.DeleteSetting(ctx, "api_key"); err != nil {
		t.Fatalf("delete setting: %v", err)
	}
	if _, ok, err := store.Setting(ctx, "api_key"); err != nil || ok {
		t.Fatalf("expected api key gone, ok=%v err=%v", ok, err)
	}
}

func TestSaveAndReadLatestRankSnapshot(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()
	seedPlayer(t, store)

	oldSnapshot := rank.Snapshot{PlayerPUUID: "p1", Region: "ap", TierName: "Gold 3", RankingInTier: 90, FetchedAt: time.Now().UTC().Add(-time.Hour)}
	newSnapshot := rank.Snapshot{PlayerPUUID: "p1", Region: "ap", TierName: "Platinum 1", RankingInTier: 12, Elo: 1312, FetchedAt: time.Now().UTC()}
	if err := store.SaveRankSnapshot(ctx, oldSnapshot); err != nil {
		t.Fatalf("save old rank: %v", err)
	}
	if err := store.SaveRankSnapshot(ctx, newSnapshot); err != nil {
		t.Fatalf("save new rank: %v", err)
	}

	latest, ok, err := store.LatestRankSnapshot(ctx, "p1")
	if err != nil {
		t.Fatalf("latest rank: %v", err)
	}
	if !ok || latest.TierName != "Platinum 1" || latest.Elo != 1312 {
		t.Fatalf("unexpected latest rank: ok=%v rank=%+v", ok, latest)
	}
}

func TestSaveReportAndResetAll(t *testing.T) {
	t.Parallel()
	store := newTestStore(t)
	ctx := context.Background()
	seedPlayer(t, store)

	report := analysisdomain.Report{
		PlayerPUUID:     "p1",
		GeneratedAt:     time.Now().UTC(),
		MatchCount:      1,
		AverageKDA:      2,
		HeadshotPercent: 20,
		AverageDamage:   3000,
		Summary:         "summary",
		Findings:        []analysisdomain.Finding{{Type: "baseline", Severity: "low", Confidence: 0.7, Title: "title", Description: "desc", Evidence: []string{"e1"}}},
		Recommendations: []analysisdomain.Recommendation{{Title: "rec", Drill: "drill", Priority: "low", Reason: "reason", Evidence: []string{"e1"}, Status: "new"}},
	}
	saved, err := store.SaveReport(ctx, report)
	if err != nil {
		t.Fatalf("save report: %v", err)
	}
	if saved.ID == 0 {
		t.Fatal("expected report id")
	}
	if err := store.SaveRankSnapshot(ctx, rank.Snapshot{PlayerPUUID: "p1", Region: "ap", TierName: "Gold 1", FetchedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("save rank before reset: %v", err)
	}

	if err := store.ResetAll(ctx); err != nil {
		t.Fatalf("reset all: %v", err)
	}
	_, ok, err := store.CurrentPlayer(ctx)
	if err != nil {
		t.Fatalf("current player after reset: %v", err)
	}
	if ok {
		t.Fatal("expected no current player after reset")
	}
	_, ok, err = store.LatestRankSnapshot(ctx, "p1")
	if err != nil {
		t.Fatalf("latest rank after reset: %v", err)
	}
	if ok {
		t.Fatal("expected no rank after reset")
	}
}

func seedPlayer(t *testing.T, store *Store) {
	t.Helper()
	account := player.Account{PUUID: "p1", Name: "Player", Tag: "VN2", Region: "ap"}
	consent := player.Consent{PlayerPUUID: "p1", Name: "Player", Tag: "VN2", Region: "ap", Provider: "provider", ConsentVersion: player.ConsentVersion, ConsentedAt: time.Now().UTC()}
	if err := store.SavePlayerWithConsent(context.Background(), account, consent); err != nil {
		t.Fatalf("seed player: %v", err)
	}
}
