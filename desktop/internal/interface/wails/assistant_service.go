package wailsiface

import (
	"context"
	"fmt"

	assistantdomain "valorant-tactical-trainer/internal/domain/assistant"
	"valorant-tactical-trainer/internal/infrastructure/storage"
)

type AssistantService struct {
	store *storage.Store
}

type AssistantQueryInput struct {
	MapName         string `json:"mapName"`
	Agent           string `json:"agent"`
	Side            string `json:"side"`
	Phase           string `json:"phase"`
	Credits         int    `json:"credits"`
	PreviousOutcome string `json:"previousOutcome"`
}

type TacticalCardDTO struct {
	ID          string `json:"id"`
	MapName     string `json:"mapName"`
	Agent       string `json:"agent"`
	Side        string `json:"side"`
	Phase       string `json:"phase"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Action      string `json:"action"`
	Priority    int    `json:"priority"`
	SafetyNotes string `json:"safetyNotes"`
}

type EconomyAdviceDTO struct {
	Plan         string `json:"plan"`
	Summary      string `json:"summary"`
	BuyThreshold int    `json:"buyThreshold"`
	NextRoundMin int    `json:"nextRoundMin"`
	Reminder     string `json:"reminder"`
}

type AssistantResultDTO struct {
	Cards         []TacticalCardDTO `json:"cards"`
	EconomyAdvice EconomyAdviceDTO  `json:"economyAdvice"`
	SafetyNotes   []string          `json:"safetyNotes"`
}

func NewAssistantService(store *storage.Store) *AssistantService {
	return &AssistantService{store: store}
}

func (service *AssistantService) QueryAssistant(input AssistantQueryInput) (AssistantResultDTO, error) {
	query := assistantdomain.Query{
		MapName:         input.MapName,
		Agent:           input.Agent,
		Side:            input.Side,
		Phase:           input.Phase,
		Credits:         input.Credits,
		PreviousOutcome: input.PreviousOutcome,
	}
	cards, err := service.store.TacticalCards(context.Background(), query)
	if err != nil {
		return AssistantResultDTO{}, fmt.Errorf("err querying assistant: %w", err)
	}

	return AssistantResultDTO{
		Cards:         toTacticalCardDTOs(cards),
		EconomyAdvice: toEconomyAdviceDTO(assistantdomain.RecommendEconomy(query)),
		SafetyNotes: []string{
			"Không đọc memory VALORANT.exe hoặc inject vào game.",
			"Dữ liệu hiện là tactical cards local, user chọn map/agent thủ công.",
			"Overlay/hotkey sẽ là story riêng để test UX và anti-cheat safety.",
		},
	}, nil
}

func toTacticalCardDTOs(cards []assistantdomain.TacticalCard) []TacticalCardDTO {
	result := make([]TacticalCardDTO, 0, len(cards))
	for _, card := range cards {
		result = append(result, TacticalCardDTO{
			ID:          card.ID,
			MapName:     card.MapName,
			Agent:       card.Agent,
			Side:        card.Side,
			Phase:       card.Phase,
			Category:    card.Category,
			Title:       card.Title,
			Summary:     card.Summary,
			Action:      card.Action,
			Priority:    card.Priority,
			SafetyNotes: card.SafetyNotes,
		})
	}

	return result
}

func toEconomyAdviceDTO(advice assistantdomain.EconomyAdvice) EconomyAdviceDTO {
	return EconomyAdviceDTO{
		Plan:         advice.Plan,
		Summary:      advice.Summary,
		BuyThreshold: advice.BuyThreshold,
		NextRoundMin: advice.NextRoundMin,
		Reminder:     advice.Reminder,
	}
}
