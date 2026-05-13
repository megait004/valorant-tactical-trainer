package wailsiface

import (
	"context"
	"fmt"
	"time"

	matchdomain "valorant-tactical-trainer/internal/domain/match"
	"valorant-tactical-trainer/internal/infrastructure/storage"
	"valorant-tactical-trainer/internal/infrastructure/valorantapi"
)

type MatchService struct {
	store *storage.Store
}

func NewMatchService(store *storage.Store) *MatchService {
	return &MatchService{store: store}
}

type RefreshMatchesInput struct {
	PUUID  string `json:"puuid"`
	Region string `json:"region"`
	Size   string `json:"size"`
	APIKey string `json:"apiKey"`
}

type MatchDTO struct {
	MatchID      string `json:"matchId"`
	MapName      string `json:"mapName"`
	Mode         string `json:"mode"`
	Queue        string `json:"queue"`
	Region       string `json:"region"`
	GameStart    int64  `json:"gameStart"`
	RoundsPlayed int    `json:"roundsPlayed"`
	Agent        string `json:"agent"`
	Team         string `json:"team"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	Headshots    int    `json:"headshots"`
	DamageMade   int    `json:"damageMade"`
}

type RefreshMatchesResult struct {
	Matches  []MatchDTO `json:"matches"`
	Imported int        `json:"imported"`
	Cached   bool       `json:"cached"`
	Message  string     `json:"message"`
}

func (service *MatchService) RefreshMatches(input RefreshMatchesInput) (RefreshMatchesResult, error) {
	ctx := context.Background()
	if input.PUUID == "" {
		return RefreshMatchesResult{}, fmt.Errorf("InvalidPlayer: puuid is required")
	}

	apiKey := input.APIKey
	if apiKey == "" {
		if storedKey, ok, err := service.store.Setting(ctx, "api_key"); err == nil && ok {
			apiKey = storedKey
		}
	}

	cacheKey := fmt.Sprintf("matches:%s:%s:%s", input.Region, input.PUUID, input.Size)
	if _, ok, err := service.store.APICache(ctx, cacheKey); err != nil {
		return RefreshMatchesResult{}, err
	} else if ok {
		matches, err := service.store.MatchesForPlayer(ctx, input.PUUID)
		if err != nil {
			return RefreshMatchesResult{}, err
		}

		return RefreshMatchesResult{
			Matches:  toMatchDTOs(matches),
			Imported: 0,
			Cached:   true,
			Message:  "cached matches loaded",
		}, nil
	}

	client := valorantapi.NewClient(valorantapi.WithAPIKey(apiKey))
	summaries, rawPayload, err := client.MatchesByPUUID(ctx, input.PUUID, input.Region, input.Size)
	if err != nil {
		return RefreshMatchesResult{}, mapProviderError(err)
	}

	imported, err := service.store.SaveMatches(ctx, summaries)
	if err != nil {
		return RefreshMatchesResult{}, err
	}

	if err := service.store.SaveAPICache(ctx, cacheKey, "matches-by-puuid", rawPayload, 5*time.Minute); err != nil {
		return RefreshMatchesResult{}, err
	}

	matches, err := service.store.MatchesForPlayer(ctx, input.PUUID)
	if err != nil {
		return RefreshMatchesResult{}, err
	}

	return RefreshMatchesResult{
		Matches:  toMatchDTOs(matches),
		Imported: imported,
		Cached:   false,
		Message:  "matches refreshed",
	}, nil
}

func (service *MatchService) ListMatches(puuid string) ([]MatchDTO, error) {
	if puuid == "" {
		return []MatchDTO{}, nil
	}

	matches, err := service.store.MatchesForPlayer(context.Background(), puuid)
	if err != nil {
		return nil, err
	}

	return toMatchDTOs(matches), nil
}

func toMatchDTOs(summaries []matchdomain.Summary) []MatchDTO {
	matches := make([]MatchDTO, 0, len(summaries))
	for _, summary := range summaries {
		matches = append(matches, MatchDTO{
			MatchID:      summary.MatchID,
			MapName:      summary.MapName,
			Mode:         summary.Mode,
			Queue:        summary.Queue,
			Region:       summary.Region,
			GameStart:    summary.GameStart,
			RoundsPlayed: summary.RoundsPlayed,
			Agent:        summary.Agent,
			Team:         summary.Team,
			Kills:        summary.Kills,
			Deaths:       summary.Deaths,
			Assists:      summary.Assists,
			Headshots:    summary.Headshots,
			DamageMade:   summary.DamageMade,
		})
	}

	return matches
}
