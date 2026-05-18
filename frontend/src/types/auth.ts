// Auth types — Riot Account-V1 login flow.

export type RiotPlayerInfo = {
  puuid: string
  gameName: string
  tagLine: string
  region: string
  shard: string
}

export type RiotLoginResult = {
  success: boolean
  error?: string
  playerInfo?: RiotPlayerInfo
}
