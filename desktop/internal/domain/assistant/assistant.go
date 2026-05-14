package assistant

import "strings"

type Query struct {
	MapName         string
	Agent           string
	Side            string
	Phase           string
	Credits         int
	PreviousOutcome string
}

type TacticalCard struct {
	ID          string
	MapName     string
	Agent       string
	Side        string
	Phase       string
	Category    string
	Title       string
	Summary     string
	Action      string
	Priority    int
	SafetyNotes string
}

type EconomyAdvice struct {
	Plan         string
	Summary      string
	BuyThreshold int
	NextRoundMin int
	Reminder     string
}

type Result struct {
	Query         Query
	Cards         []TacticalCard
	EconomyAdvice EconomyAdvice
	SafetyNotes   []string
}

func NormalizeText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func RecommendEconomy(query Query) EconomyAdvice {
	credits := query.Credits
	outcome := NormalizeText(query.PreviousOutcome)

	if credits < 2000 {
		return EconomyAdvice{
			Plan:         "Eco",
			Summary:      "Tiết kiệm round này để round sau có rifle + giáp.",
			BuyThreshold: 3900,
			NextRoundMin: 3900,
			Reminder:     "Mua pistol/utility rẻ, tránh xuống dưới ngưỡng full buy round sau.",
		}
	}
	if credits < 3900 && outcome == "loss" {
		return EconomyAdvice{
			Plan:         "Light / Half Buy",
			Summary:      "Đội vừa thua và tiền chưa đủ ổn định, ưu tiên giữ economy.",
			BuyThreshold: 3900,
			NextRoundMin: 3900,
			Reminder:     "Chỉ mua nếu team thống nhất, giữ lại tiền cho Vandal/Phantom + giáp sau đó.",
		}
	}
	if credits >= 3900 {
		return EconomyAdvice{
			Plan:         "Full Buy",
			Summary:      "Đủ tiền mua rifle, giáp và utility chính.",
			BuyThreshold: 3900,
			NextRoundMin: 0,
			Reminder:     "Đồng bộ timing dùng utility, đừng dry peek khi team còn kỹ năng mở site.",
		}
	}

	return EconomyAdvice{
		Plan:         "Force Buy",
		Summary:      "Có thể mua ép nếu team cần phá nhịp hoặc tránh reset momentum.",
		BuyThreshold: 3000,
		NextRoundMin: 0,
		Reminder:     "Chỉ force khi cả team cùng mua; nếu không thì round này dễ mất economy kép.",
	}
}

func SeedCards() []TacticalCard {
	return []TacticalCard{
		{ID: "ascent-prematch-meta", MapName: "ascent", Agent: "", Side: "both", Phase: "prematch", Category: "composition", Title: "Ascent ưu tiên kiểm soát Mid", Summary: "Đội hình ổn thường có initiator mở Mid, controller giữ smoke Cat/Market và sentinel khóa flank.", Action: "Nếu thiếu controller hoặc initiator, nhắc team cân bằng role trước khi lock agent.", Priority: 90, SafetyNotes: "Gợi ý local, không đọc dữ liệu game."},
		{ID: "ascent-attack-default", MapName: "ascent", Agent: "", Side: "attack", Phase: "ingame", Category: "default-strat", Title: "Default chiếm Mid trước khi hit site", Summary: "Ascent phạt nặng đội bỏ Mid. Lấy Top Mid/Catwalk giúp split A hoặc B an toàn hơn.", Action: "Một người giữ lurk, initiator clear close Mid, controller smoke Market hoặc Tree theo hướng hit.", Priority: 85, SafetyNotes: "Không overlay tự động, chỉ là tab tra cứu."},
		{ID: "ascent-defense-default", MapName: "ascent", Agent: "", Side: "defense", Phase: "ingame", Category: "default-strat", Title: "Setup 2-1-2 có người contest Mid", Summary: "Giữ một người có utility kiểm soát Mid để lấy info sớm, tránh bị split site miễn phí.", Action: "Smoke/flash deny Top Mid đầu round, sau đó rotate theo sound/info thay vì stack sớm.", Priority: 82, SafetyNotes: "Không che HUD, dùng cửa sổ app riêng."},
		{ID: "bind-attack-hookah", MapName: "bind", Agent: "", Side: "attack", Phase: "ingame", Category: "default-strat", Title: "Chiếm Hookah cần flash/stun trước", Summary: "Hookah có nhiều góc close, dry peek dễ mất người đầu round.", Action: "Dùng flash/stun vào Window, người thứ hai trade ngay sau entry, smoke CT/Elbow khi hit B.", Priority: 88, SafetyNotes: "Không dùng memory reading."},
		{ID: "bind-defense-a-short", MapName: "bind", Agent: "", Side: "defense", Phase: "ingame", Category: "crosshair", Title: "Giữ A Short bằng crosshair ngang đầu", Summary: "A Short thường có swing rộng sau utility. Kê sát góc dễ bị prefire.", Action: "Đứng lệch off-angle, kê ngang đầu tại mép tường, fallback sau 1 kill hoặc khi mất smoke.", Priority: 74, SafetyNotes: "Card mô tả text thay cho ảnh/video để tránh asset bản quyền."},
		{ID: "haven-defense-c-long", MapName: "haven", Agent: "", Side: "defense", Phase: "ingame", Category: "default-strat", Title: "Phòng thủ C Long bằng utility chặn nhịp", Summary: "Haven có 3 site nên mất C Long sớm làm rotate rất khó.", Action: "Dùng smoke/molly/stun ở đầu C Long, lấy info rồi lùi giữ site chéo với Garage.", Priority: 86, SafetyNotes: "Không tương tác tiến trình VALORANT."},
		{ID: "haven-attack-garage", MapName: "haven", Agent: "", Side: "attack", Phase: "ingame", Category: "default-strat", Title: "Garage là trục xoay round", Summary: "Kiểm soát Garage mở đường split C hoặc pressure B, ép defense chia người.", Action: "Clear close bằng drone/flash, giữ window sau khi chiếm, không rush C nếu chưa cắt rotate.", Priority: 80, SafetyNotes: "Manual lookup only."},
		{ID: "sova-ascent-lineup", MapName: "ascent", Agent: "sova", Side: "attack", Phase: "ingame", Category: "lineup", Title: "Sova mở A bằng recon Tree/Generator", Summary: "Recon sớm giúp entry biết vị trí close A và ép defense phá dart.", Action: "Đứng A Lobby an toàn, bắn recon qua mái vào khu Tree/Generator trước khi team dash/swing.", Priority: 91, SafetyNotes: "Lineup text, không tự động aim."},
		{ID: "viper-bind-lineup", MapName: "bind", Agent: "viper", Side: "attack", Phase: "ingame", Category: "lineup", Title: "Viper giữ post-plant bằng Snake Bite", Summary: "Bind có post-plant mạnh nếu team sống đủ lâu và giữ lineups từ xa.", Action: "Sau plant, lùi về vị trí an toàn, chỉ dùng lineup khi nghe defuse hoặc team call rõ.", Priority: 78, SafetyNotes: "Không bắn kỹ năng hộ người chơi."},
		{ID: "brimstone-bind-smoke", MapName: "bind", Agent: "brimstone", Side: "attack", Phase: "ingame", Category: "lineup", Title: "Brimstone smoke CT/Elbow khi hit B", Summary: "B site Bind cần cắt góc CT và Elbow để entry không bị crossfire.", Action: "Smoke CT + Elbow, stim entry sau flash/stun, giữ molly cho defuse hoặc retake choke.", Priority: 84, SafetyNotes: "Chỉ gợi ý chiến thuật."},
	}
}
