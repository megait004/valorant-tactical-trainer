import type { FC } from 'react'
import type { AssistantSessionState } from '../../types'

type LiveAssistantPanelProps = {
  session: AssistantSessionState
  overlayMode: boolean
  busy: boolean
  onStart: () => void
  onStop: () => void
  onRequestTip: () => void
  onRoundStart: () => void
  onToggleOverlay: () => void
}

const severityClass = (severity: string) => {
  if (severity === 'high') return 'border-tactical-red/40 bg-tactical-red/10'
  if (severity === 'medium') return 'border-tactical-cyan/30 bg-tactical-cyan/10'
  return 'border-white/10 bg-white/[0.03]'
}

const LiveAssistantPanel: FC<LiveAssistantPanelProps> = ({
  session,
  overlayMode,
  busy,
  onStart,
  onStop,
  onRequestTip,
  onRoundStart,
  onToggleOverlay,
}) => {
  const alert = session.currentAlert

  return (
    <section className="grid gap-5 lg:grid-cols-[1fr_0.9fr]">
      <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">
        <p className="text-xs font-bold uppercase tracking-[0.2em] text-tactical-red">Live Assistant</p>
        <p className="mt-2 text-2xl font-black text-white">Trợ lý cảnh báo khi đang chơi</p>
        <p className="mt-3 text-sm leading-6 text-slate-300">
          Dùng profile từ Coach (report Henrik hoặc demo). Bật overlay, kéo cửa sổ góc màn hình, vào Valorant — app nhắc theo thói quen cần sửa.
        </p>

        <ul className="mt-4 space-y-2 text-sm text-slate-400">
          <li>· Tự nhắc mỗi ~25 giây khi phiên đang bật</li>
          <li>· <span className="text-slate-200">Round mới</span> — gợi ý pistol / rotate discipline</li>
          <li>· <span className="text-slate-200">Nhắc tôi</span> — lấy tip tiếp theo ngay</li>
          <li>· Phím tắt (khi app focus): <span className="font-mono text-tactical-cyan">H</span> nhắc, <span className="font-mono text-tactical-cyan">R</span> round mới</li>
        </ul>

        <PanelActions
          session={session}
          overlayMode={overlayMode}
          busy={busy}
          onStart={onStart}
          onStop={onStop}
          onRequestTip={onRequestTip}
          onRoundStart={onRoundStart}
          onToggleOverlay={onToggleOverlay}
        />
      </div>

      <div className="space-y-4">
        <div className={`rounded-3xl border p-5 ${session.active ? 'border-tactical-cyan/30 bg-tactical-cyan/10' : 'border-tactical-line bg-tactical-panel/90'}`}>
          <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">Trạng thái</p>
          <p className="mt-2 text-xl font-black text-white">{session.active ? 'Đang bật' : 'Chưa bật'}</p>
          <p className="mt-2 text-sm text-slate-300">{session.message}</p>
          {session.active && (
            <div className="mt-3 flex flex-wrap gap-2 text-xs text-slate-400">
              <span className="rounded-full border border-white/10 px-2 py-1">Round {session.roundCount}</span>
              <span className="rounded-full border border-white/10 px-2 py-1">Tips {session.tipsShown}</span>
              <span className="rounded-full border border-white/10 px-2 py-1">Queue {session.queueSize}</span>
              {overlayMode && <span className="rounded-full border border-tactical-cyan/30 px-2 py-1 text-tactical-cyan">Overlay</span>}
            </div>
          )}
        </div>

        {alert && (
          <div className={`rounded-3xl border p-5 ${severityClass(alert.severity)}`}>
            <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">Gợi ý hiện tại</p>
            <p className="mt-2 text-lg font-black text-white">{alert.title}</p>
            <p className="mt-2 text-sm leading-6 text-slate-300">{alert.message}</p>
          </div>
        )}

        <HowToCard />
      </div>
    </section>
  )
}

const PanelActions: FC<{
  session: AssistantSessionState
  overlayMode: boolean
  busy: boolean
  onStart: () => void
  onStop: () => void
  onRequestTip: () => void
  onRoundStart: () => void
  onToggleOverlay: () => void
}> = ({ session, overlayMode, busy, onStart, onStop, onRequestTip, onRoundStart, onToggleOverlay }) => (
  <div className="mt-5 flex flex-wrap gap-2">
    {!session.active ? (
      <button
        type="button"
        disabled={busy}
        onClick={onStart}
        className="rounded-2xl bg-tactical-red px-5 py-3 text-sm font-black uppercase tracking-[0.16em] text-white shadow-lg shadow-tactical-red/20 disabled:opacity-50"
      >
        Bật Live Assistant
      </button>
    ) : (
      <>
        <button
          type="button"
          disabled={busy}
          onClick={onToggleOverlay}
          className="rounded-2xl border border-tactical-cyan/30 bg-tactical-cyan/10 px-4 py-3 text-sm font-black uppercase tracking-[0.14em] text-tactical-cyan disabled:opacity-50"
        >
          {overlayMode ? 'Thoát overlay' : 'Bật overlay (góc màn hình)'}
        </button>
        <button
          type="button"
          disabled={busy || !session.active}
          onClick={onRoundStart}
          className="rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3 text-sm font-black uppercase tracking-[0.14em] text-slate-200 disabled:opacity-50"
        >
          Round mới
        </button>
        <button
          type="button"
          disabled={busy || !session.active}
          onClick={onRequestTip}
          className="rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3 text-sm font-black uppercase tracking-[0.14em] text-slate-200 disabled:opacity-50"
        >
          Nhắc tôi
        </button>
        <button
          type="button"
          disabled={busy}
          onClick={onStop}
          className="rounded-2xl border border-tactical-red/30 bg-tactical-red/10 px-4 py-3 text-sm font-black uppercase tracking-[0.14em] text-tactical-red disabled:opacity-50"
        >
          Tắt phiên
        </button>
      </>
    )}
  </div>
)

const HowToCard: FC = () => (
  <div className="rounded-3xl border border-tactical-line bg-black/20 p-5 text-sm text-slate-400">
    <p className="font-bold text-slate-300">Cách dùng khi vào trận</p>
    <ol className="mt-2 list-decimal space-y-1 pl-4">
      <li>Fetch report ở tab Coach (hoặc dùng demo).</li>
      <li>Bật Live Assistant → Bật overlay.</li>
      <li>Alt+Tab nhanh bấm Round mới mỗi round, hoặc phím R khi cửa sổ focus.</li>
      <li>Chỉ đọc gợi ý — không auto điều khiển game (Safe mode).</li>
    </ol>
  </div>
)

export default LiveAssistantPanel
