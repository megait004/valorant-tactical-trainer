// Practice progress (checklist) + session log (timer log).

import type {
  PracticeProgressState,
  PracticeSessionInput,
  PracticeSessionState,
} from '../types'
import { bridge } from './bridge'

export const getPracticeProgress = async (): Promise<PracticeProgressState> =>
  bridge()?.PracticeService?.GetPracticeProgress?.() ?? { items: {}, updatedAt: '' }

export const setPracticeProgress = async (
  itemID: string,
  done: boolean,
): Promise<PracticeProgressState> =>
  bridge()?.PracticeService?.SetPracticeProgress?.(itemID, done) ??
  { items: { [itemID]: done }, updatedAt: new Date().toISOString() }

export const resetPracticeProgress = async (): Promise<PracticeProgressState> =>
  bridge()?.PracticeService?.ResetPracticeProgress?.() ??
  { items: {}, updatedAt: new Date().toISOString() }

export const getPracticeSessions = async (): Promise<PracticeSessionState> =>
  bridge()?.PracticeService?.GetPracticeSessions?.() ?? { sessions: [], updatedAt: '' }

export const finishPracticeSession = async (
  input: PracticeSessionInput,
): Promise<PracticeSessionState> =>
  bridge()?.PracticeService?.FinishPracticeSession?.(input) ?? {
    sessions: [{ id: `session-${Date.now()}`, finishedAt: new Date().toISOString(), ...input }],
    updatedAt: new Date().toISOString(),
  }
