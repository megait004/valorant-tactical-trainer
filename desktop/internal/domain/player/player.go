package player

import (
	"strings"
	"time"
)

const ConsentVersion = "2026-05-13"

type Account struct {
	PUUID        string `json:"puuid"`
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	Region       string `json:"region"`
	AccountLevel int    `json:"accountLevel"`
	CardSmall    string `json:"cardSmall"`
	CardLarge    string `json:"cardLarge"`
	LastUpdate   string `json:"lastUpdate"`
}

type Consent struct {
	PlayerPUUID    string    `json:"playerPuuid"`
	Name           string    `json:"name"`
	Tag            string    `json:"tag"`
	Region         string    `json:"region"`
	Provider       string    `json:"provider"`
	ConsentVersion string    `json:"consentVersion"`
	ConsentedAt    time.Time `json:"consentedAt"`
}

func NormalizeName(value string) string {
	return strings.TrimSpace(value)
}

func NormalizeTag(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "#")
}

func NormalizeRegion(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "ap"
	}

	return normalized
}

func IsValidLookup(name string, tag string) bool {
	return NormalizeName(name) != "" && NormalizeTag(tag) != ""
}
