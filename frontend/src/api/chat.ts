// Bot AI chat (LLM coach hỏi-đáp). chatSendMessage cần bridge — sẽ throw nếu
// chưa sẵn sàng.

import type { ChatState } from '../types'
import { bridge } from './bridge'
import { fallbackChatState } from './fallbacks'

export const chatIsAvailable = async (): Promise<boolean> =>
  bridge()?.ChatService?.IsAvailable?.() ?? false

export const chatGetState = async (): Promise<ChatState> =>
  bridge()?.ChatService?.GetState?.() ?? fallbackChatState

export const chatSendMessage = async (message: string): Promise<ChatState> => {
  const send = bridge()?.ChatService?.SendMessage
  if (!send) {
    throw new Error('Bot AI chưa sẵn sàng — chạy bằng wails dev/build')
  }
  return send(message)
}

export const chatReset = async (): Promise<ChatState> =>
  bridge()?.ChatService?.Reset?.() ?? fallbackChatState
