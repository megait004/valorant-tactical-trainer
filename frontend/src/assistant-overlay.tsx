import type { FC } from 'react'
import type { AssistantSessionState } from './types'

type AssistantOverlayProps = {
  session: AssistantSessionState
  onRequestTip: () => void
  onRoundStart: () => void
  onStop: () => void
  onExitOverlay: () => void
}

const severityClass = (severity: string) => {
  if (severity === 'high') return 'border-tactical-red/50 bg-tactical-red/15 text-tactical-red'
  if (severity === 'medium') return 'border-tactical-cyan/40 bg-tactical-cyan/10 text-tactical-cyan'
  return 'border-white/15 bg-white/[0.04] text-slate-200'
}

const AssistantOverlay: FC<AssistantOverlayProps> = ({
  session,
  onRequestTip,
  onRoundStart,
  onStop,
  onExitOverlay,
}) => {
  const alert = session.currentAlert

  return (
    <main className="h-screen overflow-hidden bg-tactical-bg/95 p-3 text-slate-100 backdrop-blur-sm">
      <LiveBadge />

      {alert ? (
        <div className={`mt-2 rounded-2xl border px-3 py-3 ${severityClass(alert.severity)}`}>
          <p className="text-[10px] font-bold uppercase tracking-[0.2em] opacity-80">Live Assistant</p>
          <p className="mt-1 text-sm font-black leading-snug text-white">{alert.title}</p>
          <p className="mt-2 text-xs leading-5 opacity-90">{alert.message}</p>
        </div>
      ) : (
        <p className="mt-2 text-xs text-slate-400">{session.message}</p>
      )}

      <div className="mt-2 flex gap-3 text-[10px] text-slate-500">
        <span>Round {session.roundCount}</span>
        <span>Tips {session.tipsShown}</span>
        <span>Queue {session.queueSize}</span>
      </div>

      <OverlayActions
        onRoundStart={onRoundStart}
        onRequestTip={onRequestTip}
        onExitOverlay={onExitOverlay}
        onStop={onStop}
      />
    </main>
  )
}

const LiveBadge: FC = () => (
  <div className="flex items-center justify-between gap-2">
    <div className="flex items-center gap-2">
      <span className="relative flex h-2 w-2">
        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-tactical-red opacity-60" />
        <span className="relative inline-flex h-2 w-2 rounded-full bg-tactical-red" />
      </span>
      <span className="text-[10px] font-bold uppercase tracking-[0.22em] text-tactical-red">Đang luyện</span>
    </div>
    <span className="text-[10px] text-slate-500">Safe · read-only</span>
  </div>
)

const OverlayActions: FC<{
  onRoundStart: () => void
  onRequestTip: () => void
  onExitOverlay: () => void
  onStop: () => void
}> = ({ onRoundStart, onRequestTip, onExitOverlay, onStop }) => (
  <div className="mt-3 grid grid-cols-2 gap-2">
    <OverlayButton label="Round mới" onClick={onRoundStart} variant="cyan" />
    <OverlayButton label="Nhắc tôi" onClick={onRequestTip} variant="neutral" />
    <OverlayButton label="Mở app" onClick={onExitOverlay} variant="neutral" />
    <OverlayButton label="Tắt" onClick={onStop} variant="red" />
  </div>
)

const OverlayButton: FC<{
  label: string
  onClick: () => void
  variant: 'cyan' | 'red' | 'neutral'
}> = ({ label, onClick, variant }) => {
  const classes =
    variant === 'cyan'
      ? 'border-tactical-cyan/30 bg-tactical-cyan/10 text-tactical-cyan'
      : variant === 'red'
        ? 'border-tactical-red/30 bg-tactical-red/10 text-tactical-red'
        : 'border-white/10 bg-white/[0.04] text-slate-200'

  return (
    <button
      type="button"
      onClick={onClick}
      className={`rounded-xl border px-2 py-2 text-[10px] font-black uppercase tracking-[0.12em] ${classes}`}
    >
      {label}
    </button>
  )
}

export default AssistantOverlay
