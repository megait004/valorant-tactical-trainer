package wailsiface

import (
	"context"
	"errors"
	"fmt"
	"time"

	"valorant-tactical-trainer/internal/domain/player"
	"valorant-tactical-trainer/internal/infrastructure/storage"
	"valorant-tactical-trainer/internal/infrastructure/valorantapi"
)

type PlayerService struct {
	store *storage.Store
}

func NewPlayerService(store *storage.Store) *PlayerService {
	return &PlayerService{store: store}
}

type LookupPlayerInput struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Region  string `json:"region"`
	Consent bool   `json:"consent"`
	APIKey  string `json:"apiKey"`
}

type PlayerDTO struct {
	PUUID        string `json:"puuid"`
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Region       string `json:"region"`
	AccountLevel int    `json:"accountLevel"`
	CardSmall    string `json:"cardSmall"`
	CardLarge    string `json:"cardLarge"`
	LastUpdate   string `json:"lastUpdate"`
}

type LookupPlayerResult struct {
	Player         PlayerDTO `json:"player"`
	Provider       string    `json:"provider"`
	ConsentVersion string    `json:"consentVersion"`
	Message        string    `json:"message"`
}

func (service *PlayerService) GetCurrentPlayer() (*PlayerDTO, error) {
	ctx := context.Background()
	account, ok, err := service.store.CurrentPlayer(ctx)
	if err != nil {
		return nil, fmt.Errorf("err reading current player: %w", err)
	}
	if !ok {
		return nil, nil
	}

	dto := toPlayerDTO(account)
	return &dto, nil
}

func (service *PlayerService) LookupPlayer(input LookupPlayerInput) (LookupPlayerResult, error) {
	ctx := context.Background()
	name := player.NormalizeName(input.Name)
	tag := player.NormalizeTag(input.Tag)
	region := player.NormalizeRegion(input.Region)

	if !input.Consent {
		return LookupPlayerResult{}, errors.New("ConsentRequired: confirm consent before lookup")
	}

	if !player.IsValidLookup(name, tag) {
		return LookupPlayerResult{}, errors.New("InvalidPlayer: name and tag are required")
	}

	if input.APIKey != "" {
		if err := service.store.SaveSetting(ctx, "api_key", input.APIKey); err != nil {
			return LookupPlayerResult{}, err
		}
	}

	client := valorantapi.NewBasicClient(valorantapi.WithAPIKey(input.APIKey))
	account, err := client.LookupAccount(ctx, name, tag)
	if err != nil {
		return LookupPlayerResult{}, mapProviderError(err)
	}

	if account.Region == "" {
		account.Region = region
	}

	consent := player.Consent{
		PlayerPUUID:    account.PUUID,
		Name:           name,
		Tag:            tag,
		Region:         region,
		Provider:       valorantapi.ProviderName(),
		ConsentVersion: player.ConsentVersion,
		ConsentedAt:    time.Now().UTC(),
	}

	if err := service.store.SavePlayerWithConsent(ctx, account, consent); err != nil {
		return LookupPlayerResult{}, err
	}

	return LookupPlayerResult{
		Player:         toPlayerDTO(account),
		Provider:       valorantapi.ProviderName(),
		ConsentVersion: player.ConsentVersion,
		Message:        "data received",
	}, nil
}

func toPlayerDTO(account player.Account) PlayerDTO {
	return PlayerDTO{
		PUUID:        account.PUUID,
		Name:         account.Name,
		Tag:          account.Tag,
		Region:       account.Region,
		AccountLevel: account.AccountLevel,
		CardSmall:    account.CardSmall,
		CardLarge:    account.CardLarge,
		LastUpdate:   account.LastUpdate,
	}
}

func mapProviderError(err error) error {
	switch {
	case errors.Is(err, valorantapi.ErrRateLimited):
		return errors.New("RateLimited: provider rate limit reached")
	case errors.Is(err, valorantapi.ErrUnauthorizedAPIKey):
		return errors.New("UnauthorizedApiKey: API key was rejected")
	case errors.Is(err, valorantapi.ErrNotFound):
		return errors.New("NotFound: player not found")
	case errors.Is(err, valorantapi.ErrDecodeFailed):
		return errors.New("DecodeFailed: provider payload could not be decoded")
	case errors.Is(err, valorantapi.ErrProviderUnavailable):
		return errors.New("ProviderUnavailable: provider unavailable")
	default:
		return fmt.Errorf("UnknownProviderError: %w", err)
	}
}
