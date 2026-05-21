package analysis

// generateFindings chạy từng pattern trong registry, gắn evidence/severity/
// confidence do pattern tự tính. Trường hợp đặc biệt: snapshot rỗng → trả
// đúng 1 Finding "no data" để UI biết hướng dẫn user fetch report trước.
func generateFindings(snapshot PlayerSnapshot, metrics MetricSummary) []Finding {
	if metrics.Matches == 0 {
		return []Finding{noDataFinding()}
	}

	findings := make([]Finding, 0, len(findingPatterns))
	for _, pattern := range findingPatterns {
		evidence, ok := pattern.Match(snapshot, metrics)
		if !ok {
			continue
		}
		findings = append(findings, Finding{
			ID:         pattern.ID,
			Title:      pattern.DefaultTitle,
			Severity:   pattern.SeverityFn(metrics),
			Confidence: pattern.ConfidenceFn(metrics),
			Detail:     pattern.DefaultDetail(metrics),
			Evidence:   []Evidence{evidence},
		})
	}
	return findings
}

// generateRecommendations sinh Recommendation tương ứng từng Finding bằng cách
// tra cứu pattern.DefaultRecommendation. Finding "no data" có recommendation
// riêng vì không tra từ pattern registry.
func generateRecommendations(findings []Finding, metrics MetricSummary) []Recommendation {
	recs := make([]Recommendation, 0, len(findings))
	for _, f := range findings {
		if f.ID == "finding-no-data" {
			recs = append(recs, noDataRecommendation())
			continue
		}
		pattern := FindingPatternByID(f.ID)
		if pattern == nil || pattern.DefaultRecommendation == nil {
			continue
		}
		recs = append(recs, pattern.DefaultRecommendation(metrics))
	}
	return recs
}

// noDataFinding là Finding placeholder khi user chưa có match history.
func noDataFinding() Finding {
	return Finding{
		ID:         "finding-no-data",
		Title:      "Chưa có dữ liệu trận để phân tích",
		Severity:   "low",
		Confidence: "high",
		Detail:     "Cần import match history hoặc VOD trước khi sinh giáo án cá nhân.",
		Evidence:   []Evidence{{Metric: "matches", Value: 0, SampleSize: 0}},
	}
}

// noDataRecommendation là Recommendation kèm Finding no-data.
func noDataRecommendation() Recommendation {
	return Recommendation{
		ID:        "rec-import-data",
		FindingID: "finding-no-data",
		Title:     "Import dữ liệu trước khi coach",
		Reason:    "Không có evidence thì app không nên sinh lời khuyên giả.",
		Drill:     "Login Riot và bấm Fetch report để có match history thật.",
		Cadence:   "Một lần setup",
	}
}

// collectMatchIDs gom toàn bộ match ID — dùng làm evidence khi rule kích hoạt
// trên cả sample. Dùng cho pattern không gắn riêng map (vd low HS%).
func collectMatchIDs(matches []MatchSummary) []string {
	ids := make([]string, 0, len(matches))
	for _, m := range matches {
		ids = append(ids, m.ID)
	}
	return ids
}

// matchesForMap lọc match ID theo tên map — dùng làm evidence cho pattern
// map-gap (chỉ những trận trên map yếu).
func matchesForMap(matches []MatchSummary, mapName string) []string {
	ids := []string{}
	for _, m := range matches {
		if m.Map == mapName {
			ids = append(ids, m.ID)
		}
	}
	return ids
}
