import { useEffect, useState } from 'react'
import type { FC, MouseEvent, ReactNode, RefObject } from 'react'
import { getMapCallouts } from './map-callouts'
import type { MapCallout, MapCatalogEntry, MapPlan, PlanMarker } from './types'

type TacticalBoardProps = {
  plan: MapPlan
  selectedMap?: MapCatalogEntry
  lineStart: { x: number; y: number } | null
  boardRef: RefObject<HTMLDivElement | null>
  showCallouts: boolean
  onPlanPatch: (patch: Partial<MapPlan>) => void
  onBoardClick: (event: MouseEvent<HTMLDivElement>) => void
  kindColor: (kind: string) => string
}

const TacticalBoard: FC<TacticalBoardProps> = ({
  plan,
  selectedMap,
  lineStart,
  boardRef,
  showCallouts,
  onPlanPatch,
  onBoardClick,
  kindColor,
}) => {
  const [imageRatio, setImageRatio] = useState<number | null>(null)
  const [imageError, setImageError] = useState(false)
  const mapId = selectedMap?.id ?? plan.mapId
  const mapName = selectedMap?.displayName ?? mapId
  const callouts = getMapCallouts(mapId)
  const tacticalSrc = selectedMap?.tacticalImageUrl ?? ''

  useEffect(() => {
    setImageRatio(null)
    setImageError(false)
  }, [mapId])

  return (
    <div>
      <p className="mb-2 text-xs text-slate-400">
        Minimap chính thức của Riot · click callout để ghim marker · click trên nền để vẽ marker hoặc route
      </p>
      <div
        ref={boardRef}
        onClick={onBoardClick}
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

        {/* Overlay tối nhẹ để callout dễ đọc trên minimap sáng */}
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

        {plan.markers.map((marker) => (
          <button
            key={marker.id}
            type="button"
            title={`${marker.label || marker.kind} — click xóa`}
            onClick={(event) => {
              event.stopPropagation()
              onPlanPatch({ markers: plan.markers.filter((item) => item.id !== marker.id) })
            }}
            style={{
              left: `${marker.x}%`,
              top: `${marker.y}%`,
              borderColor: kindColor(marker.kind),
              color: kindColor(marker.kind),
            }}
            className="absolute z-20 max-w-[120px] -translate-x-1/2 -translate-y-1/2 truncate rounded-lg border-2 bg-black/85 px-2 py-1 text-[10px] font-black uppercase shadow-lg"
          >
            {marker.label ? marker.label.slice(0, 14) : marker.kind.slice(0, 1)}
          </button>
        ))}
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
