// Live Assistant in-game tip engine.

export type AssistantAlert = {
  id: string
  title: string
  message: string
  severity: string
  source: string
}

export type AssistantSessionState = {
  active: boolean
  startedAt: string
  roundCount: number
  tipsShown: number
  lastAlertAt: string
  currentAlert?: AssistantAlert
  message: string
  queueSize: number
}

export type AssistantTipResult = {
  hasTip: boolean
  alert: AssistantAlert
  state: AssistantSessionState
}
