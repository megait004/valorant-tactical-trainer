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
  agentPortrait?: string
  agentRole?: string
  agentUuid?: string
  abilitySlot?: string
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

export type PlannerTool = 'marker' | 'line' | 'agent'
export type MarkerKind = 'duelist' | 'initiator' | 'controller' | 'sentinel' | 'callout' | 'agent'

export type ValorantAbility = {
  slot: string
  displayName: string
  description: string
  displayIcon: string | null
}

export type ValorantAgentRole = {
  uuid: string
  displayName: string
  description: string
  displayIcon: string
}

export type ValorantAgent = {
  uuid: string
  displayName: string
  description: string
  displayIcon: string
  displayIconSmall: string
  bustPortrait: string | null
  fullPortrait: string | null
  minimapPortrait: string | null
  background: string | null
  backgroundGradientColors: string[]
  role: ValorantAgentRole
  abilities: ValorantAbility[]
}
