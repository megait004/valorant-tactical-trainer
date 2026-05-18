package analysis

import "context"

// Coach là cổng để cá nhân hoá nội dung text trong Report bằng dịch vụ ngoài
// (LLM). Domain không biết Coach được implement bằng gì — có thể là OpenAI,
// Claude, Gemini hoặc local model. Implementation nằm trong
// internal/infrastructure/llm.
//
// Có 2 method:
//
//   - SuggestRecommendations: legacy entry, chỉ override Recommendations.
//     Vẫn giữ để code cũ (vd hooks, test) còn chạy.
//
//   - SuggestFullReport: entry mới — LLM viết lại đồng thời Findings,
//     Recommendations và PracticePlan theo agent/map cụ thể của user. Rule
//     engine vẫn chạy trước để đảm bảo evidence, severity và confidence
//     deterministic; LLM chỉ nhận ID + metrics và "thay áo" text.
type Coach interface {
	SuggestRecommendations(ctx context.Context, snapshot PlayerSnapshot, findings []Finding) ([]Recommendation, error)
	SuggestFullReport(ctx context.Context, input CoachInput) (CoachOutput, error)
}

// CoachInput là payload gọn cho LLM khi sinh full report. Chỉ những trường
// cần thiết để giữ token thấp.
type CoachInput struct {
	Snapshot  PlayerSnapshot `json:"snapshot"`
	Metrics   MetricSummary  `json:"metrics"`
	WeakMap   string         `json:"weakMap"`
	MainAgent string         `json:"mainAgent"`
	// Findings ở đây đã có ID, severity, confidence và evidence từ rule engine.
	// LLM chỉ cần override Title + Detail + viết lại Recommendation + PracticeTask.
	Findings []Finding `json:"findings"`
}

// CoachOutput là cấu trúc trả về của LLM. ID phải khớp với Findings input.
// Nếu LLM trả entry với ID không tồn tại → caller drop entry đó.
type CoachOutput struct {
	Findings        []Finding        `json:"findings"`
	Recommendations []Recommendation `json:"recommendations"`
	PracticePlan    []PracticeTask   `json:"practicePlan"`
}

// AnalyzePlayerWithCoach giống AnalyzePlayer nhưng cho phép LLM cá nhân hoá
// Findings/Recommendations/PracticePlan. Rule engine luôn chạy trước để có
// fallback an toàn — nếu coach == nil hoặc lỗi, giữ nguyên report rule-based.
//
// Pipeline:
//  1. AnalyzePlayer → ra report deterministic (metrics + findings rule + template).
//  2. Nếu có Coach: build CoachInput → gọi SuggestFullReport.
//  3. Merge output: với mỗi Finding ID có trong rule engine, nếu LLM trả
//     phiên bản khớp ID thì replace Title/Detail; tương tự với
//     Recommendation/PracticeTask. ID nào LLM không cover → giữ template.
func AnalyzePlayerWithCoach(ctx context.Context, snapshot PlayerSnapshot, coach Coach) Report {
	report := AnalyzePlayer(snapshot)
	if coach == nil || len(report.Findings) == 0 {
		return report
	}

	input := CoachInput{
		Snapshot:  snapshot,
		Metrics:   report.Metrics,
		WeakMap:   firstBreakdownName(report.MapBreakdown),
		MainAgent: firstBreakdownName(report.AgentBreakdown),
		Findings:  report.Findings,
	}

	out, err := coach.SuggestFullReport(ctx, input)
	if err != nil {
		return report
	}

	report.Findings = mergeFindings(report.Findings, out.Findings)
	if len(out.Recommendations) > 0 && looksAlignedWithFindings(out.Recommendations, report.Findings) {
		report.Recommendations = mergeRecommendations(report.Recommendations, out.Recommendations)
	}
	if len(out.PracticePlan) > 0 {
		report.PracticePlan = mergePracticePlan(report.PracticePlan, out.PracticePlan)
	}
	return report
}

// mergeFindings giữ thứ tự + evidence/severity/confidence từ rule engine,
// chỉ replace Title/Detail nếu LLM trả phiên bản cùng ID.
func mergeFindings(base, llm []Finding) []Finding {
	byID := map[string]Finding{}
	for _, f := range llm {
		if f.ID == "" {
			continue
		}
		byID[f.ID] = f
	}
	out := make([]Finding, 0, len(base))
	for _, f := range base {
		if alt, ok := byID[f.ID]; ok {
			if alt.Title != "" {
				f.Title = alt.Title
			}
			if alt.Detail != "" {
				f.Detail = alt.Detail
			}
		}
		out = append(out, f)
	}
	return out
}

// mergeRecommendations: với mỗi Recommendation từ rule engine, nếu LLM có
// recommendation cùng FindingID thì thay nguyên; nếu không thì giữ template.
func mergeRecommendations(base, llm []Recommendation) []Recommendation {
	byFindingID := map[string]Recommendation{}
	for _, r := range llm {
		if r.FindingID == "" || r.Title == "" {
			continue
		}
		byFindingID[r.FindingID] = r
	}
	out := make([]Recommendation, 0, len(base))
	seen := map[string]bool{}
	for _, r := range base {
		if alt, ok := byFindingID[r.FindingID]; ok {
			if alt.ID == "" {
				alt.ID = r.ID
			}
			out = append(out, alt)
			seen[r.FindingID] = true
			continue
		}
		out = append(out, r)
	}
	// Nếu LLM gợi ý thêm Recommendation cho FindingID mà rule engine bỏ qua
	// (hiếm), append cuối.
	for _, r := range llm {
		if !seen[r.FindingID] && r.FindingID != "" {
			out = append(out, r)
		}
	}
	return out
}

// mergePracticePlan: ưu tiên LLM nếu trả >0 task, nhưng giữ Day numbering
// liên tục để UI render gọn.
func mergePracticePlan(base, llm []PracticeTask) []PracticeTask {
	if len(llm) == 0 {
		return base
	}
	out := make([]PracticeTask, 0, len(llm))
	for i, t := range llm {
		t.Day = i + 1
		out = append(out, t)
	}
	if len(out) > 4 {
		return out[:4]
	}
	return out
}

func looksAlignedWithFindings(recs []Recommendation, findings []Finding) bool {
	known := map[string]struct{}{}
	for _, f := range findings {
		known[f.ID] = struct{}{}
	}
	hits := 0
	for _, r := range recs {
		if _, ok := known[r.FindingID]; ok {
			hits++
		}
	}
	return hits > 0
}
