package match

type Summary struct {
	MatchID      string `json:"matchId"`
	PlayerPUUID  string `json:"playerPuuid"`
	MapName      string `json:"mapName"`
	Mode         string `json:"mode"`
	Queue        string `json:"queue"`
	SeasonID     string `json:"seasonId"`
	Region       string `json:"region"`
	Cluster      string `json:"cluster"`
	GameStart    int64  `json:"gameStart"`
	GameLength   int    `json:"gameLength"`
	RoundsPlayed int    `json:"roundsPlayed"`
	Agent        string `json:"agent"`
	Team         string `json:"team"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	Headshots    int    `json:"headshots"`
	Bodyshots    int    `json:"bodyshots"`
	Legshots     int    `json:"legshots"`
	DamageMade   int    `json:"damageMade"`
	RawJSON      string `json:"-"`
}
