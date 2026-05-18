package llm

import (
	"encoding/json"
	"strings"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

// File này tập trung mọi logic parse + sanitize JSON do LLM trả về.
// LLM thường thêm ```json fence dù prompt cấm — stripCodeFence xử lý case đó.
//
// Các function:
//   - parseRecommendationsJSON: cho SuggestRecommendations (legacy).
//   - parseFullReportJSON:      cho SuggestFullReport.

type llmRecommendationsOutput struct {
	Recommendations []analysis.Recommendation `json:"recommendations"`
}

func parseRecommendationsJSON(raw string) ([]analysis.Recommendation, error) {
	cleaned := stripCodeFence(strings.TrimSpace(raw))
	var out llmRecommendationsOutput
	if err := json.Unmarshal([]byte(cleaned), &out); err != nil {
		return nil, err
	}
	return sanitizeRecommendations(out.Recommendations), nil
}

func parseFullReportJSON(raw string) (analysis.CoachOutput, error) {
	cleaned := stripCodeFence(strings.TrimSpace(raw))
	var out analysis.CoachOutput
	if err := json.Unmarshal([]byte(cleaned), &out); err != nil {
		return analysis.CoachOutput{}, err
	}
	out.Findings = sanitizeFindings(out.Findings)
	out.Recommendations = sanitizeRecommendations(out.Recommendations)
	out.PracticePlan = sanitizePracticeTasks(out.PracticePlan)
	return out, nil
}

// sanitizeFindings trim chuỗi + drop entry không có ID. KHÔNG fabricate field.
func sanitizeFindings(in []analysis.Finding) []analysis.Finding {
	out := make([]analysis.Finding, 0, len(in))
	for _, f := range in {
		f.ID = strings.TrimSpace(f.ID)
		f.Title = strings.TrimSpace(f.Title)
		f.Detail = strings.TrimSpace(f.Detail)
		if f.ID == "" {
			continue
		}
		out = append(out, f)
	}
	return out
}

// sanitizeRecommendations trim chuỗi + drop entry không có FindingID/Title.
// Tự gen ID nếu LLM quên.
func sanitizeRecommendations(in []analysis.Recommendation) []analysis.Recommendation {
	out := make([]analysis.Recommendation, 0, len(in))
	for _, r := range in {
		r.ID = strings.TrimSpace(r.ID)
		r.FindingID = strings.TrimSpace(r.FindingID)
		r.Title = strings.TrimSpace(r.Title)
		r.Reason = strings.TrimSpace(r.Reason)
		r.Drill = strings.TrimSpace(r.Drill)
		r.Cadence = strings.TrimSpace(r.Cadence)
		if r.FindingID == "" || r.Title == "" {
			continue
		}
		if r.ID == "" {
			r.ID = "rec-" + r.FindingID
		}
		out = append(out, r)
	}
	return out
}

// sanitizePracticeTasks trim chuỗi + drop task không có Focus.
func sanitizePracticeTasks(in []analysis.PracticeTask) []analysis.PracticeTask {
	out := make([]analysis.PracticeTask, 0, len(in))
	for _, t := range in {
		t.Focus = strings.TrimSpace(t.Focus)
		t.Map = strings.TrimSpace(t.Map)
		t.Agent = strings.TrimSpace(t.Agent)
		t.Duration = strings.TrimSpace(t.Duration)
		t.Evidence = strings.TrimSpace(t.Evidence)
		if t.Focus == "" {
			continue
		}
		out = append(out, t)
	}
	return out
}

// stripCodeFence loại bỏ ```json … ``` nếu LLM cố chèn.
func stripCodeFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		s = s[idx+1:]
	}
	if i := strings.LastIndex(s, "```"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
