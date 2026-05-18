import { useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import AssistantOverlay from './assistant-overlay'
import {
  fetchLiveReport,
  finishPracticeSession,
  generateDemoReport,
  getAPIStatus,
  getAssistantSession,
  getLastReport,
  getPracticeProgress,
  getPracticeSessions,
  markAssistantRoundStart,
  pollAssistantAutoTip,
  requestAssistantTip,
  resetPracticeProgress,
  riotIsLoggedIn,
  riotGetPlayerInfo,
  riotLogout,
  setPracticeProgress,
  startAssistantSession,
  stopAssistantSession,
} from './api'
import LiveAssistantPanel from './live-assistant-panel'
import MapPlannerPanel from './map-planner-panel'
import LoginScreen from './login-screen'
import CoachBot from './coach-bot'
import { disableOverlayMode, enableOverlayMode } from './overlay-window'
import type {
  APIStatus,
  AnalysisReport,
  AssistantSessionState,
  BreakdownRow,
  LiveAnalysisResult,
  PracticeProgressState,
  PracticeSessionInput,
  PracticeSessionState,
  PracticeTask,
  Recommendation,
  RiotPlayerInfo,
} from './types'

type Tab = 'coach' | 'assistant' | 'maps'

const tabs: Array<{ id: Tab; label: string }> = [
  { id: 'coach', label: 'Coach' },
  { id: 'assistant', label: 'Live' },
  { id: 'maps', label: 'Maps' },
]

const applyTipState = (result: { hasTip: boolean; state: AssistantSessionState }) => result.state

export const App = () => {
  const [activeTab, setActiveTab] = useState<Tab>('coach')
  const [report, setReport] = useState<AnalysisReport | null>(null)
  const [apiStatus, setApiStatus] = useState<APIStatus | null>(null)
  const [practiceProgress, setPracticeProgressState] = useState<PracticeProgressState>({ items: {}, updatedAt: '' })
  const [practiceSessions, setPracticeSessions] = useState<PracticeSessionState>({ sessions: [], updatedAt: '' })
  const [liveResult, setLiveResult] = useState<LiveAnalysisResult | null>(null)
  const [coachMessage, setCoachMessage] = useState('')
  const [coachError, setCoachError] = useState('')
  const [isFetching, setIsFetching] = useState(false)
  const [assistantSession, setAssistantSession] = useState<AssistantSessionState | null>(null)
  const [overlayMode, setOverlayMode] = useState(false)
  const [assistantBusy, setAssistantBusy] = useState(false)

  // Auth state
  const [authChecked, setAuthChecked] = useState(false)
  const [playerInfo, setPlayerInfo] = useState<RiotPlayerInfo | null>(null)

  // Kiểm tra auth state khi app khởi động
  useEffect(() => {
    const checkAuth = async () => {
      const loggedIn = await riotIsLoggedIn()
      if (loggedIn) {
        const info = await riotGetPlayerInfo()
        setPlayerInfo(info)
      }
      setAuthChecked(true)
    }
    void checkAuth()
  }, [])

  useEffect(() => {
    const load = async () => {
      const [lastReport, nextPracticeProgress, nextPracticeSessions, nextAPIStatus] = await Promise.all([
        getLastReport(),
        getPracticeProgress(),
        getPracticeSessions(),
        getAPIStatus(),
      ])

      if (lastReport.hasReport) {
        setReport(lastReport.result.report)
        setLiveResult(lastReport.result)
      } else {
        setReport(await generateDemoReport())
      }

      setPracticeProgressState(nextPracticeProgress)
      setPracticeSessions(nextPracticeSessions)
      setApiStatus(nextAPIStatus)
      setAssistantSession(await getAssistantSession())
    }

    void load()
  }, [])

  useEffect(() => {
    if (!assistantSession?.active) {
      return
    }
    const timer = window.setInterval(() => {
      void pollAssistantAutoTip().then((result) => {
        if (result.hasTip) {
          setAssistantSession(result.state)
        }
      })
    }, 8000)
    return () => window.clearInterval(timer)
  }, [assistantSession?.active])

  useEffect(() => {
    if (!assistantSession?.active) {
      return
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
        return
      }
      const key = event.key.toLowerCase()
      if (key === 'h') {
        event.preventDefault()
        void requestAssistantTip().then((result) => setAssistantSession(applyTipState(result)))
      }
      if (key === 'r') {
        event.preventDefault()
        void markAssistantRoundStart().then((result) => setAssistantSession(applyTipState(result)))
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [assistantSession?.active])

  const handleStartAssistant = async () => {
    setAssistantBusy(true)
    try {
      setAssistantSession(await startAssistantSession())
      setActiveTab('assistant')
    } finally {
      setAssistantBusy(false)
    }
  }

  const handleStopAssistant = async () => {
    setAssistantBusy(true)
    try {
      if (overlayMode) {
        await disableOverlayMode()
        setOverlayMode(false)
      }
      setAssistantSession(await stopAssistantSession())
    } finally {
      setAssistantBusy(false)
    }
  }

  const handleToggleOverlay = async () => {
    if (overlayMode) {
      await disableOverlayMode()
      setOverlayMode(false)
      return
    }
    await enableOverlayMode()
    setOverlayMode(true)
  }

  const handleExitOverlay = async () => {
    await disableOverlayMode()
    setOverlayMode(false)
  }

  if (overlayMode && assistantSession?.active) {
    return (
      <AssistantOverlay
        session={assistantSession}
        onRequestTip={() => void requestAssistantTip().then((result) => setAssistantSession(applyTipState(result)))}
        onRoundStart={() => void markAssistantRoundStart().then((result) => setAssistantSession(applyTipState(result)))}
        onStop={() => void handleStopAssistant()}
        onExitOverlay={() => void handleExitOverlay()}
      />
    )
  }

  if (!authChecked) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-tactical-bg text-slate-300">
        Đang tải...
      </main>
    )
  }

  if (!playerInfo) {
    return <LoginScreen onLoginSuccess={(info) => setPlayerInfo(info)} />
  }

  if (!assistantSession) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-tactical-bg text-slate-300">
        Đang tải...
      </main>
    )
  }

  return (
    <main className="min-h-screen bg-tactical-bg text-slate-100">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-5 px-4 py-5 sm:px-6 lg:px-8">
        <header className="rounded-3xl border border-tactical-line bg-gradient-to-br from-tactical-panel via-[#151825] to-[#0b0d14] p-5 shadow-2xl shadow-black/30">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div>
              <p className="w-fit rounded-full border border-tactical-red/40 bg-tactical-red/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-tactical-red">Valorant Tactical Trainer</p>
              <h1 className="mt-3 text-3xl font-black tracking-tight sm:text-5xl">Coach cá nhân cho game thủ Valorant.</h1>
              <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-300">App phân tích match history Riot, chỉ ra điểm yếu cụ thể của bạn (map, agent, vai trò), gợi ý bài tập hôm nay và hỗ trợ Live Assistant + Map Planner trong lúc thi đấu.</p>
            </div>
            <div className="flex flex-col items-end gap-3">
              <div className="grid grid-cols-2 gap-2 text-center text-xs sm:min-w-[260px]">
                <MetricPill label="Mode" value="Lean" />
                <MetricPill label="Safe" value={apiStatus?.safeMode ? 'On' : 'Off'} />
              </div>
              {playerInfo && (
                <div className="flex items-center gap-3 rounded-2xl border border-tactical-line bg-black/20 px-4 py-2">
                  <div className="h-2 w-2 rounded-full bg-green-400 shadow-sm shadow-green-400/50" />
                  <span className="text-sm font-bold text-white">{playerInfo.gameName}<span className="text-tactical-red">#{playerInfo.tagLine}</span></span>
                  <button
                    type="button"
                    onClick={() => { void riotLogout().then(() => setPlayerInfo(null)) }}
                    className="text-xs text-slate-500 transition hover:text-slate-300"
                  >
                    Logout
                  </button>
                </div>
              )}
            </div>
          </div>
        </header>

        <nav className="flex w-fit gap-2 rounded-2xl border border-tactical-line bg-tactical-panel/70 p-2">
          {tabs.map(({ id, label }) => (
            <button
              key={id}
              type="button"
              onClick={() => setActiveTab(id)}
              className={`rounded-xl px-4 py-2 text-sm font-semibold transition ${activeTab === id ? 'bg-tactical-red text-white shadow-lg shadow-tactical-red/20' : 'text-slate-300 hover:bg-white/5 hover:text-white'}`}
            >
              {label}
            </button>
          ))}
        </nav>

        {activeTab === 'coach' && (
          <CoachPanel
            report={report}
            playerInfo={playerInfo}
            apiStatus={apiStatus}
            message={coachMessage}
            error={coachError}
            isFetching={isFetching}
            progress={practiceProgress}
            sessions={practiceSessions}
            onFetchLive={async () => {
              setIsFetching(true)
              setCoachError('')
              setCoachMessage('')
              try {
                const result = await fetchLiveReport()
                setReport(result.report)
                setLiveResult(result)
                setApiStatus(await getAPIStatus())
              } catch (err) {
                setCoachError(errorMessage(err))
              } finally {
                setIsFetching(false)
              }
            }}
            onSetPracticeProgress={async (itemID, done) => setPracticeProgressState(await setPracticeProgress(itemID, done))}
            onResetPracticeProgress={async () => setPracticeProgressState(await resetPracticeProgress())}
            onFinishPracticeSession={async (input) => setPracticeSessions(await finishPracticeSession(input))}
          />
        )}

        {activeTab === 'maps' && (
          <MapPlannerPanel suggestedMapName={report?.metrics.weakestMap} />
        )}

        {activeTab === 'assistant' && (
          <LiveAssistantPanel
            session={assistantSession}
            overlayMode={overlayMode}
            busy={assistantBusy}
            onStart={() => void handleStartAssistant()}
            onStop={() => void handleStopAssistant()}
            onRequestTip={() => void requestAssistantTip().then((result) => setAssistantSession(applyTipState(result)))}
            onRoundStart={() => void markAssistantRoundStart().then((result) => setAssistantSession(applyTipState(result)))}
            onToggleOverlay={() => void handleToggleOverlay()}
          />
        )}
      </section>

      {/* Floating bot AI — chỉ hiển thị khi user đã login (có playerInfo). */}
      <CoachBot hasReport={Boolean(report && liveResult)} />
    </main>
  )
}

const CoachPanel = ({
  report,
  playerInfo,
  apiStatus,
  message,
  error,
  isFetching,
  progress,
  sessions,
  onFetchLive,
  onSetPracticeProgress,
  onResetPracticeProgress,
  onFinishPracticeSession,
}: {
  report: AnalysisReport | null
  playerInfo: RiotPlayerInfo | null
  apiStatus: APIStatus | null
  message: string
  error: string
  isFetching: boolean
  progress: PracticeProgressState
  sessions: PracticeSessionState
  onFetchLive: () => Promise<void>
  onSetPracticeProgress: (itemID: string, done: boolean) => Promise<void>
  onResetPracticeProgress: () => Promise<void>
  onFinishPracticeSession: (input: PracticeSessionInput) => Promise<void>
}) => {
  if (!report) {
    return <EmptyState title="Đang tải Coach" />
  }

  const mainFinding = report.findings[0]
  const mainRecommendation = report.recommendations[0]
  const todayTask = report.practicePlan[0]

  // Tên player ưu tiên lấy từ session login (Riot Account-V1) — tên này luôn
  // chính xác kể cả khi report chưa fetch hoặc đang dùng demo. Fallback về
  // report.player nếu chưa login (vd offline / debug).
  const displayName = playerInfo?.gameName ?? report.player.name
  const displayTag = playerInfo?.tagLine ?? report.player.tagline

  return (
    <section className="grid gap-5 lg:grid-cols-[0.95fr_1.05fr]">
      <div className="space-y-5">
        <Panel>
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">Player</p>
              <h2 className="mt-1 text-3xl font-black text-white">{displayName}#{displayTag}</h2>
              {message && <p className="mt-2 text-sm text-slate-400">{message}</p>}
              {error && <p className="mt-2 text-sm font-semibold text-tactical-red">{error}</p>}
            </div>
            <button
              type="button"
              disabled={isFetching || !apiStatus?.canFetchPersonalData}
              onClick={() => void onFetchLive()}
              className="rounded-2xl bg-tactical-red px-5 py-3 text-sm font-black uppercase tracking-[0.16em] text-white shadow-lg shadow-tactical-red/20 transition disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400 disabled:shadow-none"
            >
              {isFetching ? 'Đang tải...' : 'Fetch report'}
            </button>
          </div>
        </Panel>

        <MetricGrid report={report} />
        <CompactBreakdown title="Map cần chú ý" rows={report.mapBreakdown.slice(0, 3)} value="winRate" />
        <CompactBreakdown title="Agent đang dùng" rows={report.agentBreakdown.slice(0, 3)} value="kd" />
      </div>

      <div className="space-y-5">
        {mainFinding && <MainFindingCard finding={mainFinding} />}
        {mainRecommendation && <RecommendationCard recommendation={mainRecommendation} />}
        {todayTask && (
          <TodayPracticeCard
            task={todayTask}
            progress={progress}
            sessions={sessions}
            onSetProgress={onSetPracticeProgress}
            onResetProgress={onResetPracticeProgress}
            onFinishSession={onFinishPracticeSession}
          />
        )}
      </div>
    </section>
  )
}

const MetricGrid = ({ report }: { report: AnalysisReport }) => {
  const metrics = [
    { label: 'Matches', value: report.metrics.matches.toString() },
    { label: 'Win Rate', value: `${(report.metrics.winRate * 100).toFixed(0)}%` },
    { label: 'K/D', value: report.metrics.kd.toFixed(2) },
    { label: 'HS%', value: `${report.metrics.headshotPercent.toFixed(1)}%` },
    { label: 'First Death', value: `${(report.metrics.firstDeathRate * 100).toFixed(0)}%` },
    { label: 'Map yếu', value: report.metrics.weakestMap || 'N/A' },
  ]

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
      {metrics.map((metric) => <MetricCard key={metric.label} label={metric.label} value={metric.value} />)}
    </div>
  )
}

const MainFindingCard = ({ finding }: { finding: AnalysisReport['findings'][number] }) => (
  <Panel>
    <p className="text-xs font-bold uppercase tracking-[0.2em] text-tactical-red">Vấn đề chính</p>
    <h2 className="mt-2 text-2xl font-black text-white">{finding.title}</h2>
    <p className="mt-3 text-sm leading-6 text-slate-300">{finding.detail}</p>
    <div className="mt-4 grid gap-2 sm:grid-cols-2">
      <MetricPill label="Severity" value={finding.severity} />
      <MetricPill label="Confidence" value={finding.confidence} />
    </div>
  </Panel>
)

const RecommendationCard = ({ recommendation }: { recommendation: Recommendation }) => (
  <Panel>
    <p className="text-xs font-bold uppercase tracking-[0.2em] text-tactical-cyan">Cần làm</p>
    <h2 className="mt-2 text-2xl font-black text-white">{recommendation.title}</h2>
    <p className="mt-3 text-sm leading-6 text-slate-300">{recommendation.drill}</p>
    <p className="mt-4 w-fit rounded-full bg-tactical-cyan/10 px-3 py-2 text-xs font-black text-tactical-cyan">{recommendation.cadence}</p>
  </Panel>
)

const TodayPracticeCard = ({
  task,
  progress,
  sessions,
  onSetProgress,
  onResetProgress,
  onFinishSession,
}: {
  task: PracticeTask
  progress: PracticeProgressState
  sessions: PracticeSessionState
  onSetProgress: (itemID: string, done: boolean) => Promise<void>
  onResetProgress: () => Promise<void>
  onFinishSession: (input: PracticeSessionInput) => Promise<void>
}) => {
  const [startedAt, setStartedAt] = useState('')
  const [elapsedSeconds, setElapsedSeconds] = useState(0)

  useEffect(() => {
    if (!startedAt) {
      return
    }
    const timer = window.setInterval(() => setElapsedSeconds(Math.max(0, Math.floor((Date.now() - Date.parse(startedAt)) / 1000))), 1000)
    return () => window.clearInterval(timer)
  }, [startedAt])

  const ids = task.checklist.map((item, index) => practiceItemID(task, item, index))
  const done = ids.filter((id) => progress.items[id]).length

  return (
    <Panel>
      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">Bài tập hôm nay</p>
          <h2 className="mt-2 text-2xl font-black text-white">{task.focus}</h2>
          <p className="mt-2 text-sm text-slate-400">{task.map ? `${task.map} ` : ''}{task.agent ? `/ ${task.agent} ` : ''}/ {task.duration}</p>
        </div>
        <span className="rounded-full bg-tactical-cyan/10 px-3 py-2 text-xs font-black text-tactical-cyan">{done}/{ids.length}</span>
      </div>

      <div className="mt-4 space-y-2">
        {task.checklist.map((item, index) => {
          const itemID = practiceItemID(task, item, index)
          const checked = Boolean(progress.items[itemID])
          return (
            <label key={itemID} className={`flex cursor-pointer items-start gap-3 rounded-xl border px-3 py-2 text-sm ${checked ? 'border-tactical-cyan/30 bg-tactical-cyan/10 text-tactical-cyan' : 'border-white/10 bg-white/[0.03] text-slate-300'}`}>
              <input type="checkbox" checked={checked} onChange={(event) => void onSetProgress(itemID, event.target.checked)} className="mt-1 h-4 w-4 accent-tactical-cyan" />
              <span className={checked ? 'line-through decoration-tactical-cyan/70' : ''}>{item}</span>
            </label>
          )
        })}
      </div>

      <div className="mt-4 grid gap-2 sm:grid-cols-3">
        <button type="button" onClick={() => setStartedAt(startedAt ? '' : new Date().toISOString())} className="rounded-xl border border-tactical-cyan/30 bg-tactical-cyan/10 px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-tactical-cyan">
          {startedAt ? `Stop ${formatDuration(elapsedSeconds)}` : 'Start timer'}
        </button>
        <button
          type="button"
          disabled={!startedAt}
          onClick={() => {
            void onFinishSession({ taskId: taskID(task), focus: task.focus, map: task.map ?? '', agent: task.agent ?? '', durationSeconds: elapsedSeconds, startedAt }).then(() => {
              setStartedAt('')
              setElapsedSeconds(0)
            })
          }}
          className="rounded-xl border border-tactical-red/30 bg-tactical-red/10 px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-tactical-red disabled:cursor-not-allowed disabled:opacity-50"
        >
          Finish
        </button>
        <button type="button" onClick={() => void onResetProgress()} className="rounded-xl border border-white/10 bg-white/[0.03] px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-slate-300">Reset</button>
      </div>

      {sessions.sessions.length > 0 && <p className="mt-3 text-xs text-slate-500">Session gần nhất: {formatDuration(sessions.sessions[0].durationSeconds)} / {formatDateTime(sessions.sessions[0].finishedAt)}</p>}
    </Panel>
  )
}

const CompactBreakdown = ({ title, rows, value }: { title: string; rows: BreakdownRow[]; value: 'winRate' | 'kd' }) => (
  <Panel>
    <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">{title}</p>
    <div className="mt-3 space-y-2">
      {rows.length === 0 ? <p className="text-sm text-slate-400">Chưa đủ dữ liệu.</p> : rows.map((row) => {
        const score = value === 'winRate' ? row.winRate : Math.min(1, row.kd / 1.6)
        const label = value === 'winRate' ? `${(row.winRate * 100).toFixed(0)}% WR` : `${row.kd.toFixed(2)} KD`
        return (
          <div key={`${title}-${row.name}`}>
            <div className="mb-1 flex items-center justify-between text-sm">
              <span className="font-bold text-white">{row.name}</span>
              <span className="text-xs text-slate-400">{label}</span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-white/10">
              <div className="h-full rounded-full bg-tactical-cyan" style={{ width: `${Math.max(5, score * 100)}%` }} />
            </div>
          </div>
        )
      })}
    </div>
  </Panel>
)

const Panel = ({ children }: { children: ReactNode }) => <section className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">{children}</section>

const MetricPill = ({ label, value }: { label: string; value: string }) => (
  <div className="rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3">
    <div className="text-[11px] uppercase tracking-[0.2em] text-slate-500">{label}</div>
    <div className="mt-1 font-bold text-white">{value}</div>
  </div>
)

const MetricCard = ({ label, value }: { label: string; value: string }) => (
  <div className="rounded-2xl border border-tactical-line bg-black/20 p-4">
    <p className="text-xs uppercase tracking-[0.18em] text-slate-500">{label}</p>
    <p className="mt-2 text-2xl font-black text-white">{value}</p>
  </div>
)

const EmptyState = ({ title }: { title: string }) => <div className="rounded-3xl border border-tactical-line bg-tactical-panel p-6 text-slate-300">{title}</div>

const practiceItemID = (task: PracticeTask, item: string, index: number) => [task.day, task.focus, task.map ?? 'all-map', task.agent ?? 'all-agent', index, item].map(slugPart).join('__')

const taskID = (task: PracticeTask) => [task.day, task.focus, task.map ?? 'all-map', task.agent ?? 'all-agent'].map(slugPart).join('__')

const slugPart = (value: string | number) => String(value).trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '') || 'empty'

const formatDuration = (seconds: number) => {
  const safeSeconds = Math.max(0, Math.floor(seconds))
  const minutes = Math.floor(safeSeconds / 60)
  const rest = safeSeconds % 60
  return `${String(minutes).padStart(2, '0')}:${String(rest).padStart(2, '0')}`
}

const formatDateTime = (value: string) => {
  if (!value) return 'N/A'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
}

const errorMessage = (err: unknown) => {
  if (err instanceof Error && err.message) return err.message
  if (typeof err === 'string' && err.trim()) return err
  if (err && typeof err === 'object') {
    const maybeMessage = 'message' in err ? err.message : undefined
    if (typeof maybeMessage === 'string' && maybeMessage.trim()) return maybeMessage
  }
  return 'err fetch report'
}
