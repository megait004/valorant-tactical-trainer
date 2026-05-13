package analysis

import (
	"testing"

	matchdomain "valorant-tactical-trainer/internal/domain/match"
)

func TestGenerateReportNoMatches(t *testing.T) {
	t.Parallel()

	report := GenerateReport("p1", nil)
	if report.MatchCount != 0 {
		t.Fatalf("expected 0 matches, got %d", report.MatchCount)
	}
	if len(report.Recommendations) == 0 {
		t.Fatal("expected import recommendation")
	}
}

func TestGenerateReportLowHeadshotFinding(t *testing.T) {
	t.Parallel()

	matches := []matchdomain.Summary{
		{MatchID: "1", Kills: 10, Deaths: 10, Assists: 2, Headshots: 4, Bodyshots: 50, Legshots: 10, DamageMade: 3000, MapName: "Ascent", Agent: "Sova"},
		{MatchID: "2", Kills: 8, Deaths: 12, Assists: 4, Headshots: 3, Bodyshots: 45, Legshots: 12, DamageMade: 2500, MapName: "Ascent", Agent: "Sova"},
		{MatchID: "3", Kills: 7, Deaths: 11, Assists: 3, Headshots: 2, Bodyshots: 35, Legshots: 8, DamageMade: 2200, MapName: "Bind", Agent: "Sova"},
	}

	report := GenerateReport("p1", matches)
	if report.HeadshotPercent >= 18 {
		t.Fatalf("expected low headshot percent, got %.2f", report.HeadshotPercent)
	}
	if len(report.Findings) == 0 {
		t.Fatal("expected findings")
	}
}
