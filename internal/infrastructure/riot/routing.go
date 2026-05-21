package riot

import "strings"

// File này gom logic map region → routing cluster / shard / VAL platform host.
// Riot chia API ra:
//
//   - Account-V1: gọi qua routing cluster (americas/europe/asia).
//   - VAL-MATCH-V1: gọi qua VAL platform host (na/eu/kr/ap.api.riotgames.com).
//
// Cùng 1 region (vd "ap") map ra cùng shard ("ap") nhưng routing cluster
// có thể khác (vd "kr" routing = asia, shard = kr).

// resolveRoutingShard cho Account-V1 (xem account.go).
func resolveRoutingShard(region string) (routing, shard string) {
	switch region {
	case "na", "br", "latam":
		return "americas", "na"
	case "eu", "euw", "eune", "tr", "ru":
		return "europe", "eu"
	case "kr":
		return "asia", "kr"
	case "ap", "jp", "oce":
		return "asia", "ap"
	default:
		return "asia", "ap"
	}
}

// defaultMatchHost cho VAL-MATCH-V1 (xem match.go). Map shard → host.
// Riot doc: developer.riotgames.com/apis#val-match-v1
func defaultMatchHost(shard string) string {
	switch strings.ToLower(strings.TrimSpace(shard)) {
	case "na", "br", "latam":
		return "na.api.riotgames.com"
	case "eu":
		return "eu.api.riotgames.com"
	case "kr":
		return "kr.api.riotgames.com"
	case "ap":
		return "ap.api.riotgames.com"
	default:
		return "ap.api.riotgames.com"
	}
}
