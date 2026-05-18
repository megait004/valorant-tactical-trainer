package wailsiface

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
	"valorant-tactical-trainer/desktop/internal/infrastructure/llm"
	"valorant-tactical-trainer/desktop/internal/infrastructure/localstore"
)

// ChatService cung cấp endpoint cho icon bot AI ở UI: nhận message từ user,
// ghép context (player + metrics + findings + breakdown) từ report cuối cùng
// và gọi OpenAI/Claude qua llm.Coach. Lịch sử hội thoại giữ in-memory để LLM
// có context turn-by-turn.
type ChatService struct {
	coach       *llm.Coach
	reportStore *localstore.ReportStore

	mu       sync.Mutex
	messages []llm.ChatMessage
}

// ChatMessage là DTO trả về frontend.
type ChatMessage struct {
	Role      string `json:"role"` // "user" | "assistant"
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// ChatState là state hiện tại của cuộc hội thoại (cho UI render).
type ChatState struct {
	Available bool          `json:"available"`
	Message   string        `json:"message"`
	History   []ChatMessage `json:"history"`
}

// IsAvailable cho UI biết coach có hoạt động không (có LLM_API_KEY trong .env).
func (s *ChatService) IsAvailable() bool {
	return s != nil && s.coach != nil
}

// GetState trả lịch sử hiện tại để UI restore khi mở panel.
func (s *ChatService) GetState() ChatState {
	state := ChatState{Available: s.IsAvailable()}
	if !state.Available {
		state.Message = "Chưa có LLM_API_KEY trong .env — bot AI tạm tắt."
		return state
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state.History = toDTOHistory(s.messages)
	return state
}

// Reset xoá lịch sử hội thoại.
func (s *ChatService) Reset() ChatState {
	s.mu.Lock()
	s.messages = nil
	s.mu.Unlock()
	return s.GetState()
}

// SendMessage gửi 1 message từ user, ghép context report + history rồi gọi LLM.
// Trả về state mới (đã có cả user message lẫn assistant reply).
func (s *ChatService) SendMessage(message string) (ChatState, error) {
	if s == nil || s.coach == nil {
		return ChatState{Available: false, Message: "Bot AI chưa cấu hình LLM_API_KEY."}, errors.New("coach unavailable")
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return s.GetState(), errors.New("message rỗng")
	}

	s.mu.Lock()
	s.messages = append(s.messages, llm.ChatMessage{Role: "user", Content: message})
	// Giữ tối đa 16 turn (8 round QA) để control token.
	if len(s.messages) > 16 {
		s.messages = s.messages[len(s.messages)-16:]
	}
	historyCopy := append([]llm.ChatMessage(nil), s.messages...)
	s.mu.Unlock()

	system := s.buildSystemPrompt()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	reply, err := s.coach.Complete(ctx, system, historyCopy, false)
	if err != nil {
		// Rollback message user vì call lỗi → user có thể retry mà không bị
		// double trong history.
		s.mu.Lock()
		if n := len(s.messages); n > 0 && s.messages[n-1].Role == "user" {
			s.messages = s.messages[:n-1]
		}
		s.mu.Unlock()
		return s.GetState(), fmt.Errorf("LLM lỗi: %w", err)
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		reply = "Mình chưa có câu trả lời rõ ràng — thử hỏi lại cụ thể hơn nhé."
	}

	s.mu.Lock()
	s.messages = append(s.messages, llm.ChatMessage{Role: "assistant", Content: reply})
	s.mu.Unlock()

	return s.GetState(), nil
}

// buildSystemPrompt ghép context từ last_report.json (nếu có) thành system
// prompt giàu thông tin. Mỗi lần SendMessage là load lại để lấy report mới
// nhất nếu user vừa fetch.
func (s *ChatService) buildSystemPrompt() string {
	base := `Bạn là HLV Valorant cá nhân hoá cho 1 player cụ thể. Người dùng đặt câu hỏi về kỹ năng, map, agent, drill luyện tập của họ.

Nguyên tắc:
- Trả lời tiếng Việt tự nhiên, thẳng vào trọng tâm, ngắn gọn.
- Bám sát số liệu trong CONTEXT bên dưới — không bịa metric.
- Nếu user hỏi điều CONTEXT không có (vd map họ chưa chơi), trả lời theo hiểu biết Valorant chung và nói rõ là chưa có data.
- Khi đưa drill, nêu rõ thời lượng, map/agent gợi ý và cách đo tiến bộ.
- Không dùng emoji quá nhiều, không format markdown phức tạp; bullet đơn giản OK.`

	stored, ok, err := s.reportStore.LoadLastReport()
	if err != nil || !ok {
		return base + "\n\nCONTEXT: Player chưa có report — khuyên user fetch report trước trong tab Coach để có dữ liệu phân tích."
	}
	report := stored.Report
	return base + "\n\n" + summarizeReportForChat(report)
}

func summarizeReportForChat(r analysis.Report) string {
	var b strings.Builder
	b.WriteString("CONTEXT (report gần nhất từ Henrik/Riot):\n")
	fmt.Fprintf(&b, "- Player: %s#%s · region %s · vai trò chính: %s\n",
		r.Player.Name, r.Player.Tagline, r.Player.Region, r.Metrics.PrimaryRoleObserved)
	fmt.Fprintf(&b, "- Tổng quan: %d trận / %d round · K/D %.2f · KDA %.2f · HS%% %.1f · WinRate %.0f%%\n",
		r.Metrics.Matches, r.Metrics.Rounds, r.Metrics.KD, r.Metrics.KDA,
		r.Metrics.HeadshotPercent, r.Metrics.WinRate*100)
	fmt.Fprintf(&b, "- First Blood %.0f%% · First Death %.0f%%\n",
		r.Metrics.FirstBloodRate*100, r.Metrics.FirstDeathRate*100)
	if r.Metrics.WeakestMap != "" {
		fmt.Fprintf(&b, "- Map yếu: %s (winrate %.0f%% trên %d trận)\n",
			r.Metrics.WeakestMap, r.Metrics.WeakestMapWinRate*100, r.Metrics.WeakestMapSample)
	}

	if len(r.MapBreakdown) > 0 {
		b.WriteString("- Map breakdown:\n")
		for i, m := range r.MapBreakdown {
			if i >= 5 {
				break
			}
			fmt.Fprintf(&b, "    · %s: %d trận, K/D %.2f, WR %.0f%%, HS %.1f%%\n",
				m.Name, m.Matches, m.KD, m.WinRate*100, m.HeadshotPercent)
		}
	}
	if len(r.AgentBreakdown) > 0 {
		b.WriteString("- Agent breakdown:\n")
		for i, a := range r.AgentBreakdown {
			if i >= 5 {
				break
			}
			fmt.Fprintf(&b, "    · %s: %d trận, K/D %.2f, WR %.0f%%, HS %.1f%%\n",
				a.Name, a.Matches, a.KD, a.WinRate*100, a.HeadshotPercent)
		}
	}
	if len(r.Findings) > 0 {
		b.WriteString("- Findings (rule engine + LLM):\n")
		for _, f := range r.Findings {
			fmt.Fprintf(&b, "    · [%s/%s] %s — %s\n", f.Severity, f.Confidence, f.Title, f.Detail)
		}
	}
	if len(r.Player.RecentMatches) > 0 {
		b.WriteString("- Trận gần nhất:\n")
		for i, m := range r.Player.RecentMatches {
			if i >= 5 {
				break
			}
			result := "thua"
			if m.Won {
				result = "thắng"
			}
			fmt.Fprintf(&b, "    · %s / %s: %d-%d-%d · %d round · %s\n",
				m.Map, m.Agent, m.Kills, m.Deaths, m.Assists, m.RoundsPlayed, result)
		}
	}
	return b.String()
}

func toDTOHistory(messages []llm.ChatMessage) []ChatMessage {
	out := make([]ChatMessage, 0, len(messages))
	now := time.Now().Format(time.RFC3339)
	for _, m := range messages {
		out = append(out, ChatMessage{Role: m.Role, Content: m.Content, CreatedAt: now})
	}
	return out
}
