package wailsiface

import (
	"context"
	"fmt"
	"time"

	"valorant-tactical-trainer/internal/domain/rank"
	"valorant-tactical-trainer/internal/infrastructure/storage"
	"valorant-tactical-trainer/internal/infrastructure/valorantapi"
)

type RankService struct {
	store *storage.Store
}

func NewRankService(store *storage.Store) *RankService {
	return &RankService{store: store}
}

type RefreshRankInput struct {
	PUUID  string `json:"puuid"`
	Region string `json:"region"`
	APIKey string `json:"apiKey"`
}

type RankDTO struct {
	Tier            int    `json:"tier"`
	TierName        string `json:"tierName"`
	RankingInTier   int    `json:"rankingInTier"`
	MMRChangeToLast int    `json:"mmrChangeToLast"`
	Elo             int    `json:"elo"`
	SeasonID        string `json:"seasonId"`
	Region          string `json:"region"`
	FetchedAt       string `json:"fetchedAt"`
}

type RefreshRankResult struct {
	Rank    RankDTO `json:"rank"`
	Cached  bool    `json:"cached"`
	Message string  `json:"message"`
}

func (service *RankService) RefreshRank(input RefreshRankInput) (RefreshRankResult, error) {
	ctx := context.Background()
	if input.PUUID == "" {
		return RefreshRankResult{}, fmt.Errorf("InvalidPlayer: puuid is required")
	}

	apiKey := input.APIKey
	if apiKey == "" {
		if storedKey, ok, err := service.store.Setting(ctx, "api_key"); err == nil && ok {
			apiKey = storedKey
		}
	}

	cacheKey := fmt.Sprintf("mmr:%s:%s", input.Region, input.PUUID)
	if _, ok, err := service.store.APICache(ctx, cacheKey); err != nil {
		return RefreshRankResult{}, err
	} else if ok {
		snapshot, hasRank, err := service.store.LatestRankSnapshot(ctx, input.PUUID)
		if err != nil {
			return RefreshRankResult{}, err
		}
		if hasRank {
			return RefreshRankResult{Rank: toRankDTO(snapshot), Cached: true, Message: "cached rank loaded"}, nil
		}
	}

	client := valorantapi.NewBasicClient(valorantapi.WithAPIKey(apiKey))
	snapshot, rawPayload, err := client.MMRByPUUID(ctx, input.PUUID, input.Region)
	if err != nil {
		return RefreshRankResult{}, mapProviderError(err)
	}

	if err := service.store.SaveRankSnapshot(ctx, snapshot); err != nil {
		return RefreshRankResult{}, err
	}
	if err := service.store.SaveAPICache(ctx, cacheKey, "mmr-by-puuid", rawPayload, 10*time.Minute); err != nil {
		return RefreshRankResult{}, err
	}

	return RefreshRankResult{Rank: toRankDTO(snapshot), Cached: false, Message: "rank refreshed"}, nil
}

func (service *RankService) LatestRank(puuid string) (*RankDTO, error) {
	if puuid == "" {
		return nil, nil
	}

	snapshot, ok, err := service.store.LatestRankSnapshot(context.Background(), puuid)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	dto := toRankDTO(snapshot)
	return &dto, nil
}

func toRankDTO(snapshot rank.Snapshot) RankDTO {
	return RankDTO{
		Tier:            snapshot.Tier,
		TierName:        snapshot.TierName,
		RankingInTier:   snapshot.RankingInTier,
		MMRChangeToLast: snapshot.MMRChangeToLast,
		Elo:             snapshot.Elo,
		SeasonID:        snapshot.SeasonID,
		Region:          snapshot.Region,
		FetchedAt:       snapshot.FetchedAt.Format(time.RFC3339),
	}
}
