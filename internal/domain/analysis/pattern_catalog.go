package analysis

import "fmt"

// File này khai báo từng FindingPattern cụ thể. Mỗi pattern là một function
// constructor để có thể closure các threshold riêng và dễ unit test.
//
// Tất cả "magic number" ngưỡng kích hoạt được khai báo ngay đầu function (const
// local) để dễ tinker/A-B test không phải tìm khắp file.

func patternFirstDeathNonDuelist() FindingPattern {
	const (
		threshold     = 0.18 // first death rate ≥ 18% là cảnh báo
		highThreshold = 0.24 // ≥ 24% là severity high
		baselineFD    = 0.18
	)
	return FindingPattern{
		ID:           "finding-first-death-non-duelist",
		DefaultTitle: "First Death cao so với vai trò không phải Đối Đầu",
		Match: func(snapshot PlayerSnapshot, m MetricSummary) (Evidence, bool) {
			if m.FirstDeathRate < threshold || m.PrimaryRoleObserved == "duelist" {
				return Evidence{}, false
			}
			return Evidence{
				MatchIDs:           collectMatchIDs(snapshot.RecentMatches),
				Metric:             "first_death_rate",
				Value:              m.FirstDeathRate,
				SampleSize:         m.Rounds,
				ComparisonBaseline: baselineFD,
			}, true
		},
		DefaultDetail: func(m MetricSummary) string {
			return fmt.Sprintf("Tỉ lệ chết đầu round là %.0f%% trong khi vai trò chính là %s.",
				m.FirstDeathRate*100, roleDisplayName(m.PrimaryRoleObserved))
		},
		SeverityFn:   func(m MetricSummary) string { return severity(m.FirstDeathRate, highThreshold) },
		ConfidenceFn: func(m MetricSummary) string { return confidence(m.Rounds, 60) },
		DefaultRecommendation: func(_ MetricSummary) Recommendation {
			return Recommendation{
				ID:        "rec-survive-contact",
				FindingID: "finding-first-death-non-duelist",
				Title:     "Giảm chết sớm bằng utility-before-peek drill",
				Reason:    "Kiểm Soát/Hộ Vệ cần sống tới mid-late round để giữ smoke/trap và call rotate.",
				Drill:     "Custom map yếu, 10 round defense: trước mỗi peek phải đặt smoke/trap/info hoặc gọi teammate trade.",
				Cadence:   "15 phút/ngày trong 5 ngày",
			}
		},
		DefaultPracticeTask: func(m MetricSummary, weakMap, mainAgent string, day int) PracticeTask {
			return PracticeTask{
				Day:      day,
				Focus:    "Survive first contact",
				Map:      weakMap,
				Agent:    mainAgent,
				Duration: "20 phút",
				Checklist: []string{
					"Không peek đầu round nếu chưa có utility/info",
					"Gọi teammate trade trước mọi contact chính",
					"Đánh dấu mọi death trước 1:20 trong 5 trận tiếp theo",
				},
				Evidence: fmt.Sprintf("First death rate %.0f%% trên %d round.", m.FirstDeathRate*100, m.Rounds),
			}
		},
	}
}

func patternLowHeadshot() FindingPattern {
	const (
		threshold     = 18.0 // HS% < 18 là cảnh báo
		highThreshold = 5.0  // gap ≥ 5 là severity high
		baselineHS    = 18.0
	)
	return FindingPattern{
		ID:           "finding-low-headshot",
		DefaultTitle: "Headshot thấp, khả năng crosshair placement chưa ổn định",
		Match: func(snapshot PlayerSnapshot, m MetricSummary) (Evidence, bool) {
			if m.HeadshotPercent >= threshold {
				return Evidence{}, false
			}
			return Evidence{
				MatchIDs:           collectMatchIDs(snapshot.RecentMatches),
				Metric:             "headshot_percent",
				Value:              m.HeadshotPercent,
				SampleSize:         m.Rounds,
				ComparisonBaseline: baselineHS,
			}, true
		},
		DefaultDetail: func(m MetricSummary) string {
			return fmt.Sprintf("HS trung bình %.1f%%, thấp hơn baseline MVP 18%%.", m.HeadshotPercent)
		},
		SeverityFn:   func(m MetricSummary) string { return severity(threshold-m.HeadshotPercent, highThreshold) },
		ConfidenceFn: func(m MetricSummary) string { return confidence(m.Rounds, 50) },
		DefaultRecommendation: func(_ MetricSummary) Recommendation {
			return Recommendation{
				ID:        "rec-crosshair-placement",
				FindingID: "finding-low-headshot",
				Title:     "Crosshair placement theo tuyến di chuyển",
				Reason:    "HS thấp thường tới từ tâm đặt ngang ngực/chân hoặc pre-aim chưa theo góc phổ biến.",
				Drill:     "Deathmatch 2 game tap Sheriff/Vandal, không spray; sau đó custom 5 phút pre-aim các góc default trên map yếu.",
				Cadence:   "25 phút/ngày, đo lại sau 10 trận",
			}
		},
		DefaultPracticeTask: func(m MetricSummary, weakMap, mainAgent string, day int) PracticeTask {
			return PracticeTask{
				Day:      day,
				Focus:    "Crosshair placement",
				Map:      weakMap,
				Agent:    mainAgent,
				Duration: "20 phút",
				Checklist: []string{
					"Deathmatch 1 game chỉ tap/burst, không spray dài",
					"Custom pre-aim 15 góc default trên map yếu",
					"Review 3 death do tâm đặt thấp hoặc swing thiếu pre-aim",
				},
				Evidence: fmt.Sprintf("HS trung bình %.1f%%, baseline 18%%.", m.HeadshotPercent),
			}
		},
	}
}

func patternLowSurvivalImpact() FindingPattern {
	const (
		threshold     = 0.9 // KD < 0.9 là cảnh báo
		highThreshold = 0.2 // gap ≥ 0.2 là severity high
		baselineKD    = 0.9
	)
	return FindingPattern{
		ID:           "finding-low-survival-impact",
		DefaultTitle: "Impact giao tranh thấp do trade/survival chưa ổn",
		Match: func(snapshot PlayerSnapshot, m MetricSummary) (Evidence, bool) {
			if m.KD >= threshold {
				return Evidence{}, false
			}
			return Evidence{
				MatchIDs:           collectMatchIDs(snapshot.RecentMatches),
				Metric:             "kd",
				Value:              m.KD,
				SampleSize:         m.Matches,
				ComparisonBaseline: baselineKD,
			}, true
		},
		DefaultDetail: func(m MetricSummary) string {
			return fmt.Sprintf("K/D hiện tại %.2f, cần giảm death vô ích trước khi tăng playmaking.", m.KD)
		},
		SeverityFn:   func(m MetricSummary) string { return severity(threshold-m.KD, highThreshold) },
		ConfidenceFn: func(m MetricSummary) string { return confidence(m.Matches, 5) },
		DefaultRecommendation: func(_ MetricSummary) Recommendation {
			return Recommendation{
				ID:        "rec-trade-discipline",
				FindingID: "finding-low-survival-impact",
				Title:     "Trade discipline và second-contact drill",
				Reason:    "K/D thấp cần sửa death không được trade trước khi ép aim duel khó.",
				Drill:     "Trong 5 trận tiếp theo, ghi lại mọi death không có teammate trong 3 giây trade range; mục tiêu giảm 30% death cô lập.",
				Cadence:   "Theo dõi 5 trận ranked/scrim",
			}
		},
		DefaultPracticeTask: func(m MetricSummary, weakMap, mainAgent string, day int) PracticeTask {
			return PracticeTask{
				Day:      day,
				Focus:    "Trade discipline",
				Map:      weakMap,
				Agent:    mainAgent,
				Duration: "Theo dõi 5 trận",
				Checklist: []string{
					"Chỉ nhận duel nếu có trade path hoặc exit path",
					"Sau mỗi death ghi có teammate trade trong 3 giây không",
					"Mục tiêu giảm death cô lập 30%",
				},
				Evidence: fmt.Sprintf("K/D %.2f, baseline 0.90.", m.KD),
			}
		},
	}
}

func patternMapGap() FindingPattern {
	const (
		threshold     = 0.35 // map win rate ≤ 35% trong sample đủ lớn
		baselineWR    = 0.5
		minSampleHigh = 3
	)
	return FindingPattern{
		ID:           "finding-map-gap",
		DefaultTitle: "Có map gap cần ưu tiên luyện riêng",
		Match: func(snapshot PlayerSnapshot, m MetricSummary) (Evidence, bool) {
			if m.WeakestMap == "" || m.WeakestMapWinRate > threshold {
				return Evidence{}, false
			}
			return Evidence{
				MatchIDs:           matchesForMap(snapshot.RecentMatches, m.WeakestMap),
				Map:                m.WeakestMap,
				Metric:             "map_win_rate",
				Value:              m.WeakestMapWinRate,
				SampleSize:         m.WeakestMapSample,
				ComparisonBaseline: baselineWR,
			}, true
		},
		DefaultDetail: func(m MetricSummary) string {
			return fmt.Sprintf("Map %s có win rate %.0f%% trong %d trận gần đây.",
				m.WeakestMap, m.WeakestMapWinRate*100, m.WeakestMapSample)
		},
		SeverityFn:   func(_ MetricSummary) string { return "medium" },
		ConfidenceFn: func(m MetricSummary) string { return confidence(m.WeakestMapSample, minSampleHigh) },
		DefaultRecommendation: func(_ MetricSummary) Recommendation {
			return Recommendation{
				ID:        "rec-map-playbook",
				FindingID: "finding-map-gap",
				Title:     "Tạo mini-playbook cho map yếu",
				Reason:    "Win rate map thấp cần sửa route, setup và mid-round protocol thay vì luyện aim chung chung.",
				Drill:     "Viết 2 default attack, 2 setup defense và 1 retake protocol cho map này; chạy custom walkthrough 20 phút.",
				Cadence:   "2 buổi/tuần",
			}
		},
		DefaultPracticeTask: func(m MetricSummary, weakMap, mainAgent string, day int) PracticeTask {
			return PracticeTask{
				Day:      day,
				Focus:    "Map gap protocol",
				Map:      weakMap,
				Agent:    mainAgent,
				Duration: "25 phút",
				Checklist: []string{
					"Viết 2 default route attack cho map yếu",
					"Viết 2 setup defense giữ site/choke chính",
					"Chạy custom walkthrough tới khi timing smoke/rotate ổn định",
				},
				Evidence: fmt.Sprintf("%s win rate %.0f%% trên %d trận.",
					m.WeakestMap, m.WeakestMapWinRate*100, m.WeakestMapSample),
			}
		},
	}
}
