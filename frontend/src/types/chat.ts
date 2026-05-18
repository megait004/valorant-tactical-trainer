// Bot AI chat (LLM coach) state.

export type ChatRole = 'user' | 'assistant'

export type ChatMessage = {
  role: ChatRole
  content: string
  createdAt: string
}

export type ChatState = {
  available: boolean
  message?: string
  history: ChatMessage[]
}
