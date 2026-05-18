package analysis

import "testing"

func TestAnalyzePlayerGeneratesEvidenceLinkedRecommendations(t *testing.T) {
	report := AnalyzePlayer(DemoSnapshot())

	if len(report.Findings) == 0 {
		t.Fatal("expected findings")
	}
	if len(report.Recommendations) == 0 {
		t.Fatal("expected recommendations")
	}
	if len(report.MapBreakdown) == 0 || len(report.AgentBreakdown) == 0 {
		t.Fatal("expected map and agent breakdown")
	}
	if report.MapBreakdown[0].Name != "Ascent" {
		t.Fatalf("expected weakest/frequent map first, got %+v", report.MapBreakdown[0])
	}
	if len(report.PracticePlan) == 0 {
		t.Fatal("expected practice plan")
	}
	if report.PracticePlan[0].Day != 1 || len(report.PracticePlan[0].Checklist) == 0 {
		t.Fatalf("unexpected practice task: %+v", report.PracticePlan[0])
	}

	findingIDs := map[string]bool{}
	for _, finding := range report.Findings {
		findingIDs[finding.ID] = true
		if len(finding.Evidence) == 0 {
			t.Fatalf("expected evidence for %s", finding.ID)
		}
	}

	for _, recommendation := range report.Recommendations {
		if !findingIDs[recommendation.FindingID] {
			t.Fatalf("recommendation %s is not linked to a finding", recommendation.ID)
		}
		if recommendation.Drill == "" {
			t.Fatalf("recommendation %s missing drill", recommendation.ID)
		}
	}
}

func TestAnalyzePlayerHandlesMissingData(t *testing.T) {
	report := AnalyzePlayer(PlayerSnapshot{Name: "empty"})

	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 missing-data finding, got %d", len(report.Findings))
	}
	if report.Findings[0].ID != "finding-no-data" {
		t.Fatalf("unexpected finding: %s", report.Findings[0].ID)
	}
	if len(report.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(report.Recommendations))
	}
	if len(report.PracticePlan) != 1 || report.PracticePlan[0].Focus != "Setup dữ liệu" {
		t.Fatalf("unexpected missing-data practice plan: %+v", report.PracticePlan)
	}
}
