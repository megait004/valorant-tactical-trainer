import { useEffect, useRef, useState } from 'react'
import { chatGetState, chatIsAvailable, chatReset, chatSendMessage } from './api'
import type { ChatState } from './types'

type CoachBotProps = {
  /**
   * Có report đã fetch chưa — bot chỉ ý nghĩa khi có context. Nếu chưa có report
   * thì icon vẫn hiển thị nhưng panel sẽ hiện hint "fetch report trước".
   */
  hasReport: boolean
}

const SUGGESTED_QUESTIONS = [
  'Vấn đề lớn nhất hiện tại của tôi là gì?',
  'Tập gì 30 phút hôm nay để khắc phục?',
  'Map nào tôi đang yếu nhất và lý do?',
  'Agent nào hợp với style của tôi?',
]

const CoachBot = ({ hasReport }: CoachBotProps) => {
  const [open, setOpen] = useState(false)
  const [available, setAvailable] = useState<boolean | null>(null)
  const [state, setState] = useState<ChatState>({ available: false, history: [] })
  const [input, setInput] = useState('')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState('')
  const scrollRef = useRef<HTMLDivElement | null>(null)
  const inputRef = useRef<HTMLTextAreaElement | null>(null)

  useEffect(() => {
    void chatIsAvailable().then(setAvailable)
  }, [])

  useEffect(() => {
    if (!open) return
    void chatGetState().then(setState)
    setTimeout(() => inputRef.current?.focus(), 50)
  }, [open])

  useEffect(() => {
    if (!scrollRef.current) return
    scrollRef.current.scrollTop = scrollRef.current.scrollHeight
  }, [state.history.length, sending])

  const handleSend = async (text?: string) => {
    const message = (text ?? input).trim()
    if (!message || sending) return
    setSending(true)
    setError('')
    setInput('')
    try {
      const next = await chatSendMessage(message)
      setState(next)
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setSending(false)
    }
  }

  const handleReset = async () => {
    setError('')
    setState(await chatReset())
  }

  // Floating button luôn hiển thị (khi available). Click toggle panel.
  if (available === false) {
    // Coach không cấu hình LLM_API_KEY — ẩn icon hoàn toàn để UI gọn.
    return null
  }

  return (
    <>
      <button
        type="button"
        aria-label="Mở Coach AI"
        onClick={() => setOpen((v) => !v)}
        className="fixed bottom-6 right-6 z-40 flex h-14 w-14 items-center justify-center rounded-full bg-tactical-red text-white shadow-2xl shadow-tactical-red/30 transition hover:scale-105 active:scale-95"
      >
        <BotIcon className="h-7 w-7" />
        {state.history.length > 0 && !open && (
          <span className="absolute -right-1 -top-1 flex h-5 w-5 items-center justify-center rounded-full bg-tactical-cyan text-[10px] font-black text-tactical-bg">
            {Math.min(state.history.length, 99)}
          </span>
        )}
      </button>

      {open && (
        <div className="fixed bottom-24 right-6 z-40 flex h-[560px] w-[380px] max-w-[calc(100vw-2rem)] flex-col rounded-3xl border border-tactical-line bg-tactical-panel shadow-2xl shadow-black/40 sm:w-[420px]">
          <header className="flex items-center justify-between gap-2 rounded-t-3xl border-b border-tactical-line bg-gradient-to-br from-[#1a1d2c] to-[#0e1018] px-4 py-3">
            <div className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-tactical-red/20 text-tactical-red">
                <BotIcon className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-bold text-white">Coach AI</p>
                <p className="text-[11px] text-slate-500">
                  {hasReport ? 'Đang dùng report gần nhất làm context' : 'Fetch report trước để có context tốt hơn'}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-1">
              {state.history.length > 0 && (
                <button
                  type="button"
                  onClick={() => void handleReset()}
                  className="rounded-lg px-2 py-1 text-[11px] text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
                  title="Xoá hội thoại"
                >
                  Reset
                </button>
              )}
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="rounded-lg px-2 py-1 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
                aria-label="Đóng"
              >
                ✕
              </button>
            </div>
          </header>

          <div ref={scrollRef} className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
            {state.history.length === 0 && (
              <div className="space-y-3">
                <div className="rounded-2xl border border-tactical-line bg-black/20 p-3 text-sm leading-6 text-slate-300">
                  Hỏi mình bất cứ gì về stat, map, agent, drill luyện tập của bạn. Mình sẽ trả lời dựa trên report gần nhất.
                </div>
                <div className="space-y-2">
                  {SUGGESTED_QUESTIONS.map((q) => (
                    <button
                      key={q}
                      type="button"
                      onClick={() => void handleSend(q)}
                      disabled={sending}
                      className="w-full rounded-xl border border-white/10 bg-white/[0.03] px-3 py-2 text-left text-xs text-slate-300 transition hover:border-tactical-cyan/40 hover:bg-tactical-cyan/5 disabled:opacity-50"
                    >
                      {q}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {state.history.map((m, idx) => (
              <MessageBubble key={`${m.createdAt}-${idx}`} role={m.role} content={m.content} />
            ))}

            {sending && (
              <MessageBubble role="assistant" content="Đang nghĩ..." pending />
            )}

            {error && (
              <p className="rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs leading-5 text-red-300">
                {error}
              </p>
            )}
          </div>

          <form
            onSubmit={(e) => {
              e.preventDefault()
              void handleSend()
            }}
            className="border-t border-tactical-line bg-black/20 px-3 py-3"
          >
            <div className="flex items-end gap-2">
              <textarea
                ref={inputRef}
                rows={2}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault()
                    void handleSend()
                  }
                }}
                placeholder="Hỏi về drill, map yếu, agent..."
                disabled={sending}
                className="flex-1 resize-none rounded-xl border border-tactical-line bg-tactical-bg px-3 py-2 text-sm text-slate-100 placeholder-slate-600 outline-none transition focus:border-tactical-cyan/40 focus:ring-1 focus:ring-tactical-cyan/30 disabled:opacity-50"
              />
              <button
                type="submit"
                disabled={sending || !input.trim()}
                className="rounded-xl bg-tactical-red px-4 py-2 text-xs font-black uppercase tracking-[0.16em] text-white shadow-lg shadow-tactical-red/20 transition hover:bg-tactical-red/80 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400 disabled:shadow-none"
              >
                Gửi
              </button>
            </div>
            <p className="mt-2 text-[10px] text-slate-600">Enter để gửi · Shift+Enter xuống dòng</p>
          </form>
        </div>
      )}
    </>
  )
}

const MessageBubble = ({
  role,
  content,
  pending,
}: {
  role: 'user' | 'assistant'
  content: string
  pending?: boolean
}) => {
  const isUser = role === 'user'
  return (
    <div className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}>
      <div
        className={`max-w-[85%] whitespace-pre-wrap rounded-2xl px-3 py-2 text-sm leading-6 ${
          isUser
            ? 'bg-tactical-red/90 text-white'
            : pending
              ? 'border border-tactical-line bg-black/30 text-slate-500 italic'
              : 'border border-tactical-line bg-black/30 text-slate-200'
        }`}
      >
        {content}
      </div>
    </div>
  )
}

const BotIcon = ({ className = '' }: { className?: string }) => (
  <svg
    className={className}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    aria-hidden="true"
  >
    <rect x="4" y="7" width="16" height="12" rx="3" />
    <path d="M12 3v4" />
    <circle cx="9" cy="13" r="1" fill="currentColor" />
    <circle cx="15" cy="13" r="1" fill="currentColor" />
    <path d="M9 17h6" />
  </svg>
)

export default CoachBot
