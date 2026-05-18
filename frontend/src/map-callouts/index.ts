import type { MapCallout } from '../types'
import { buildMapCallouts, mapData } from './raw'

const cache: Record<string, MapCallout[]> = {}

const genericTwoSiteFallback: MapCallout[] = [
  { id: 'site-a', label: 'A', x: 24, y: 46, kind: 'site' },
  { id: 'site-b', label: 'B', x: 78, y: 46, kind: 'site' },
  { id: 'att-spawn', label: 'SPAWN TẤN CÔNG', x: 50, y: 16, kind: 'spawn' },
  { id: 'def-spawn', label: 'SPAWN PHÒNG THỦ', x: 50, y: 86, kind: 'spawn' },
]

export const getMapCallouts = (mapId: string): MapCallout[] => {
  if (cache[mapId]) return cache[mapId]
  if (!mapData[mapId]) return genericTwoSiteFallback
  const callouts = buildMapCallouts(mapId)
  cache[mapId] = callouts
  return callouts
}
