import { useCallback, useEffect, useRef, useState } from 'react'
import type { FC, MouseEvent, ReactNode, RefObject } from 'react'
import { deleteMapPlan, listTacticalMaps, loadMapPlan, saveMapPlan } from './api'
import { fallbackMapCatalog, mapIdFromName } from './map-catalog-fallback'
import TacticalBoard from './tactical-board'
import type { MapCatalogEntry, MapCallout, MapPlan, MarkerKind, PlanLine, PlanMarker, PlannerTool } from './types'
import { getMapCallouts } from './map-callouts'

type MapPlannerPanelProps = {
  suggestedMapName?: string
}

const markerKinds: Array<{ id: MarkerKind; label: string; color: string }> = [
  { id: 'duelist', label: 'Đối Đầu (Duelist)', color: '#ff4655' },
  { id: 'initiator', label: 'Khởi Tranh (Initiator)', color: '#fbbf24' },
  { id: 'controller', label: 'Kiểm Soát (Controller)', color: '#42e8f3' },
  { id: 'sentinel', label: 'Hộ Vệ (Sentinel)', color: '#4ade80' },
  { id: 'callout', label: 'Callout', color: '#e2e8f0' },
]

const markerLabel = (kind: MarkerKind) => markerKinds.find((item) => item.id === kind)?.label ?? kind

const kindColor = (kind: string) => markerKinds.find((item) => item.id === kind)?.color ?? '#42e8f3'

const catalogHasMap = (name: string) => {
  const id = mapIdFromName(name)
  return fallbackMapCatalog.some((map) => map.id === id || map.name.toLowerCase() === name.toLowerCase())
}

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
    setPlan(await loadMapPlan(mapId))
  }

  const updatePlan = (patch: Partial<MapPlan>) => {
    setPlan((current) => (current ? { ...current, ...patch, mapId: selectedMapId } : current))
  }

  const pointerPercent = (event: MouseEvent<HTMLDivElement>) => {
    const rect = boardRef.current?.getBoundingClientRect()
    if (!rect?.width || !rect.height) {
      return null
    }
    return {
      x: Math.min(100, Math.max(0, Number((((event.clientX - rect.left) / rect.width) * 100).toFixed(2)))),
      y: Math.min(100, Math.max(0, Number((((event.clientY - rect.top) / rect.height) * 100).toFixed(2)))),
    }
  }

  const handleBoardClick = (event: MouseEvent<HTMLDivElement>) => {
    if (!plan) {
      return
    }
    let point = pointerPercent(event)
    if (!point) {
      return
    }

    const snapped = snapToCallout(point.x, point.y, getMapCallouts(selectedMapId))
    if (snapped) {
      point = { x: snapped.x, y: snapped.y }
    }

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

  const handleSave = async () => {
    if (!plan) {
      return
    }
    setBusy(true)
    try {
      setPlan(await saveMapPlan({ ...plan, mapId: selectedMapId }))
      setMessage(`Đã lưu kế hoạch ${selectedMap?.displayName ?? selectedMapId}.`)
    } finally {
      setBusy(false)
    }
  }

  const handleReset = async () => {
    if (!window.confirm('Xóa kế hoạch đã lưu cho map này?')) {
      return
    }
    setBusy(true)
    try {
      await deleteMapPlan(selectedMapId)
      setPlan(await loadMapPlan(selectedMapId))
      setLineStart(null)
      setMessage('Đã xóa kế hoạch map.')
    } finally {
      setBusy(false)
    }
  }

  if (!plan) {
    return <div className="rounded-3xl border border-tactical-line bg-tactical-panel p-6 text-slate-300">Đang tải bản đồ...</div>
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[280px_1fr]">
      <aside className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-4 shadow-xl shadow-black/20">
        <p className="text-xs font-bold uppercase tracking-[0.2em] text-tactical-cyan">Map pool</p>
        <p className="mt-1 text-sm text-slate-400">12 map competitive · tactical mô phỏng · lưu local</p>
        {suggestedMapName && (
          <p className="mt-2 rounded-xl border border-tactical-red/30 bg-tactical-red/10 px-3 py-2 text-xs text-tactical-red">
            Coach gợi ý: {suggestedMapName}
          </p>
        )}
        <div className="mt-3 max-h-[520px] space-y-1 overflow-y-auto pr-1">
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
              <img
                src={map.imageUrl}
                alt=""
                className="h-10 w-16 rounded-lg object-cover bg-black/40"
              />
              <span className="text-sm font-bold text-white">
                {map.displayName}
                {map.hasTacticalLayout && <span className="ml-1 text-[10px] text-tactical-cyan">TAC</span>}
              </span>
            </button>
          ))}
        </div>
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
              setMessage('Đã xóa nháp (chưa lưu).')
            }}
          />
          <PlannerForm
            plan={plan}
            tool={tool}
            markerKind={markerKind}
            showCallouts={showCallouts}
            hasTactical={selectedMap?.hasTacticalLayout ?? false}
            onToolChange={setTool}
            onMarkerKindChange={setMarkerKind}
            onToggleCallouts={() => setShowCallouts((value) => !value)}
            onPlanPatch={updatePlan}
          />
          <textarea
            value={plan.notes}
            onChange={(event) => updatePlan({ notes: event.target.value })}
            placeholder="Ghi chú: timing smoke, default call, rotate..."
            className="mt-3 min-h-[72px] w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan"
          />
          {message && <p className="mt-2 text-sm font-semibold text-tactical-cyan">{message}</p>}
          {lineStart && tool === 'line' && (
            <button
              type="button"
              onClick={() => {
                setLineStart(null)
                setMessage('')
              }}
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
          kindColor={kindColor}
        />

        <div className="rounded-2xl border border-tactical-line bg-black/20 px-4 py-3 text-xs text-slate-400">
          {plan.markers.length} marker · {plan.lines.length} route
          {plan.updatedAt ? ` · Lưu: ${new Date(plan.updatedAt).toLocaleString()}` : ' · Chưa lưu'}
        </div>
      </div>
    </section>
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
      <p className="mt-2 text-sm text-slate-400">Click map: marker hoặc route (2 click).</p>
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
  <div className="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
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
        <option value="attack">Attack</option>
        <option value="defense">Defense</option>
        <option value="both">Both</option>
      </select>
    </Field>
    <Field label="Công cụ">
      <ToolButtons tool={tool} onToolChange={onToolChange} />
    </Field>
    <Field label="Loại marker">
      <select
        value={markerKind}
        disabled={tool !== 'marker'}
        onChange={(event) => onMarkerKindChange(event.target.value as MarkerKind)}
        className="w-full rounded-xl border border-tactical-line bg-black/30 px-3 py-2 text-sm text-white outline-none focus:border-tactical-cyan disabled:opacity-50"
      >
        {markerKinds.map((item) => (
          <option key={item.id} value={item.id}>
            {item.label}
          </option>
        ))}
      </select>
    </Field>
    {hasTactical && (
      <Field label="Callout">
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
)

const ToolButtons: FC<{ tool: PlannerTool; onToolChange: (tool: PlannerTool) => void }> = ({ tool, onToolChange }) => (
  <div className="flex gap-2">
    <ToolButton active={tool === 'marker'} label="Marker" onClick={() => onToolChange('marker')} />
    <ToolButton active={tool === 'line'} label="Route" onClick={() => onToolChange('line')} />
  </div>
)

const Field: FC<{ label: string; children: ReactNode }> = ({ label, children }) => (
  <label className="block">
    <span className="text-xs font-bold uppercase tracking-[0.18em] text-slate-500">{label}</span>
    <div className="mt-2">{children}</div>
  </label>
)

const ToolButton: FC<{ active: boolean; label: string; onClick: () => void }> = ({ active, label, onClick }) => (
  <button
    type="button"
    onClick={onClick}
    className={`flex-1 rounded-xl border px-3 py-2 text-xs font-black uppercase tracking-[0.12em] ${
      active ? 'border-tactical-cyan/40 bg-tactical-cyan/10 text-tactical-cyan' : 'border-white/10 text-slate-400'
    }`}
  >
    {label}
  </button>
)

export default MapPlannerPanel
