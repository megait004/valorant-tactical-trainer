package analysis

// FindingPattern là khai báo deterministic của một "trigger". Mỗi pattern bao
// gồm:
//
//   - ID: định danh ổn định, dùng để match Finding ↔ Recommendation ↔ PracticeTask.
//   - Match: nhận snapshot + metrics, trả về Evidence + true nếu kích hoạt.
//   - SeverityFn / ConfidenceFn: tính severity và confidence từ metrics
//     (vẫn deterministic, dễ test).
//   - DefaultTitle/Detail: chuỗi tiếng Việt fallback dùng khi LLM không có
//     hoặc trả lỗi. LLM có thể override.
//
// Toàn bộ ngưỡng số (vd FirstDeathRate >= 0.18) tập trung trong file này, code
// khác không hard-code threshold trực tiếp nữa.
type FindingPattern struct {
	ID            string
	DefaultTitle  string
	DefaultDetail func(metrics MetricSummary) string
	Match         func(snapshot PlayerSnapshot, metrics MetricSummary) (Evidence, bool)
	SeverityFn    func(metrics MetricSummary) string
	ConfidenceFn  func(metrics MetricSummary) string
	// DefaultRecommendation và DefaultPracticeTask là template fallback —
	// chỉ dùng khi LLM Coach không khả dụng. Nhận metrics + map/agent
	// chính để thay placeholder cơ bản.
	DefaultRecommendation func(metrics MetricSummary) Recommendation
	DefaultPracticeTask   func(metrics MetricSummary, weakMap, mainAgent string, day int) PracticeTask
}

// findingPatterns là registry duy nhất chứa toàn bộ luật. Thêm/đổi luật chỉ
// cần edit slice này — không cần đụng generateFindings/generateRecommendations.
var findingPatterns = []FindingPattern{
	patternFirstDeathNonDuelist(),
	patternLowHeadshot(),
	patternLowSurvivalImpact(),
	patternMapGap(),
}

// FindingPatterns trả về snapshot read-only của registry để code ngoài
// (ví dụ LLM coach hoặc test) liệt kê được toàn bộ ID.
func FindingPatterns() []FindingPattern {
	out := make([]FindingPattern, len(findingPatterns))
	copy(out, findingPatterns)
	return out
}

// FindingPatternByID tra cứu pattern theo ID. Trả về nil nếu không có.
func FindingPatternByID(id string) *FindingPattern {
	for i := range findingPatterns {
		if findingPatterns[i].ID == id {
			return &findingPatterns[i]
		}
	}
	return nil
}
