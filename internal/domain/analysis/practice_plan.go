package analysis

// generatePracticePlan sinh giáo án 1-4 ngày. Mỗi Finding tương ứng 1 task
// (lấy template từ pattern.DefaultPracticeTask). Nếu không có Finding nào kích
// hoạt thì sinh 1 task "maintain form" mặc định.
//
// Max 4 task để UI gọn — user có thể luyện tiếp tuần sau và rule engine sẽ
// sinh plan mới theo metrics mới.
func generatePracticePlan(metrics MetricSummary, mapBreakdown, agentBreakdown []BreakdownRow, findings []Finding) []PracticeTask {
	if metrics.Matches == 0 {
		return []PracticeTask{noDataPracticeTask()}
	}

	weakMap := firstBreakdownName(mapBreakdown)
	mainAgent := firstBreakdownName(agentBreakdown)

	plan := []PracticeTask{}
	for _, f := range findings {
		pattern := FindingPatternByID(f.ID)
		if pattern == nil || pattern.DefaultPracticeTask == nil {
			continue
		}
		plan = append(plan, pattern.DefaultPracticeTask(metrics, weakMap, mainAgent, len(plan)+1))
	}

	if len(plan) == 0 {
		plan = append(plan, maintainFormTask(weakMap, mainAgent))
	}
	if len(plan) > 4 {
		return plan[:4]
	}
	return plan
}

// noDataPracticeTask là task placeholder khi user chưa có match history.
func noDataPracticeTask() PracticeTask {
	return PracticeTask{
		Day:       1,
		Focus:     "Setup dữ liệu",
		Duration:  "10 phút",
		Checklist: []string{"Login Riot", "Bấm Fetch report", "Kiểm tra report có evidence"},
		Evidence:  "Chưa có match history để coach.",
	}
}

// maintainFormTask là task fallback khi không Finding nào vượt ngưỡng — vẫn
// duy trì thói quen warmup + review.
func maintainFormTask(weakMap, mainAgent string) PracticeTask {
	return PracticeTask{
		Day:      1,
		Focus:    "Maintain form",
		Map:      weakMap,
		Agent:    mainAgent,
		Duration: "20 phút",
		Checklist: []string{
			"Warmup aim 10 phút",
			"Review 1 win và 1 loss gần nhất",
			"Ghi 1 adjustment cho map đang gặp nhiều nhất",
		},
		Evidence: "Không có finding vượt ngưỡng MVP trong sample hiện tại.",
	}
}
