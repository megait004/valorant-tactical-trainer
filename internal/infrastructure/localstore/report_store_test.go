package localstore

import (
	"path/filepath"
	"testing"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

func TestReportStorePersistsLastReport(t *testing.T) {
	store := NewReportStoreAt(filepath.Join(t.TempDir(), "last_report.json"))

	_, ok, err := store.LoadLastReport()
	if err != nil {
		t.Fatalf("expected missing report to be safe, got %v", err)
	}
	if ok {
		t.Fatal("expected missing report")
	}

	expected := StoredReport{
		Report:    analysis.AnalyzePlayer(analysis.DemoSnapshot()),
		Source:    "https://api.henrikdev.xyz/valorant/v3/matches/ap/name/tag?size=5",
		Cached:    false,
		FetchedAt: "2026-05-15T10:00:00Z",
		Message:   "Fetch Henrik thành công qua Go core.",
	}
	if err := store.SaveLastReport(expected); err != nil {
		t.Fatalf("expected save last report, got %v", err)
	}

	loaded, ok, err := store.LoadLastReport()
	if err != nil {
		t.Fatalf("expected load last report, got %v", err)
	}
	if !ok {
		t.Fatal("expected report to exist")
	}
	if loaded.Report.Player.Name != expected.Report.Player.Name || loaded.Source != expected.Source {
		t.Fatalf("unexpected stored report: %+v", loaded)
	}
}
