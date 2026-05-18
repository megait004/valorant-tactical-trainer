// Data & Consent settings + Henrik API status.

export type DataSettings = {
  consentPersonalData: boolean
  riotName: string
  riotTag: string
  region: string
  apiKey: string
  apiKeyHeader: string
  rateLimitTier: string
  matchCount: number
  cacheTTLMinutes: number
  lastUpdatedAt: string
}

export type APIStatus = {
  baseURL: string
  consentGranted: boolean
  canFetchPersonalData: boolean
  rateLimitPerMinute: number
  cacheTTLMinutes: number
  safeMode: boolean
  message: string
  nextStep: string
}
