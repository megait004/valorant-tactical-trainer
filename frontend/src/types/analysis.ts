// Types liên quan đến Analysis Report (sinh từ Henrik/Riot match history).
// Mirror các struct trong internal/domain/analysis (Go).

export type Evidence = {
  matchIds?: string[]
  map?: string
  agent?: string
  metric: string
  value: number
  sampleSize: number
  comparisonBaseline: number
}

export type Finding = {
  id: string
  title: string
  severity: string
  confidence: string
  detail: string
  evidence: Evidence[]
}

export type Recommendation = {
  id: string
  findingId: string
  title: string
  reason: string
  drill: string
  cadence: string
}

export type BreakdownRow = {
  name: string
  matches: number
  rounds: number
  kd: number
  winRate: number
  headshotPercent: number
}

export type PracticeTask = {
  day: number
  focus: string
  map?: string
  agent?: string
  duration: string
  checklist: string[]
  evidence: string
}

export type MetricSummary = {
  matches: number
  rounds: number
  kd: number
  kda: number
  headshotPercent: number
  firstBloodRate: number
  firstDeathRate: number
  winRate: number
  weakestMap: string
  weakestMapWinRate: number
  weakestMapSample: number
  primaryRoleObserved: string
}

export type PlayerSnapshot = {
  name: string
  tagline: string
  region: string
  primaryRole: string
}

export type AnalysisReport = {
  player: PlayerSnapshot
  metrics: MetricSummary
  mapBreakdown: BreakdownRow[]
  agentBreakdown: BreakdownRow[]
  practicePlan: PracticeTask[]
  findings: Finding[]
  recommendations: Recommendation[]
}

export type LiveAnalysisResult = {
  report: AnalysisReport
  source: string
  cached: boolean
  fetchedAt: string
  message: string
}

export type LastReportResult = {
  hasReport: boolean
  result: LiveAnalysisResult
}
