package assistant

import (
	"fmt"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

// BuildAlertsFromReport convert Report (rule engine) thành queue Alert cho
// Live Assistant. Mỗi alert sẽ được hiển thị xoay vòng (cooldown 25s).
//
// Thứ tự ưu tiên trong queue:
//  1. Opener — mục tiêu phiên (lấy từ Recommendation đầu tiên)
//  2. Từng Finding rule engine phát hiện → alert tương ứng (mapping cố định)
//  3. Practice focus hôm nay (PracticePlan[0])
//  4. Map cần chú ý (weakest map)
//  5. Default maintain (chỉ khi queue rỗng — vd report chưa có data)
//
// Sau đó dedupe theo ID để tránh hiển thị 2 lần cùng alert.
func BuildAlertsFromReport(report analysis.Report) []Alert {
	alerts := make([]Alert, 0, 8)

	if len(report.Recommendations) > 0 {
		rec := report.Recommendations[0]
		alerts = append(alerts, Alert{
			ID:       "opener-recommendation",
			Title:    "Mục tiêu phiên",
			Message:  truncate(rec.Title+" — "+rec.Drill, 160),
			Severity: "medium",
			Source:   rec.FindingID,
		})
	}

	for _, finding := range report.Findings {
		if alert, ok := alertFromFinding(finding, report); ok {
			alerts = append(alerts, alert)
		}
	}

	if len(report.PracticePlan) > 0 {
		task := report.PracticePlan[0]
		checklist := ""
		if len(task.Checklist) > 0 {
			checklist = task.Checklist[0]
		}
		alerts = append(alerts, Alert{
			ID:       "practice-focus",
			Title:    "Bài tập hôm nay",
			Message:  truncate(task.Focus+": "+checklist, 160),
			Severity: "medium",
			Source:   "practice-plan",
		})
	}

	if report.Metrics.WeakestMap != "" {
		alerts = append(alerts, Alert{
			ID:       "map-awareness",
			Title:    "Map cần chú ý",
			Message:  fmt.Sprintf("%s — chơi default đã luyện, tránh pick fight solo.", report.Metrics.WeakestMap),
			Severity: "low",
			Source:   "map-gap",
		})
	}

	if len(alerts) == 0 {
		alerts = append(alerts, defaultMaintainAlert())
	}

	return dedupeAlerts(alerts)
}

// alertFromFinding map mỗi Finding ID sang alert content cụ thể. Mapping cố
// định cho 4 pattern hiện có + finding-no-data. Pattern mới chưa map ở đây
// sẽ bị skip (return ok=false) — UI vẫn fallback dùng Recommendation.
func alertFromFinding(finding analysis.Finding, report analysis.Report) (Alert, bool) {
	switch finding.ID {
	case "finding-first-death-non-duelist":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Round mới: smoke/info trước khi peek; nếu không phải Đối Đầu thì đừng entry sớm.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-low-headshot":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Pre-aim góc default, tap/burst — tránh wide swing khi crosshair chưa đặt đầu.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-low-survival-impact":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Chỉ nhận duel khi có trade path — sau death hỏi: có teammate trade trong 3 giây?",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-map-gap":
		mapName := report.Metrics.WeakestMap
		if mapName == "" {
			mapName = "map yếu"
		}
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  fmt.Sprintf("Ưu tiên protocol %s — default attack/defense, không hero pick mid.", mapName),
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	case "finding-no-data":
		return Alert{
			ID:       finding.ID,
			Title:    finding.Title,
			Message:  "Chưa có report thật — vào Settings bật consent và fetch Henrik trước khi luyện.",
			Severity: finding.Severity,
			Source:   finding.ID,
		}, true
	default:
		return Alert{}, false
	}
}

// defaultMaintainAlert hiển thị khi queue rỗng (vd user chưa fetch report).
func defaultMaintainAlert() Alert {
	return Alert{
		ID:       "default-maintain",
		Title:    "Giữ form",
		Message:  "Utility trước peek, gọi trade trước contact chính.",
		Severity: "low",
		Source:   "default",
	}
}

// dedupeAlerts loại bỏ alert trùng ID — giữ thứ tự xuất hiện đầu tiên.
func dedupeAlerts(alerts []Alert) []Alert {
	seen := map[string]bool{}
	out := make([]Alert, 0, len(alerts))
	for _, alert := range alerts {
		if seen[alert.ID] {
			continue
		}
		seen[alert.ID] = true
		out = append(out, alert)
	}
	return out
}
