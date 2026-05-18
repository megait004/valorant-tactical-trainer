package assistant

import (
	"testing"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

func TestBuildAlertsFromReportIncludesFinding(t *testing.T) {
	report := analysis.AnalyzePlayer(analysis.DemoSnapshot())
	alerts := BuildAlertsFromReport(report)

	if len(alerts) == 0 {
		t.Fatal("expected alerts from demo report")
	}

	found := false
	for _, alert := range alerts {
		if alert.ID == "finding-first-death-non-duelist" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected first-death alert, got %#v", alerts)
	}
}

func TestEngineStartAndRequestTip(t *testing.T) {
	engine := NewEngine()
	report := analysis.AnalyzePlayer(analysis.DemoSnapshot())
	state := engine.Start(report)

	if !state.Active {
		t.Fatal("expected active session")
	}
	if state.CurrentAlert == nil {
		t.Fatal("expected opening alert")
	}

	result := engine.RequestTip()
	if !result.HasTip {
		t.Fatal("expected forced tip")
	}
	if result.Alert.Message == "" {
		t.Fatal("expected alert message")
	}
}

func TestEngineMarkRoundStart(t *testing.T) {
	engine := NewEngine()
	engine.Start(analysis.AnalyzePlayer(analysis.DemoSnapshot()))

	result := engine.MarkRoundStart()
	if !result.HasTip {
		t.Fatal("expected round tip")
	}
	if engine.State().RoundCount != 1 {
		t.Fatalf("expected round count 1, got %d", engine.State().RoundCount)
	}
}

func TestPollAutoTipRespectsCooldown(t *testing.T) {
	engine := NewEngine()
	engine.Start(analysis.AnalyzePlayer(analysis.DemoSnapshot()))

	if engine.PollAutoTip().HasTip {
		t.Fatal("expected cooldown to block auto tip right after session start")
	}

	if !engine.RequestTip().HasTip {
		t.Fatal("expected forced tip to bypass cooldown")
	}
}
