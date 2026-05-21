// Data & Consent settings + Henrik API status. Fallback giá trị từ fallbacks.ts
// khi Wails bridge chưa sẵn sàng (dev mode browser thuần).

import type { APIStatus, DataSettings } from '../types'
import { bridge } from './bridge'
import { fallbackAPIStatus, fallbackDataSettings } from './fallbacks'

export const getDataSettings = async (): Promise<DataSettings> =>
  bridge()?.SettingsService?.GetDataSettings?.() ?? fallbackDataSettings

export const saveDataSettings = async (settings: DataSettings): Promise<DataSettings> =>
  bridge()?.SettingsService?.SaveDataSettings?.(settings) ??
  { ...settings, lastUpdatedAt: new Date().toISOString() }

export const getAPIStatus = async (): Promise<APIStatus> =>
  bridge()?.SettingsService?.GetAPIStatus?.() ?? fallbackAPIStatus
