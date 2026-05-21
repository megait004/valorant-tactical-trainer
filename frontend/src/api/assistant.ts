// Live Assistant — in-game tip engine. Khi không có bridge thì trả về session
// demo để UI vẫn render preview được.

import type { AssistantSessionState } from '../types'
import { bridge } from './bridge'
import { emptyAssistantTip, fallbackAssistantSession } from './fallbacks'

export const getAssistantSession = async () =>
  bridge()?.AssistantService?.GetSessionState?.() ?? fallbackAssistantSession

export const startAssistantSession = async (): Promise<AssistantSessionState> =>
  bridge()?.AssistantService?.StartSession?.() ?? {
    ...fallbackAssistantSession,
    active: true,
    message: 'Demo: chạy wails dev để rule engine thật.',
    currentAlert: {
      id: 'demo',
      title: 'Demo Live Assistant',
      message: 'Utility trước peek — fetch report ở Coach để gợi ý cá nhân.',
      severity: 'medium',
      source: 'demo',
    },
    queueSize: 1,
    tipsShown: 1,
    startedAt: new Date().toISOString(),
  }

export const stopAssistantSession = async () =>
  bridge()?.AssistantService?.StopSession?.() ?? {
    ...fallbackAssistantSession,
    message: 'Đã tắt (demo).',
  }

export const requestAssistantTip = async () =>
  bridge()?.AssistantService?.RequestTip?.() ?? emptyAssistantTip

export const markAssistantRoundStart = async () =>
  bridge()?.AssistantService?.MarkRoundStart?.() ?? emptyAssistantTip

export const pollAssistantAutoTip = async () =>
  bridge()?.AssistantService?.PollAutoTip?.() ?? emptyAssistantTip
