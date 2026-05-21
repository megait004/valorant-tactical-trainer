package analysis

// AnalyzePlayer là entry chính của rule engine — nhận snapshot match history,
// trả về Report đầy đủ (deterministic, không gọi LLM).
//
// Pipeline:
//
//  1. calculateMetrics    → MetricSummary tổng
//  2. generateFindings    → các Finding từ pattern registry
//  3. generateRecommendations → từng Recommendation tương ứng
//  4. calculateBreakdown × 2  → map / agent breakdown
//  5. generatePracticePlan    → giáo án 4 ngày
//
// Bản pure → unit test deterministic. Coach LLM hook được tách riêng trong
// AnalyzePlayerWithCoach (xem coach.go).
func AnalyzePlayer(snapshot PlayerSnapshot) Report {
	metrics := calculateMetrics(snapshot)
	findings := generateFindings(snapshot, metrics)
	recommendations := generateRecommendations(findings, metrics)
	mapBreakdown := calculateBreakdown(snapshot.RecentMatches, func(m MatchSummary) string { return m.Map }, true)
	agentBreakdown := calculateBreakdown(snapshot.RecentMatches, func(m MatchSummary) string { return m.Agent }, false)

	return Report{
		Player:          snapshot,
		Metrics:         metrics,
		MapBreakdown:    mapBreakdown,
		AgentBreakdown:  agentBreakdown,
		PracticePlan:    generatePracticePlan(metrics, mapBreakdown, agentBreakdown, findings),
		Findings:        findings,
		Recommendations: recommendations,
	}
}
