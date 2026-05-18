package settings

type DataSettings struct {
	ConsentPersonalData bool   `json:"consentPersonalData"`
	RiotName            string `json:"riotName"`
	RiotTag             string `json:"riotTag"`
	PUUID               string `json:"puuid"`
	Region              string `json:"region"`
	Shard               string `json:"shard"`
	APIKey              string `json:"apiKey"`
	APIKeyHeader        string `json:"apiKeyHeader"`
	RateLimitTier       string `json:"rateLimitTier"`
	MatchCount          int    `json:"matchCount"`
	CacheTTLMinutes     int    `json:"cacheTTLMinutes"`
	LastUpdatedAt       string `json:"lastUpdatedAt"`
}

func DefaultDataSettings() DataSettings {
	return DataSettings{
		ConsentPersonalData: false,
		Region:              "ap",
		APIKeyHeader:        "Authorization",
		RateLimitTier:       "basic",
		MatchCount:          10,
		CacheTTLMinutes:     30,
	}
}
