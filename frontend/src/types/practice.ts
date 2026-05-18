// Practice progress (checklist) + session log (timer log).

export type PracticeProgressState = {
  items: Record<string, boolean>
  updatedAt: string
}

export type PracticeSession = {
  id: string
  taskId: string
  focus: string
  map: string
  agent: string
  durationSeconds: number
  startedAt: string
  finishedAt: string
}

export type PracticeSessionInput = {
  taskId: string
  focus: string
  map: string
  agent: string
  durationSeconds: number
  startedAt: string
}

export type PracticeSessionState = {
  sessions: PracticeSession[]
  updatedAt: string
}
