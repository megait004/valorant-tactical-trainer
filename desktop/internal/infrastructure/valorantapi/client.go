package valorantapi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	matchdomain "valorant-tactical-trainer/internal/domain/match"
	"valorant-tactical-trainer/internal/domain/player"
	"valorant-tactical-trainer/internal/domain/rank"
)

const providerName = "Henrik unofficial VALORANT API"

var (
	ErrNotFound            = errors.New("not found")
	ErrRateLimited         = errors.New("rate limited")
	ErrUnauthorizedAPIKey  = errors.New("unauthorized api key")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrDecodeFailed        = errors.New("decode failed")
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	limiter    *RateLimiter
}

type RateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

type Option func(*Client)

func WithAPIKey(apiKey string) Option {
	return func(client *Client) {
		client.apiKey = strings.TrimSpace(apiKey)
	}
}

func WithBaseURL(baseURL string) Option {
	return func(client *Client) {
		client.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		if httpClient != nil {
			client.httpClient = httpClient
		}
	}
}

func WithRateLimiter(limiter *RateLimiter) Option {
	return func(client *Client) {
		client.limiter = limiter
	}
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	if interval <= 0 {
		return nil
	}

	return &RateLimiter{interval: interval}
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		baseURL: "https://api.henrikdev.xyz/valorant/v1",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

var basicRateLimiter = NewRateLimiter(2 * time.Second)

func NewBasicClient(opts ...Option) *Client {
	options := append([]Option{WithRateLimiter(basicRateLimiter)}, opts...)
	return NewClient(options...)
}

func (client *Client) LookupAccount(ctx context.Context, name string, tag string) (player.Account, error) {
	escapedName := url.PathEscape(player.NormalizeName(name))
	escapedTag := url.PathEscape(player.NormalizeTag(tag))
	endpoint := fmt.Sprintf("%s/account/%s/%s", client.baseURL, escapedName, escapedTag)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return player.Account{}, fmt.Errorf("create account request: %w", err)
	}

	if client.apiKey != "" {
		req.Header.Set("Authorization", client.apiKey)
	}

	resp, err := client.do(req)
	if err != nil {
		return player.Account{}, fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return player.Account{}, ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return player.Account{}, ErrUnauthorizedAPIKey
	case http.StatusNotFound:
		return player.Account{}, ErrNotFound
	default:
		if resp.StatusCode >= 500 {
			return player.Account{}, ErrProviderUnavailable
		}
		return player.Account{}, fmt.Errorf("provider status %d", resp.StatusCode)
	}

	var payload accountResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return player.Account{}, fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	if payload.Status == http.StatusNotFound || len(payload.Errors) > 0 && payload.Data.PUUID == "" {
		return player.Account{}, ErrNotFound
	}

	return player.Account{
		PUUID:        payload.Data.PUUID,
		Name:         payload.Data.Name,
		Tag:          payload.Data.Tag,
		Region:       payload.Data.Region,
		AccountLevel: payload.Data.AccountLevel,
		CardSmall:    payload.Data.Card.Small,
		CardLarge:    payload.Data.Card.Large,
		LastUpdate:   payload.Data.LastUpdate,
	}, nil
}

func (client *Client) MatchesByPUUID(ctx context.Context, puuid string, region string, size string) ([]matchdomain.Summary, string, error) {
	normalizedRegion := player.NormalizeRegion(region)
	normalizedSize := strings.TrimSpace(size)
	if normalizedSize == "" {
		normalizedSize = "10"
	}

	endpoint := fmt.Sprintf("%s/by-puuid/matches/%s/%s?size=%s", strings.Replace(client.baseURL, "/v1", "/v3", 1), url.PathEscape(normalizedRegion), url.PathEscape(strings.TrimSpace(puuid)), url.QueryEscape(normalizedSize))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create matches request: %w", err)
	}

	if client.apiKey != "" {
		req.Header.Set("Authorization", client.apiKey)
	}

	resp, err := client.do(req)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return nil, "", ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, "", ErrUnauthorizedAPIKey
	case http.StatusNotFound:
		return nil, "", ErrNotFound
	default:
		if resp.StatusCode >= 500 {
			return nil, "", ErrProviderUnavailable
		}
		return nil, "", fmt.Errorf("provider status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	var payload matchesResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	summaries := make([]matchdomain.Summary, 0, len(payload.Data))
	for _, item := range payload.Data {
		summary := matchdomain.Summary{
			MatchID:      item.Metadata.MatchID,
			PlayerPUUID:  puuid,
			MapName:      item.Metadata.Map,
			Mode:         item.Metadata.Mode,
			Queue:        item.Metadata.Queue,
			SeasonID:     item.Metadata.SeasonID,
			Region:       item.Metadata.Region,
			Cluster:      item.Metadata.Cluster,
			GameStart:    item.Metadata.GameStart,
			GameLength:   item.Metadata.GameLength,
			RoundsPlayed: item.Metadata.RoundsPlayed,
			RawJSON:      string(body),
		}

		for _, participant := range item.Players.AllPlayers {
			if participant.PUUID == puuid {
				summary.Agent = participant.Character
				summary.Team = participant.Team
				summary.Kills = participant.Stats.Kills
				summary.Deaths = participant.Stats.Deaths
				summary.Assists = participant.Stats.Assists
				summary.Headshots = participant.Stats.Headshots
				summary.Bodyshots = participant.Stats.Bodyshots
				summary.Legshots = participant.Stats.Legshots
				summary.DamageMade = participant.DamageMade
				break
			}
		}

		if summary.MatchID == "" {
			summary.MatchID = fallbackMatchID(body, len(summaries))
		}

		summaries = append(summaries, summary)
	}

	return summaries, string(body), nil
}

func (client *Client) MMRByPUUID(ctx context.Context, puuid string, region string) (rank.Snapshot, string, error) {
	normalizedRegion := player.NormalizeRegion(region)
	endpoint := fmt.Sprintf("%s/by-puuid/mmr/%s/%s", strings.Replace(client.baseURL, "/v1", "/v2", 1), url.PathEscape(normalizedRegion), url.PathEscape(strings.TrimSpace(puuid)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return rank.Snapshot{}, "", fmt.Errorf("create mmr request: %w", err)
	}

	if client.apiKey != "" {
		req.Header.Set("Authorization", client.apiKey)
	}

	resp, err := client.do(req)
	if err != nil {
		return rank.Snapshot{}, "", fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return rank.Snapshot{}, "", ErrRateLimited
	case http.StatusUnauthorized, http.StatusForbidden:
		return rank.Snapshot{}, "", ErrUnauthorizedAPIKey
	case http.StatusNotFound:
		return rank.Snapshot{}, "", ErrNotFound
	default:
		if resp.StatusCode >= 500 {
			return rank.Snapshot{}, "", ErrProviderUnavailable
		}
		return rank.Snapshot{}, "", fmt.Errorf("provider status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rank.Snapshot{}, "", fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	var payload mmrResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return rank.Snapshot{}, "", fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}

	snapshot := rank.Snapshot{
		PlayerPUUID:     strings.TrimSpace(puuid),
		Region:          normalizedRegion,
		Tier:            payload.Data.CurrentTier,
		TierName:        payload.Data.CurrentTierPatched,
		RankingInTier:   payload.Data.RankingInTier,
		MMRChangeToLast: payload.Data.MMRChangeToLastGame,
		Elo:             payload.Data.Elo,
		SeasonID:        payload.Data.SeasonID,
		FetchedAt:       time.Now().UTC(),
		RawJSON:         string(body),
	}

	if snapshot.TierName == "" && payload.Data.Images.Small != "" {
		snapshot.TierName = payload.Data.Images.Small
	}

	return snapshot, string(body), nil
}

func ProviderName() string {
	return providerName
}

func (client *Client) do(req *http.Request) (*http.Response, error) {
	if client.limiter != nil {
		if err := client.limiter.Wait(req.Context()); err != nil {
			return nil, err
		}
	}

	return client.httpClient.Do(req)
}

func (limiter *RateLimiter) Wait(ctx context.Context) error {
	limiter.mu.Lock()
	now := time.Now()
	wait := limiter.next.Sub(now)
	if wait < 0 {
		wait = 0
	}
	limiter.next = now.Add(wait).Add(limiter.interval)
	limiter.mu.Unlock()

	if wait == 0 {
		return nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}

	return nil
}

type accountResponse struct {
	Status int `json:"status"`
	Data   struct {
		PUUID        string `json:"puuid"`
		Region       string `json:"region"`
		AccountLevel int    `json:"account_level"`
		Name         string `json:"name"`
		Tag          string `json:"tag"`
		Card         struct {
			Small string `json:"small"`
			Large string `json:"large"`
		} `json:"card"`
		LastUpdate string `json:"last_update"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Details string `json:"details"`
	} `json:"errors"`
}

type matchesResponse struct {
	Status int `json:"status"`
	Data   []struct {
		Metadata struct {
			Map          string `json:"map"`
			GameLength   int    `json:"game_length"`
			GameStart    int64  `json:"game_start"`
			RoundsPlayed int    `json:"rounds_played"`
			Mode         string `json:"mode"`
			Queue        string `json:"queue"`
			SeasonID     string `json:"season_id"`
			MatchID      string `json:"matchid"`
			Region       string `json:"region"`
			Cluster      string `json:"cluster"`
		} `json:"metadata"`
		Players struct {
			AllPlayers []struct {
				PUUID     string `json:"puuid"`
				Team      string `json:"team"`
				Character string `json:"character"`
				Stats     struct {
					Kills     int `json:"kills"`
					Deaths    int `json:"deaths"`
					Assists   int `json:"assists"`
					Headshots int `json:"headshots"`
					Bodyshots int `json:"bodyshots"`
					Legshots  int `json:"legshots"`
				} `json:"stats"`
				DamageMade int `json:"damage_made"`
			} `json:"all_players"`
		} `json:"players"`
	} `json:"data"`
}

type mmrResponse struct {
	Status int `json:"status"`
	Data   struct {
		CurrentTier         int    `json:"currenttier"`
		CurrentTierPatched  string `json:"currenttierpatched"`
		RankingInTier       int    `json:"ranking_in_tier"`
		MMRChangeToLastGame int    `json:"mmr_change_to_last_game"`
		Elo                 int    `json:"elo"`
		SeasonID            string `json:"season_id"`
		Images              struct {
			Small string `json:"small"`
		} `json:"images"`
	} `json:"data"`
}

func fallbackMatchID(body []byte, index int) string {
	hash := sha256.Sum256(append(body, byte(index)))
	return fmt.Sprintf("generated-%x", hash[:8])
}
