package rank

import "time"

type Snapshot struct {
	PlayerPUUID     string
	Region          string
	Tier            int
	TierName        string
	RankingInTier   int
	MMRChangeToLast int
	Elo             int
	SeasonID        string
	FetchedAt       time.Time
	RawJSON         string
}
