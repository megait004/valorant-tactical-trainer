import { useCallback, useEffect, useRef, useState } from 'react'
import type { FC, MouseEvent, ReactNode } from 'react'
import { deleteMapPlan, listTacticalMaps, loadMapPlan, saveMapPlan } from '../../api'
import { fallbackMapCatalog, mapIdFromName } from './catalog-fallback'
import TacticalBoard from './TacticalBoard'
import type {
  MapCatalogEntry,
  MapCallout,
  MapPlan,
  MarkerKind,
  PlanLine,
  PlanMarker,
  PlannerTool,
  ValorantAbility,
  ValorantAgent,
} from '../../types'
import { getMapCallouts } from './callouts'
import { ROLE_COLORS, ROLE_ORDER, useAgents } from '../agents/useAgents'

type MapPlannerPanelProps = {
  suggestedMapName?: string
}

const SLOT_KEY: Record<string, string> = {
  Ability1: 'E',
  Ability2: 'Q',
  Grenade: 'C',
  Ultimate: 'X',
}

const markerKinds: Array<{ id: MarkerKind; label: string; color: string }> = [
  { id: 'duelist', label: 'Đối đầu', color: '#ff4655' },
  { id: 'initiator', label: 'Khởi tranh', color: '#fbbf24' },
  { id: 'controller', label: 'Kiểm soát', color: '#42e8f3' },
  { id: 'sentinel', label: 'Hộ vệ', color: '#4ade80' },
  { id: 'callout', label: 'Callout', color: '#e2e8f0' },
]

const markerLabel = (kind: MarkerKind) => markerKinds.find((item) => item.id === kind)?.label ?? kind

const kindColor = (kind: string, agentRole?: string) => {
  if (kind === 'agent' && agentRole) return ROLE_COLORS[agentRole] ?? '#42e8f3'
  if (kind === 'ability' && agentRole) return ROLE_COLORS[agentRole] ?? '#42e8f3'
  return markerKinds.find((item) => item.id === kind)?.color ?? '#42e8f3'
}

const catalogHasMap = (name: string) => {
  const id = mapIdFromName(name)
  return fallbackMapCatalog.some((map) => map.id === id || map.name.toLowerCase() === name.toLowerCase())
}

// ─── Ability Context Menu (fixed position popup) ───────────────────────────────

type AbilityMenu = {
  marker: PlanMarker
  agent: ValorantAgent
  screenX: number
  screenY: number
}

const AbilityContextMenu: FC<{
  menu: AbilityMenu
  onSelect: (ability: ValorantAbility) => void
  onClose: () => void
}> = ({ menu, onSelect, onClose }) => {
  const { agent, screenX, screenY } = menu
  const roleColor = ROLE_COLORS[agent.role?.displayName ?? ''] ?? '#42e8f3'

  const abilityOrder = ['Ability2', 'Grenade', 'Ability1', 'Ultimate']
  const sorted = [...agent.abilities].sort(
    (a, b) => abilityOrder.indexOf(a.slot) - abilityOrder.indexOf(b.slot),
  )

  // Adjust so popup doesn't go off screen
  const left = Math.min(screenX, window.innerWidth - 260)
  const top = Math.min(screenY, window.innerHeight - 280)

  return (
    <>
      {/* Backdrop */}
      <button
        type="button"
        aria-label="Đóng menu"
        className="fixed inset-0 z-[100] cursor-default bg-transparent"
        onClick={onClose}
        onContextMenu={(e) => { e.preventDefault(); onClose() }}
      />

      {/* Popup */}
      <div
        className="fixed z-[101] w-56 overflow-hidden rounded-2xl border border-tactical-line bg-[#0f1117] shadow-2xl shadow-black/60"
        style={{ left, top }}
      >
        {/* Agent header */}
        <div
          className="flex items-center gap-2 px-3 py-2"
          style={{ borderBottom: `1px solid ${roleColor}30`, background: `${roleColor}12` }}
        >
          {agent.displayIcon && (
            <img src={agent.displayIcon} alt="" className="h-8 w-8 rounded-lg object-cover bg-black/40" />
          )}
          <div>
            <p className="text-xs font-black text-white">{agent.displayName}</p>
            <p className="text-[10px] font-bold" style={{ color: roleColor }}>
              {agent.role?.displayName} · Chọn chiêu để đặt
            </p>
          </div>
        </div>

        {/* Ability list */}
        <div className="p-2 space-y-1">
          {sorted.map((ability) => {
            const key = SLOT_KEY[ability.slot] ?? ability.slot
            const isUlt = ability.slot === 'Ultimate'
            return (
              <button
                key={ability.slot}
                type="button"
                onClick={() => onSelect(ability)}
                className="group flex w-full items-center gap-2.5 rounded-xl border border-transparent px-2 py-2 text-left transition hover:border-white/10 hover:bg-white/5"
              >
                <div
                  className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg"
                  style={{ backgroundColor: isUlt ? `${roleColor}25` : 'rgba(255,255,255,0.05)' }}
                >
                  {ability.displayIcon ? (
                    <img
                      src={ability.displayIcon}
                      alt={ability.displayName}
                      className="h-6 w-6 object-contain"
                    />
                  ) : (
                    <span className="text-sm font-black text-slate-400">{key}</span>
                  )}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-1.5">
                    <span
                      className="rounded px-1 py-0.5 text-[9px] font-black"
                      style={isUlt
                        ? { backgroundColor: `${roleColor}30`, color: roleColor }
                        : { backgroundColor: 'rgba(255,255,255,0.08)', color: '#94a3b8' }
                      }
                    >
                      {key}
                    </span>
                    <span className="truncate text-xs font-bold text-white">{ability.displayName}</span>
                  </div>
                  <p className="mt-0.5 line-clamp-2 text-[9px] leading-relaxed text-slate-500">
                    {ability.description.slice(0, 80)}…
                  </p>
                </div>
              </button>
            )
          })}
        </div>

        <div className="border-t border-white/5 px-3 py-2 text-[10px] text-slate-600">
          Click chiêu → click bản đồ để đánh dấu
        </div>
      </div>
    </>
  )
}

// ─── Agent Picker Panel ────────────────────────────────────────────────────────

const AgentPickerPanel: FC<{
  agents: ValorantAgent[]
  byRole: Record<string, ValorantAgent[]>
  selectedAgent: ValorantAgent | null
  onSelect: (agent: ValorantAgent | null) => void
}> = ({ agents, byRole, selectedAgent, onSelect }) => {
  const [activeRole, setActiveRole] = useState<string | null>(null)
  const roles = ROLE_ORDER.filter((r) => byRole[r]?.length)
  const displayed = activeRole ? (byRole[activeRole] ?? []) : agents

  const roleIcons = Object.fromEntries(
    agents.filter((a) => a.role?.displayName).map((a) => [a.role.displayName, a.role.displayIcon]),
  ) as Record<string, string>

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-1">
        <button
          type="button"
          onClick={() => setActiveRole(null)}
          className={`rounded-lg border px-2 py-1 text-[10px] font-black uppercase tracking-wider transition ${
            activeRole ? 'border-transparent text-slate-600 hover:text-slate-400' : 'border-white/20 bg-white/10 text-white'
          }`}
        >
          Tất cả
        </button>
        {roles.map((role) => {
          const color = ROLE_COLORS[role] ?? '#42e8f3'
          const active = activeRole === role
          return (
            <button
              key={role}
              type="button"
              onClick={() => setActiveRole(active ? null : role)}
              title={role}
              className="rounded-lg border px-2 py-1 text-[10px] font-black uppercase tracking-wider transition"
              style={
                active
                  ? { backgroundColor: `${color}22`, borderColor: `${color}55`, color }
                  : { borderColor: 'transparent', color: '#64748b' }
              }
            >
              {roleIcons[role] && (
                <img src={roleIcons[role]} alt={role} className="h-3 w-3 object-contain inline-block mr-1" />
              )}
              {role.slice(0, 3)}
            </button>
          )
        })}
      </div>

      {selectedAgent && (
        <div
          className="rounded-xl border p-3"
          style={{
            borderColor: `${ROLE_COLORS[selectedAgent.role?.displayName ?? ''] ?? '#42e8f3'}40`,
            backgroundColor: `${ROLE_COLORS[selectedAgent.role?.displayName ?? ''] ?? '#42e8f3'}10`,
          }}
        >
          <div className="flex items-center gap-2 mb-2">
            {selectedAgent.displayIcon && (
              <img src={selectedAgent.displayIcon} alt="" className="h-8 w-8 rounded-lg object-cover bg-black/40" />
            )}
            <div>
              <p className="text-xs font-black text-white">{selectedAgent.displayName}</p>
              <p className="text-[10px]" style={{ color: ROLE_COLORS[selectedAgent.role?.displayName ?? ''] ?? '#42e8f3' }}>
                {selectedAgent.role?.displayName}
              </p>
            </div>
            <button type="button" onClick={() => onSelect(null)} className="ml-auto text-xs text-slate-600 hover:text-slate-400">
              ✕
            </button>
          </div>
          <div className="grid grid-cols-4 gap-1">
            {selectedAgent.abilities.map((ab) => (
              <div
                key={ab.slot}
                title={`${ab.displayName}: ${ab.description}`}
                className="flex flex-col items-center gap-0.5 rounded-lg bg-black/30 p-1.5"
              >
                {ab.displayIcon ? (
                  <img src={ab.displayIcon} alt={ab.displayName} className="h-5 w-5 object-contain" />
                ) : (
                  <span className="text-xs text-slate-500">?</span>
                )}
                <span className="text-[8px] text-slate-500 truncate w-full text-center">{ab.displayName.slice(0, 8)}</span>
              </div>
            ))}
          </div>
          <p className="mt-2 text-[10px] text-slate-500">
            Click bản đồ để đặt · Chuột phải agent trên map → chọn chiêu
          </p>
        </div>
      )}

      <div className="grid grid-cols-4 gap-1.5 max-h-[340px] overflow-y-auto pr-1">
        {displayed.map((agent) => {
          const role = agent.role?.displayName ?? ''
          const color = ROLE_COLORS[role] ?? '#42e8f3'
          const isSelected = selectedAgent?.uuid === agent.uuid
          return (
            <button
              key={agent.uuid}
              type="button"
              onClick={() => onSelect(isSelected ? null : agent)}
              title={agent.displayName}
              className="group relative flex flex-col items-center rounded-xl border p-1.5 transition"
              style={
                isSelected
                  ? { backgroundColor: `${color}22`, borderColor: `${color}60` }
                  : { borderColor: 'transparent', backgroundColor: 'rgba(0,0,0,0.2)' }
              }
            >
              <div className="relative h-10 w-full overflow-hidden rounded-lg bg-black/40">
                {agent.displayIcon ? (
                  <img
                    src={agent.displayIcon}
                    alt={agent.displayName}
                    className="h-full w-full object-cover transition group-hover:scale-110"
                  />
                ) : (
                  <div className="h-full w-full bg-slate-800" />
                )}
                <div className="absolute bottom-0 left-0 right-0 h-0.5 rounded-b" style={{ backgroundColor: color }} />
              </div>
              <span
                className="mt-1 w-full truncate text-center text-[9px] font-bold leading-tight"
                style={{ color: isSelected ? color : '#94a3b8' }}
              >
                {agent.displayName}
              </span>
            </button>
          )
        })}
      </div>
    </div>
  )
}

// ─── Main Panel ────────────────────────────────────────────────────────────────

const MapPlannerPanel: FC<MapPlannerPanelProps> = ({ suggestedMapName }) => {
  const boardRef = useRef<HTMLDivElement>(null)
  const [maps, setMaps] = useState<MapCatalogEntry[]>([])
  const [selectedMapId, setSelectedMapId] = useState('ascent')
  const [plan, setPlan] = useState<MapPlan | null>(null)
  const [tool, setTool] = useState<PlannerTool>('marker')
  const [markerKind, setMarkerKind] = useState<MarkerKind>('controller')
  const [lineStart, setLineStart] = useState<{ x: number; y: number } | null>(null)
  const [message, setMessage] = useState('')
  const [busy, setBusy] = useState(false)
  const [showCallouts, setShowCallouts] = useState(true)
  const [selectedAgent, setSelectedAgent] = useState<ValorantAgent | null>(null)

  // Ability placement state
  const [abilityMenu, setAbilityMenu] = useState<AbilityMenu | null>(null)
  const [pendingAbility, setPendingAbility] = useState<{
    ability: ValorantAbility
    sourceMarker: PlanMarker
  } | null>(null)

  const { agents, byRole } = useAgents()

  const selectedMap = maps.find((item) => item.id === selectedMapId)

  const loadMapsAndPlan = useCallback(async (mapId: string) => {
    setMaps(await listTacticalMaps())
    setPlan(await loadMapPlan(mapId))
  }, [])

  useEffect(() => {
    const initial =
      suggestedMapName && catalogHasMap(suggestedMapName) ? mapIdFromName(suggestedMapName) : 'ascent'
    setSelectedMapId(initial)
    void loadMapsAndPlan(initial)
  }, [suggestedMapName, loadMapsAndPlan])

  const selectMap = async (mapId: string) => {
    setSelectedMapId(mapId)
    setLineStart(null)
    setPendingAbility(null)
    setAbilityMenu(null)
    setPlan(await loadMapPlan(mapId))
  }

  const updatePlan = (patch: Partial<MapPlan>) => {
    setPlan((current) => (current ? { ...current, ...patch, mapId: selectedMapId } : current))
  }

  const pointerPercent = (event: MouseEvent<HTMLDivElement>) => {
    const rect = boardRef.current?.getBoundingClientRect()
    if (!rect?.width || !rect.height) return null
    return {
      x: Math.min(100, Math.max(0, Number((((event.clientX - rect.left) / rect.width) * 100).toFixed(2)))),
      y: Math.min(100, Math.max(0, Number((((event.clientY - rect.top) / rect.height) * 100).toFixed(2)))),
    }
  }

  // Right-click on agent marker → look up agent → show ability menu.
  // Auto-cancel any pending ability from a different agent so user doesn't
  // need to press Esc when switching between agents.
  const handleAgentRightClick = (marker: PlanMarker, clientX: number, clientY: number) => {
    const agent = agents.find((a) => a.uuid === marker.agentUuid || a.displayName === marker.label)
    if (!agent) return
    setAbilityMenu({ marker, agent, screenX: clientX, screenY: clientY })
    setPendingAbility(null)
    setMessage('')
  }

  // Picking an agent in the sidebar should also cancel any pending ability
  // from the previously selected agent.
  const handleAgentPickerSelect = (agent: ValorantAgent | null) => {
    setSelectedAgent(agent)
    if (pendingAbility) {
      setPendingAbility(null)
      setMessage(agent ? `Đã chọn ${agent.displayName}. Click bản đồ để đặt.` : '')
    }
    setAbilityMenu(null)
  }

  // User selected an ability from the context menu
  const handleAbilitySelect = (ability: ValorantAbility) => {
    if (!abilityMenu) return
    setPendingAbility({ ability, sourceMarker: abilityMenu.marker })
    setAbilityMenu(null)
    setMessage(`Đang đặt chiêu "${ability.displayName}" — click vào bản đồ để chọn vị trí.`)
  }

  const handleBoardClick = (event: MouseEvent<HTMLDivElement>) => {
    if (!plan) return
    let point = pointerPercent(event)
    if (!point) return

    // ── Ability placement mode ──
    if (pendingAbility) {
      const { ability, sourceMarker } = pendingAbility
      updatePlan({
        markers: [
          ...plan.markers,
          {
            id: `ability-${Date.now()}`,
            kind: 'ability',
            label: ability.displayName,
            x: point.x,
            y: point.y,
            agentPortrait: ability.displayIcon ?? undefined,
            agentRole: sourceMarker.agentRole,
            agentUuid: sourceMarker.agentUuid,
            abilitySlot: ability.slot,
          },
        ],
      })
      setMessage(`Đã đặt chiêu "${ability.displayName}" lên bản đồ. Click tiếp để đặt thêm hoặc nhấn Esc để hủy.`)
      return
    }

    // ── Agent placement mode ──
    const snapped = snapToCallout(point.x, point.y, getMapCallouts(selectedMapId))
    if (snapped) point = { x: snapped.x, y: snapped.y }

    if (tool === 'agent') {
      if (!selectedAgent) {
        setMessage('Chọn một agent trước khi đặt lên bản đồ.')
        return
      }
      const agentRole = selectedAgent.role?.displayName ?? ''
      updatePlan({
        markers: [
          ...plan.markers,
          {
            id: `agent-${selectedAgent.uuid}-${Date.now()}`,
            kind: 'agent',
            label: selectedAgent.displayName,
            x: point.x,
            y: point.y,
            agentPortrait: selectedAgent.minimapPortrait ?? selectedAgent.displayIcon,
            agentRole,
            agentUuid: selectedAgent.uuid,
          },
        ],
      })
      setMessage(`Đã đặt ${selectedAgent.displayName}. Chuột phải lên avatar để chọn chiêu.`)
      return
    }

    // ── Marker mode ──
    if (tool === 'marker') {
      updatePlan({
        markers: [
          ...plan.markers,
          {
            id: `marker-${Date.now()}`,
            kind: markerKind,
            label: snapped?.label ?? markerLabel(markerKind),
            x: point.x,
            y: point.y,
          },
        ],
      })
      return
    }

    // ── Line mode ──
    if (!lineStart) {
      setLineStart(point)
      setMessage('Chọn điểm kết thúc cho route.')
      return
    }

    const line: PlanLine = {
      id: `line-${Date.now()}`,
      label: 'route',
      x1: lineStart.x,
      y1: lineStart.y,
      x2: point.x,
      y2: point.y,
    }
    updatePlan({ lines: [...plan.lines, line] })
    setLineStart(null)
    setMessage('Đã thêm route.')
  }

  // Esc cancels pending ability
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        if (pendingAbility) {
          setPendingAbility(null)
          setMessage('Đã hủy đặt chiêu.')
        }
        setAbilityMenu(null)
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [pendingAbility])

  const handleSave = async () => {
    if (!plan) return
    setBusy(true)
    try {
      setPlan(await saveMapPlan({ ...plan, mapId: selectedMapId }))
      setMessage(`Đã lưu kế hoạch ${selectedMap?.displayName ?? selectedMapId}.`)
    } finally {
      setBusy(false)
    }
  }

  const handleReset = async () => {
    if (!window.confirm('Xóa kế hoạch đã lưu cho map này?')) return
    setBusy(true)
    try {
      await deleteMapPlan(selectedMapId)
      setPlan(await loadMapPlan(selectedMapId))
      setLineStart(null)
      setPendingAbility(null)
      setMessage('Đã xóa kế hoạch map.')
    } finally {
      setBusy(false)
    }
  }

  if (!plan) {
    return <div className="rounded-3xl border border-tactical-line bg-tactical-panel p-6 text-slate-300">Đang tải bản đồ...</div>
  }

  const agentCount = plan.markers.filter((m) => m.kind === 'agent').length
  const abilityCount = plan.markers.filter((m) => m.kind === 'ability').length
  const otherCount = plan.markers.filter((m) => m.kind !== 'agent' && m.kind !== 'ability').length

  return (
    <>
      {/* Ability context menu portal */}
      {abilityMenu && (
        <AbilityContextMenu
          menu={abilityMenu}
          onSelect={handleAbilitySelect}
          onClose={() => setAbilityMenu(null)}
        />
      )}

      <section className="grid gap-5 xl:grid-cols-[280px_1fr]">
        {/* Left sidebar */}
        <aside className="space-y-4">
          {/* Map pool */}
          <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-4 shadow-xl shadow-black/20">
            <p className="text-xs font-bold uppercase tracking-[0.2em] text-tactical-cyan">Map pool</p>
            <p className="mt-1 text-sm text-slate-400">12 map competitive · tactical mô phỏng · lưu local</p>
            {suggestedMapName && (
              <p className="mt-2 rounded-xl border border-tactical-red/30 bg-tactical-red/10 px-3 py-2 text-xs text-tactical-red">
                Coach gợi ý: {suggestedMapName}
              </p>
            )}
            <div className="mt-3 max-h-[380px] space-y-1 overflow-y-auto pr-1">
              {maps.map((map) => (
                <button
                  key={map.id}
                  type="button"
                  onClick={() => void selectMap(map.id)}
                  className={`flex w-full items-center gap-3 rounded-xl border px-2 py-2 text-left transition ${
                    selectedMapId === map.id
                      ? 'border-tactical-cyan/40 bg-tactical-cyan/10'
                      : 'border-transparent hover:border-white/10 hover:bg-white/[0.03]'
                  }`}
                >
                  <img src={map.imageUrl} alt="" className="h-10 w-16 rounded-lg object-cover bg-black/40" />
                  <span className="text-sm font-bold text-white">
                    {map.displayName}
                    {map.hasTacticalLayout && <span className="ml-1 text-[10px] text-tactical-cyan">TAC</span>}
                  </span>
                </button>
              ))}
            </div>
          </div>

          {/* Agent picker */}
          {tool === 'agent' && (
            <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-4 shadow-xl shadow-black/20">
              <p className="mb-1 text-xs font-bold uppercase tracking-[0.18em] text-tactical-cyan">Chọn Agent</p>
              <p className="mb-3 text-[11px] text-slate-500">
                Chọn agent → click bản đồ để đặt → chuột phải avatar trên map → chọn chiêu
              </p>
              {agents.length === 0 ? (
                <div className="flex items-center justify-center py-8">
                  <div className="h-6 w-6 animate-spin rounded-full border-2 border-tactical-cyan border-t-transparent" />
                </div>
              ) : (
                <AgentPickerPanel
                  agents={agents}
                  byRole={byRole}
                  selectedAgent={selectedAgent}
                  onSelect={handleAgentPickerSelect}
                />
              )}
            </div>
          )}
        </aside>

        <div className="space-y-4">
          <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">
            <PlannerHeader
              mapName={selectedMap?.displayName ?? plan.mapId}
              busy={busy}
              onSave={() => void handleSave()}
              onReset={() => void handleReset()}
              onClearDraft={() => {
                updatePlan({ markers: [], lines: [] })
                setLineStart(null)
                setPendingAbility(null)
                setMessage('Đã xóa nháp (chưa lưu).')
              }}
            />
            <PlannerForm
              plan={plan}
              tool={tool}
              markerKind={markerKind}
              showCallouts={showCallouts}
              hasTactical={selectedMap?.hasTacticalLayout ?? false}
              onToolChange={(t) => {
                setTool(t)
                setLineStart(null)
                setPendingAbility(null)
                if (t !== 'agent') setSelectedAgent(null)
              }}
              onMarkerKindChange={setMarkerKind}
              onToggleCallouts={() => setShowCallouts((v) => !v)}
              onPlanPatch={updatePlan}
            />
            <textarea
              value={plan.notes}
              onChange={(event) => updatePlan({ notes: event.target.value })}
              placeholder="Ghi chú: timing smoke, default call, rotate..."
              className="mt-3 min-h-[72px] w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan"
            />

            {/* Pending ability banner */}
            {pendingAbility && (
              <div className="mt-3 flex items-center justify-between rounded-xl border border-tactical-cyan/30 bg-tactical-cyan/10 px-3 py-2">
                <div className="flex items-center gap-2">
                  {pendingAbility.ability.displayIcon && (
                    <img src={pendingAbility.ability.displayIcon} alt="" className="h-5 w-5 object-contain" />
                  )}
                  <span className="text-xs font-bold text-tactical-cyan">
                    Đang đặt: {pendingAbility.ability.displayName}
                  </span>
                  <span className="text-[10px] text-slate-500">
                    ({pendingAbility.sourceMarker.label})
                  </span>
                </div>
                <button
                  type="button"
                  onClick={() => { setPendingAbility(null); setMessage('') }}
                  className="text-xs text-slate-500 hover:text-white"
                >
                  Hủy (Esc)
                </button>
              </div>
            )}

            {message && !pendingAbility && (
              <p className="mt-2 text-sm font-semibold text-tactical-cyan">{message}</p>
            )}
            {lineStart && tool === 'line' && (
              <button
                type="button"
                onClick={() => { setLineStart(null); setMessage('') }}
                className="mt-2 text-xs font-bold uppercase tracking-[0.14em] text-slate-400 hover:text-white"
              >
                Hủy route đang vẽ
              </button>
            )}
          </div>

          <TacticalBoard
            plan={plan}
            selectedMap={selectedMap}
            lineStart={lineStart}
            boardRef={boardRef}
            showCallouts={showCallouts}
            onPlanPatch={updatePlan}
            onBoardClick={handleBoardClick}
            onAgentRightClick={handleAgentRightClick}
            kindColor={kindColor}
          />

          <div className="rounded-2xl border border-tactical-line bg-black/20 px-4 py-3 text-xs text-slate-400">
            {otherCount} marker · {agentCount} agent · {abilityCount} chiêu · {plan.lines.length} route
            {plan.updatedAt ? ` · Lưu: ${new Date(plan.updatedAt).toLocaleString()}` : ' · Chưa lưu'}
          </div>
        </div>
      </section>
    </>
  )
}

const PlannerHeader: FC<{
  mapName: string
  busy: boolean
  onSave: () => void
  onReset: () => void
  onClearDraft: () => void
}> = ({ mapName, busy, onSave, onReset, onClearDraft }) => (
  <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
    <div>
      <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">Kế hoạch trước trận</p>
      <p className="mt-1 text-2xl font-black text-white">{mapName}</p>
      <p className="mt-2 text-sm text-slate-400">
        Click: đặt marker/route/agent · <span className="text-slate-300">Chuột phải agent trên map → chọn chiêu</span>
      </p>
    </div>
    <div className="flex flex-wrap gap-2">
      <button type="button" disabled={busy} onClick={onSave} className="rounded-2xl bg-tactical-red px-4 py-2 text-xs font-black uppercase tracking-[0.14em] text-white disabled:opacity-50">
        Lưu
      </button>
      <button type="button" disabled={busy} onClick={onClearDraft} className="rounded-2xl border border-white/10 px-4 py-2 text-xs font-black uppercase tracking-[0.14em] text-slate-300">
        Xóa nháp
      </button>
      <button type="button" disabled={busy} onClick={onReset} className="rounded-2xl border border-tactical-red/30 px-4 py-2 text-xs font-black uppercase tracking-[0.14em] text-tactical-red">
        Xóa đã lưu
      </button>
    </div>
  </div>
)

const snapToCallout = (x: number, y: number, callouts: MapCallout[]) => {
  let best: MapCallout | null = null
  let bestDistance = 5
  for (const callout of callouts) {
    const distance = Math.hypot(callout.x - x, callout.y - y)
    if (distance < bestDistance) {
      best = callout
      bestDistance = distance
    }
  }
  return best
}

const PlannerForm: FC<{
  plan: MapPlan
  tool: PlannerTool
  markerKind: MarkerKind
  showCallouts: boolean
  hasTactical: boolean
  onToolChange: (tool: PlannerTool) => void
  onMarkerKindChange: (kind: MarkerKind) => void
  onToggleCallouts: () => void
  onPlanPatch: (patch: Partial<MapPlan>) => void
}> = ({ plan, tool, markerKind, showCallouts, hasTactical, onToolChange, onMarkerKindChange, onToggleCallouts, onPlanPatch }) => (
  <div className="mt-4 space-y-3">
    <div className="grid gap-3 sm:grid-cols-2">
    <Field label="Tên kế hoạch">
      <input
        value={plan.title}
        onChange={(event) => onPlanPatch({ title: event.target.value })}
        className="w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan"
      />
    </Field>
    <Field label="Phe">
      <select
        value={plan.side}
        onChange={(event) => onPlanPatch({ side: event.target.value })}
        className="w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan"
      >
        <option value="attack">Tấn công</option>
        <option value="defense">Phòng thủ</option>
        <option value="both">Cả hai</option>
      </select>
    </Field>
    </div>

    <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end">
      <Field label="Công cụ" className="w-full sm:w-auto sm:min-w-[min(100%,17rem)] sm:flex-1">
        <ToolButtons tool={tool} onToolChange={onToolChange} />
      </Field>
      {tool === 'marker' && (
      <Field label="Loại marker" className="w-full sm:w-auto sm:min-w-[min(100%,11rem)] sm:max-w-xs sm:flex-1">
        <select
          value={markerKind}
          onChange={(event) => onMarkerKindChange(event.target.value as MarkerKind)}
          className="w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan"
        >
          {markerKinds.map((item) => (
            <option key={item.id} value={item.id}>
              {item.label}
            </option>
          ))}
        </select>
      </Field>
      )}
      {hasTactical && (
      <Field label="Callout" className="w-full sm:w-auto sm:min-w-[min(100%,9rem)]">
        <button
          type="button"
          onClick={onToggleCallouts}
          className={`w-full rounded-xl border px-3 py-2 text-xs font-black uppercase tracking-[0.12em] ${
            showCallouts ? 'border-tactical-cyan/40 bg-tactical-cyan/10 text-tactical-cyan' : 'border-white/10 text-slate-400'
          }`}
        >
          {showCallouts ? 'Ẩn callout' : 'Hiện callout'}
        </button>
      </Field>
      )}
    </div>
  </div>
)

const ToolButtons: FC<{ tool: PlannerTool; onToolChange: (tool: PlannerTool) => void }> = ({ tool, onToolChange }) => (
  <div className="flex flex-wrap gap-2">
    <ToolButton active={tool === 'marker'} label="Marker" onClick={() => onToolChange('marker')} />
    <ToolButton active={tool === 'line'} label="Route" onClick={() => onToolChange('line')} />
    <ToolButton active={tool === 'agent'} label="Agent" onClick={() => onToolChange('agent')} />
  </div>
)

const Field: FC<{ label: string; children: ReactNode; className?: string }> = ({ label, children, className }) => (
  <label className={`block min-w-0 ${className ?? ''}`}>
    <span className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">{label}</span>
    <div className="mt-2 min-w-0">{children}</div>
  </label>
)

const ToolButton: FC<{ active: boolean; label: string; onClick: () => void }> = ({ active, label, onClick }) => (
  <button
    type="button"
    onClick={onClick}
    className={`min-w-[4.5rem] flex-1 rounded-xl border px-3 py-2 text-xs font-black uppercase tracking-[0.12em] ${
      active ? 'border-tactical-cyan/40 bg-tactical-cyan/10 text-tactical-cyan' : 'border-white/10 text-slate-400'
    }`}
  >
    {label}
  </button>
)

export default MapPlannerPanel
