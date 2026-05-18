// Map planner state — markers, lines, notes per map.

export type MapCatalogEntry = {
  id: string
  uuid: string
  name: string
  displayName: string
  imageUrl: string
  tacticalImageUrl: string
  hasTacticalLayout: boolean
}

export type MapCallout = {
  id: string
  label: string
  x: number
  y: number
  kind: 'site' | 'spawn' | 'lane' | 'mid'
}

export type PlanMarker = {
  id: string
  kind: string
  label: string
  x: number
  y: number
}

export type PlanLine = {
  id: string
  label: string
  x1: number
  y1: number
  x2: number
  y2: number
}

export type MapPlan = {
  mapId: string
  title: string
  side: string
  notes: string
  markers: PlanMarker[]
  lines: PlanLine[]
  updatedAt: string
}

export type PlannerTool = 'marker' | 'line'
export type MarkerKind = 'duelist' | 'initiator' | 'controller' | 'sentinel' | 'callout'
