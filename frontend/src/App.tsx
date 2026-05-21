// App shell — orchestrate top-level state (auth, report, practice, assistant)
// và route giữa 3 tab: Coach / Live Assistant / Maps. Mọi UI con nằm trong
// `features/<feature>/` và shared trong `components/`.
//
// Lifecycle quan trọng:
//   1. Khởi động → check auth (riotIsLoggedIn) → show LoginScreen nếu chưa.
//   2. Login OK → load lastReport / practiceProgress / sessions / apiStatus
//      / assistantSession song song.
//   3. Live Assistant active → poll auto tip mỗi 8s + bind hotkey H/R.
//   4. Overlay mode → render full-screen AssistantOverlay thay shell.

import { useEffect, useState } from 'react'
import AssistantOverlay from './features/assistant/AssistantOverlay'
import LiveAssistantPanel from './features/assistant/LiveAssistantPanel'
import { disableOverlayMode, enableOverlayMode } from './features/assistant/overlay-window'
import LoginScreen from './features/auth/LoginScreen'
import CoachBot from './features/coach/CoachBot'
import { CoachPanel } from './features/coach/CoachPanel'
import MapPlannerPanel from './features/maps/MapPlannerPanel'
import AgentBrowserPanel from './features/agents/AgentBrowserPanel'
import { AppLogo } from './components/AppLogo'
import { MetricPill } from './components/Panel'
import { errorMessage } from './lib/format'
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
  riotGetPlayerInfo,
  riotIsLoggedIn,
  riotLogout,
  setPracticeProgress,
  startAssistantSession,
  stopAssistantSession,
} from './api'
import type {
  APIStatus,
  AnalysisReport,
  AssistantSessionState,
  LiveAnalysisResult,
  PracticeProgressState,
  PracticeSessionState,
  RiotPlayerInfo,
} from './types'

type Tab = 'coach' | 'assistant' | 'maps' | 'agents'

const tabs: Array<{ id: Tab; label: string }> = [
  { id: 'coach', label: 'Coach' },
  { id: 'assistant', label: 'Live' },
  { id: 'maps', label: 'Maps' },
  { id: 'agents', label: 'Agents' },
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

  const [authChecked, setAuthChecked] = useState(false)
  const [playerInfo, setPlayerInfo] = useState<RiotPlayerInfo | null>(null)

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

  const handleFetchLive = async () => {
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
    return <Loading />
  }

  if (!playerInfo) {
    return (
      <LoginScreen
        onLoginSuccess={(info) => {
          setPlayerInfo(info)
          // Login ghi lại consent + Riot ID vào settings.json (authSink) nhưng
          // apiStatus đã load lúc mount — phải refresh để bật Fetch report.
          void getAPIStatus().then(setApiStatus)
        }}
      />
    )
  }

  if (!assistantSession) {
    return <Loading />
  }

  return (
    <main className="min-h-screen bg-tactical-bg text-slate-100">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-5 px-4 py-5 sm:px-6 lg:px-8">
        <AppHeader playerInfo={playerInfo} apiStatus={apiStatus} onLogout={() => { void riotLogout().then(() => setPlayerInfo(null)) }} />

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
            onFetchLive={handleFetchLive}
            onSetPracticeProgress={async (itemID, done) => setPracticeProgressState(await setPracticeProgress(itemID, done))}
            onResetPracticeProgress={async () => setPracticeProgressState(await resetPracticeProgress())}
            onFinishPracticeSession={async (input) => setPracticeSessions(await finishPracticeSession(input))}
          />
        )}

        {activeTab === 'maps' && <MapPlannerPanel suggestedMapName={report?.metrics.weakestMap} />}

        {activeTab === 'agents' && <AgentBrowserPanel />}

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

const Loading = () => (
  <main className="flex min-h-screen flex-col items-center justify-center gap-4 bg-tactical-bg text-slate-300">
    <AppLogo size={72} />
    <p>Đang tải...</p>
  </main>
)

interface AppHeaderProps {
  playerInfo: RiotPlayerInfo
  apiStatus: APIStatus | null
  onLogout: () => void
}

const AppHeader = ({ playerInfo, apiStatus, onLogout }: AppHeaderProps) => (
  <header className="rounded-3xl border border-tactical-line bg-gradient-to-br from-tactical-panel via-[#151825] to-[#0b0d14] p-5 shadow-2xl shadow-black/30">
    <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
      <div className="flex gap-4 sm:items-center">
        <AppLogo size={64} className="hidden shrink-0 sm:block" />
        <div>
          <p className="w-fit rounded-full border border-tactical-red/40 bg-tactical-red/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-tactical-red">
            Valorant Tactical Trainer
          </p>
          <h1 className="mt-3 text-3xl font-black tracking-tight sm:text-5xl">Coach cá nhân cho game thủ Valorant.</h1>
          <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-300">
            App phân tích match history Riot, chỉ ra điểm yếu cụ thể của bạn (map, agent, vai trò), gợi ý bài tập hôm nay và hỗ trợ Live Assistant + Map Planner trong lúc thi đấu.
          </p>
        </div>
      </div>
      <div className="flex flex-col items-end gap-3">
        <div className="grid grid-cols-2 gap-2 text-center text-xs sm:min-w-[260px]">
          <MetricPill label="Mode" value="Lean" />
          <MetricPill label="Safe" value={apiStatus?.safeMode ? 'On' : 'Off'} />
        </div>
        <div className="flex items-center gap-3 rounded-2xl border border-tactical-line bg-black/20 px-4 py-2">
          <div className="h-2 w-2 rounded-full bg-green-400 shadow-sm shadow-green-400/50" />
          <span className="text-sm font-bold text-white">
            {playerInfo.gameName}
            <span className="text-tactical-red">#{playerInfo.tagLine}</span>
          </span>
          <button type="button" onClick={onLogout} className="text-xs text-slate-500 transition hover:text-slate-300">
            Logout
          </button>
        </div>
      </div>
    </div>
  </header>
)
