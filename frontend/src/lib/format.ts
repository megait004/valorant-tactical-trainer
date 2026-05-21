// Format helpers dùng chung — id slug cho Practice items, format thời lượng
// timer, datetime hiển thị, và rút message từ Error/Wails reject.

import type { PracticeTask } from '../types'

export const slugPart = (value: string | number) =>
  String(value).trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') || 'empty'

export const practiceItemID = (task: PracticeTask, item: string, index: number) =>
  [task.day, task.focus, task.map ?? 'all-map', task.agent ?? 'all-agent', index, item]
    .map(slugPart)
    .join('__')

export const taskID = (task: PracticeTask) =>
  [task.day, task.focus, task.map ?? 'all-map', task.agent ?? 'all-agent'].map(slugPart).join('__')

export const formatDuration = (seconds: number) => {
  const safeSeconds = Math.max(0, Math.floor(seconds))
  const minutes = Math.floor(safeSeconds / 60)
  const rest = safeSeconds % 60
  return `${String(minutes).padStart(2, '0')}:${String(rest).padStart(2, '0')}`
}

export const formatDateTime = (value: string) => {
  if (!value) return 'N/A'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

export const errorMessage = (err: unknown) => {
  if (err instanceof Error && err.message) return err.message
  if (typeof err === 'string' && err.trim()) return err
  if (err && typeof err === 'object') {
    const maybeMessage = 'message' in err ? (err as { message?: unknown }).message : undefined
    if (typeof maybeMessage === 'string' && maybeMessage.trim()) return maybeMessage
  }
  return 'err fetch report'
}
