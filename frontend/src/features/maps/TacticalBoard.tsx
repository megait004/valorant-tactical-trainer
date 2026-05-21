import { useEffect, useRef, useState } from 'react'
import type {
  FC,
  MouseEvent,
  PointerEvent as ReactPointerEvent,
  ReactNode,
  RefObject,
} from 'react'
import { getMapCallouts } from './callouts'
import type { MapCallout, MapCatalogEntry, MapPlan, PlanMarker } from '../../types'
import { ROLE_COLORS } from '../agents/useAgents'

type TacticalBoardProps = {
  plan: MapPlan
  selectedMap?: MapCatalogEntry
  lineStart: { x: number; y: number } | null
  boardRef: RefObject<HTMLDivElement | null>
  showCallouts: boolean
  onPlanPatch: (patch: Partial<MapPlan>) => void
  onBoardClick: (event: MouseEvent<HTMLDivElement>) => void
  onAgentRightClick?: (marker: PlanMarker, clientX: number, clientY: number) => void
  kindColor: (kind: string, agentRole?: string) => string
}

// Threshold (px) before a pointer-down is treated as a drag instead of a click.
const DRAG_THRESHOLD = 4

type DragHandlers = {
  draggingId: string | null
  startDrag: (event: ReactPointerEvent<HTMLButtonElement>, markerId: string) => void
}

// ─── Ability Marker ───────────────────────────────────────────────────────────

const AbilityMarker: FC<{
  marker: PlanMarker
  isDragging: boolean
  onPointerDown: (event: ReactPointerEvent<HTMLButtonElement>) => void
  onRemove: () => void
}> = ({ marker, isDragging, onPointerDown, onRemove }) => {
  const roleColor = ROLE_COLORS[marker.agentRole ?? ''] ?? '#42e8f3'
  return (
    <button
      type="button"
      title={`${marker.label} · kéo để di chuyển · click để xóa`}
      onPointerDown={onPointerDown}
      onClick={(e) => {
        e.stopPropagation()
        onRemove()
      }}
      onContextMenu={(e) => e.stopPropagation()}
      style={{ left: `${marker.x}%`, top: `${marker.y}%`, touchAction: 'none' }}
      className={`absolute z-20 -translate-x-1/2 -translate-y-1/2 group ${
        isDragging ? 'cursor-grabbing' : 'cursor-grab'
      }`}
    >
      <div
        className={`flex h-8 w-8 items-center justify-center rounded-lg border-2 border-dashed bg-black/80 shadow-lg transition group-hover:scale-110 ${
          isDragging ? 'scale-110 ring-2 ring-white/50' : ''
        }`}
        style={{ borderColor: roleColor, boxShadow: `0 0 6px ${roleColor}50` }}
      >
        {marker.agentPortrait ? (
          <img
            src={marker.agentPortrait}
            alt={marker.label}
            className="h-5 w-5 object-contain"
            draggable={false}
          />
        ) : (
          <span className="text-[9px] font-black" style={{ color: roleColor }}>
            {marker.abilitySlot?.slice(0, 1) ?? '?'}
          </span>
        )}
        <div className="pointer-events-none absolute inset-0 hidden items-center justify-center rounded-lg bg-black/70 group-hover:flex">
          <span className="text-xs font-black text-white">✕</span>
        </div>
      </div>
      <div
        className="mt-0.5 max-w-[80px] truncate rounded px-1 py-0.5 text-center text-[8px] font-black uppercase leading-tight"
        style={{ backgroundColor: `${roleColor}22`, color: roleColor }}
      >
        {marker.label.slice(0, 12)}
      </div>
    </button>
  )
}

// ─── Agent Marker ─────────────────────────────────────────────────────────────

const AgentMarker: FC<{
  marker: PlanMarker
  isDragging: boolean
  onPointerDown: (event: ReactPointerEvent<HTMLButtonElement>) => void
  onRemove: () => void
  onRightClick: (clientX: number, clientY: number) => void
}> = ({ marker, isDragging, onPointerDown, onRemove, onRightClick }) => {
  const roleColor = ROLE_COLORS[marker.agentRole ?? ''] ?? '#42e8f3'
  return (
    <button
      type="button"
      title={`${marker.label} · kéo để di chuyển · click để xóa · chuột phải → chọn chiêu`}
      onPointerDown={onPointerDown}
      onClick={(e) => {
        e.stopPropagation()
        onRemove()
      }}
      onContextMenu={(e) => {
        e.preventDefault()
        e.stopPropagation()
        onRightClick(e.clientX, e.clientY)
      }}
      style={{ left: `${marker.x}%`, top: `${marker.y}%`, touchAction: 'none' }}
      className={`absolute z-20 -translate-x-1/2 -translate-y-1/2 group ${
        isDragging ? 'cursor-grabbing' : 'cursor-grab'
      }`}
    >
      <div
        className={`relative flex h-10 w-10 items-center justify-center rounded-full border-2 bg-black shadow-lg transition group-hover:scale-110 ${
          isDragging ? 'scale-110 ring-2 ring-white/60' : ''
        }`}
        style={{ borderColor: roleColor, boxShadow: `0 0 10px ${roleColor}60` }}
      >
        {marker.agentPortrait ? (
          <img
            src={marker.agentPortrait}
            alt={marker.label}
            className="h-full w-full rounded-full object-cover"
            draggable={false}
          />
        ) : (
          <span className="text-xs font-black" style={{ color: roleColor }}>
            {marker.label.slice(0, 2)}
          </span>
        )}
        {/* Right-click hint ring on hover */}
        <div className="pointer-events-none absolute -inset-1 hidden rounded-full ring-2 ring-white/30 ring-offset-1 ring-offset-transparent group-hover:block" />
      </div>
      <div
        className="mt-0.5 max-w-[72px] truncate rounded px-1 py-0.5 text-center text-[8px] font-black uppercase leading-tight"
        style={{ backgroundColor: `${roleColor}22`, color: roleColor }}
      >
        {marker.label}
      </div>
    </button>
  )
}

// ─── Tactical Board ───────────────────────────────────────────────────────────

const TacticalBoard: FC<TacticalBoardProps> = ({
  plan,
  selectedMap,
  lineStart,
  boardRef,
  showCallouts,
  onPlanPatch,
  onBoardClick,
  onAgentRightClick,
  kindColor,
}) => {
  const [imageRatio, setImageRatio] = useState<number | null>(null)
  const [imageError, setImageError] = useState(false)
  const [draggingId, setDraggingId] = useState<string | null>(null)

  const mapId = selectedMap?.id ?? plan.mapId
  const mapName = selectedMap?.displayName ?? mapId
  const callouts = getMapCallouts(mapId)
  const tacticalSrc = selectedMap?.tacticalImageUrl ?? ''

  // Refs to access latest data inside the global pointer listeners without
  // re-binding them on every render.
  const markersRef = useRef(plan.markers)
  const onPlanPatchRef = useRef(onPlanPatch)
  const suppressClickRef = useRef<string | null>(null)

  useEffect(() => {
    markersRef.current = plan.markers
  }, [plan.markers])

  useEffect(() => {
    onPlanPatchRef.current = onPlanPatch
  }, [onPlanPatch])

  useEffect(() => {
    setImageRatio(null)
    setImageError(false)
  }, [mapId])

  const removeMarker = (markerId: string) => {
    if (suppressClickRef.current === markerId) {
      suppressClickRef.current = null
      return
    }
    onPlanPatch({ markers: plan.markers.filter((m) => m.id !== markerId) })
  }

  const startDrag = (event: ReactPointerEvent<HTMLButtonElement>, markerId: string) => {
    if (event.button !== 0) return
    event.stopPropagation()

    const startX = event.clientX
    const startY = event.clientY
    let moved = false

    const handleMove = (ev: PointerEvent) => {
      const dx = ev.clientX - startX
      const dy = ev.clientY - startY
      if (!moved && Math.hypot(dx, dy) > DRAG_THRESHOLD) {
        moved = true
        setDraggingId(markerId)
      }
      if (!moved) return

      const rect = boardRef.current?.getBoundingClientRect()
      if (!rect?.width || !rect.height) return
      const x = Math.min(
        100,
        Math.max(0, Number((((ev.clientX - rect.left) / rect.width) * 100).toFixed(2))),
      )
      const y = Math.min(
        100,
        Math.max(0, Number((((ev.clientY - rect.top) / rect.height) * 100).toFixed(2))),
      )
      const updated = markersRef.current.map((m) =>
        m.id === markerId ? { ...m, x, y } : m,
      )
      onPlanPatchRef.current({ markers: updated })
    }

    const handleUp = () => {
      window.removeEventListener('pointermove', handleMove)
      window.removeEventListener('pointerup', handleUp)
      window.removeEventListener('pointercancel', handleUp)
      if (moved) {
        // Suppress the click that follows a drag so we don't delete the marker.
        suppressClickRef.current = markerId
        setTimeout(() => {
          if (suppressClickRef.current === markerId) suppressClickRef.current = null
        }, 0)
      }
      setDraggingId(null)
    }

    window.addEventListener('pointermove', handleMove)
    window.addEventListener('pointerup', handleUp)
    window.addEventListener('pointercancel', handleUp)
  }

  const drag: DragHandlers = { draggingId, startDrag }

  return (
    <div>
      <p className="mb-2 text-xs text-slate-400">
        Click: đặt marker/route/agent · Kéo marker để di chuyển · Chuột phải agent: chọn chiêu
      </p>
      <div
        ref={boardRef}
        role="button"
        tabIndex={0}
        onClick={onBoardClick}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ')
            onBoardClick(e as unknown as MouseEvent<HTMLDivElement>)
        }}
        style={imageRatio ? { aspectRatio: String(imageRatio) } : undefined}
        className={`relative w-full cursor-crosshair overflow-hidden rounded-3xl border border-tactical-line bg-[#08090d] shadow-2xl shadow-black/30 ${
          imageRatio ? '' : 'aspect-square'
        }`}
      >
        {tacticalSrc && !imageError ? (
          <img
            src={tacticalSrc}
            alt={`Minimap ${mapName}`}
            className="absolute inset-0 h-full w-full object-cover"
            draggable={false}
            onLoad={(event) => {
              const img = event.currentTarget
              if (img.naturalWidth && img.naturalHeight) {
                setImageRatio(img.naturalWidth / img.naturalHeight)
              }
            }}
            onError={() => setImageError(true)}
          />
        ) : (
          <div className="absolute inset-0 flex items-center justify-center text-sm text-slate-500">
            Không tải được minimap {mapName}
          </div>
        )}

        <div className="pointer-events-none absolute inset-0 bg-black/15" />

        {showCallouts && (
          <CalloutLayer
            callouts={callouts}
            onPickCallout={(callout) => {
              const marker: PlanMarker = {
                id: `callout-${callout.id}-${Date.now()}`,
                kind: 'callout',
                label: callout.label,
                x: callout.x,
                y: callout.y,
              }
              onPlanPatch({ markers: [...plan.markers, marker] })
            }}
          />
        )}

        <svg className="pointer-events-none absolute inset-0 h-full w-full">
          {plan.lines.map((line) => (
            <line
              key={line.id}
              x1={`${line.x1}%`}
              y1={`${line.y1}%`}
              x2={`${line.x2}%`}
              y2={`${line.y2}%`}
              stroke="#42e8f3"
              strokeWidth="3"
              strokeOpacity="0.9"
              markerEnd="url(#tactical-arrow)"
            />
          ))}
          <defs>
            <marker id="tactical-arrow" markerWidth="8" markerHeight="8" refX="6" refY="3" orient="auto">
              <polygon points="0 0, 8 3, 0 6" fill="#42e8f3" />
            </marker>
          </defs>
          {lineStart && <circle cx={`${lineStart.x}%`} cy={`${lineStart.y}%`} r="6" fill="#ff4655" />}
        </svg>

        {plan.markers.map((marker) => {
          const isDragging = drag.draggingId === marker.id
          if (marker.kind === 'agent') {
            return (
              <AgentMarker
                key={marker.id}
                marker={marker}
                isDragging={isDragging}
                onPointerDown={(e) => drag.startDrag(e, marker.id)}
                onRemove={() => removeMarker(marker.id)}
                onRightClick={(cx, cy) => onAgentRightClick?.(marker, cx, cy)}
              />
            )
          }
          if (marker.kind === 'ability') {
            return (
              <AbilityMarker
                key={marker.id}
                marker={marker}
                isDragging={isDragging}
                onPointerDown={(e) => drag.startDrag(e, marker.id)}
                onRemove={() => removeMarker(marker.id)}
              />
            )
          }
          return (
            <button
              key={marker.id}
              type="button"
              title={`${marker.label || marker.kind} · kéo để di chuyển · click để xóa`}
              onPointerDown={(e) => drag.startDrag(e, marker.id)}
              onClick={(event) => {
                event.stopPropagation()
                removeMarker(marker.id)
              }}
              style={{
                left: `${marker.x}%`,
                top: `${marker.y}%`,
                borderColor: kindColor(marker.kind, marker.agentRole),
                color: kindColor(marker.kind, marker.agentRole),
                touchAction: 'none',
              }}
              className={`absolute z-20 max-w-[120px] -translate-x-1/2 -translate-y-1/2 truncate rounded-lg border-2 bg-black/85 px-2 py-1 text-[10px] font-black uppercase shadow-lg ${
                isDragging ? 'cursor-grabbing scale-110 ring-2 ring-white/40' : 'cursor-grab'
              }`}
            >
              {marker.label ? marker.label.slice(0, 14) : marker.kind.slice(0, 1)}
            </button>
          )
        })}
      </div>
    </div>
  )
}

const CalloutLayer: FC<{
  callouts: MapCallout[]
  onPickCallout: (callout: MapCallout) => void
}> = ({ callouts, onPickCallout }) => (
  <CalloutRoot>
    {callouts.map((callout) => (
      <button
        key={callout.id}
        type="button"
        style={{ left: `${callout.x}%`, top: `${callout.y}%` }}
        title={callout.label}
        onClick={(event) => {
          event.stopPropagation()
          onPickCallout(callout)
        }}
        className={calloutClass(callout.kind)}
      >
        <span className={callout.kind === 'site' ? 'relative z-10' : 'block truncate'}>
          {callout.label}
        </span>
      </button>
    ))}
  </CalloutRoot>
)

const CalloutRoot: FC<{ children: ReactNode }> = ({ children }) => (
  <div className="pointer-events-none absolute inset-0 z-[5]">{children}</div>
)

const calloutClass = (kind: MapCallout['kind']) => {
  const base = 'pointer-events-auto absolute z-[6] -translate-x-1/2 -translate-y-1/2'
  if (kind === 'site') {
    return `${base} flex h-10 w-10 items-center justify-center rounded-full bg-tactical-red text-lg font-black text-white shadow-lg shadow-tactical-red/50 ring-2 ring-white/40`
  }
  if (kind === 'spawn') {
    return `${base} max-w-[170px] rounded-md border border-white/40 bg-black/80 px-2 py-1 text-[9px] font-black uppercase leading-tight text-slate-100 shadow-md shadow-black/40 backdrop-blur-sm`
  }
  if (kind === 'mid') {
    return `${base} max-w-[120px] rounded border border-violet-300/70 bg-black/75 px-1.5 py-0.5 text-[8px] font-black uppercase leading-tight text-violet-100 shadow shadow-black/40 backdrop-blur-sm`
  }
  return `${base} max-w-[120px] rounded border border-cyan-200/50 bg-black/75 px-1.5 py-0.5 text-[8px] font-black uppercase leading-tight text-slate-100 shadow shadow-black/40 backdrop-blur-sm`
}

export default TacticalBoard
