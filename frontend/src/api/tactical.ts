// Tactical map planner — catalog + per-map plan (markers, lines, notes).
// Fallback dùng catalog tĩnh để UI hiển thị maps ngay cả khi bridge chưa
// load xong.

import type { MapCatalogEntry, MapPlan } from '../types'
import { bridge } from './bridge'
import { defaultMapPlan } from './fallbacks'

import { fallbackMapCatalog } from '../features/maps/catalog-fallback'

export const listTacticalMaps = async (): Promise<MapCatalogEntry[]> =>
  bridge()?.TacticalService?.ListMaps?.() ?? fallbackMapCatalog

export const loadMapPlan = async (mapId: string): Promise<MapPlan> =>
  bridge()?.TacticalService?.LoadMapPlan?.(mapId) ?? defaultMapPlan(mapId)

export const saveMapPlan = async (plan: MapPlan): Promise<MapPlan> =>
  bridge()?.TacticalService?.SaveMapPlan?.(plan) ??
  { ...plan, updatedAt: new Date().toISOString() }

export const deleteMapPlan = async (mapId: string): Promise<void> => {
  const del = bridge()?.TacticalService?.DeleteMapPlan
  if (del) {
    await del(mapId)
  }
}
