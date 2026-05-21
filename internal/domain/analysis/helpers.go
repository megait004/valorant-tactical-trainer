package analysis

import "math"

// Helpers chung cho rule engine. Tách riêng để các file logic (metrics,
// breakdown, patterns) chỉ tập trung vào nghiệp vụ.

// ratio = value / total, an toàn với total = 0.
func ratio(value, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(value) / float64(total)
}

// round2 làm tròn 2 chữ số thập phân để JSON output gọn (vd 0.234 → 0.23).
func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

// severity trả "high" nếu value vượt ngưỡng highThreshold, không thì "medium".
// Dùng cho pattern muốn so sánh "khoảng cách so với baseline".
func severity(value, highThreshold float64) string {
	if value >= highThreshold {
		return "high"
	}
	return "medium"
}

// confidence chia 3 mức (high/medium/low) theo sample size so với target. Mỗi
// pattern truyền target khác nhau (vd: 60 round cho first-death rate, 5 trận
// cho KD).
func confidence(sampleSize, target int) string {
	if sampleSize >= target {
		return "high"
	}
	if sampleSize >= target/2 {
		return "medium"
	}
	return "low"
}

// roleDisplayName chuyển role kỹ thuật (controller, duelist, ...) thành tên
// tiếng Việt hiển thị UI. Trả nguyên role nếu không nằm trong danh sách.
func roleDisplayName(role string) string {
	switch role {
	case "duelist":
		return "Đối Đầu (Duelist)"
	case "initiator":
		return "Khởi Tranh (Initiator)"
	case "controller":
		return "Kiểm Soát (Controller)"
	case "sentinel":
		return "Hộ Vệ (Sentinel)"
	default:
		return role
	}
}
