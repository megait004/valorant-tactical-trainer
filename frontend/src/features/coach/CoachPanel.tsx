// Coach panel — tab "Coach" trong shell App. Hiển thị:
//   • Player header + nút Fetch live report (Henrik).
//   • Metric grid 6 ô (Matches, WR, K/D, HS%, First Death, Map yếu).
//   • Top 3 maps + agents (breakdown).
//   • Main finding + recommendation chính.
//   • Today practice card (checklist + timer + finish session).
//
// Toàn bộ data về từ props — không tự fetch để App giữ source of truth.

import { useEffect, useState } from 'react'
import { EmptyState, MetricCard, MetricPill, Panel } from '../../components/Panel'
import { formatDateTime, formatDuration, practiceItemID, taskID } from '../../lib/format'
import type {
  APIStatus,
  AnalysisReport,
  BreakdownRow,
  PracticeProgressState,
  PracticeSessionInput,
  PracticeSessionState,
  PracticeTask,
  Recommendation,
  RiotPlayerInfo,
} from '../../types'

export interface CoachPanelProps {
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
}

export const CoachPanel = ({
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
}: CoachPanelProps) => {
  if (!report) {
    return <EmptyState title="Đang tải Coach" />
  }

  const mainFinding = report.findings[0]
  const mainRecommendation = report.recommendations[0]
  const todayTask = report.practicePlan[0]

  // Tên player ưu tiên lấy từ session login (Riot Account-V1) — luôn chính xác
  // kể cả khi report chưa fetch hoặc đang dùng demo.
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
      {metrics.map((metric) => (
        <MetricCard key={metric.label} label={metric.label} value={metric.value} />
      ))}
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

interface TodayPracticeCardProps {
  task: PracticeTask
  progress: PracticeProgressState
  sessions: PracticeSessionState
  onSetProgress: (itemID: string, done: boolean) => Promise<void>
  onResetProgress: () => Promise<void>
  onFinishSession: (input: PracticeSessionInput) => Promise<void>
}

const TodayPracticeCard = ({
  task,
  progress,
  sessions,
  onSetProgress,
  onResetProgress,
  onFinishSession,
}: TodayPracticeCardProps) => {
  const [startedAt, setStartedAt] = useState('')
  const [elapsedSeconds, setElapsedSeconds] = useState(0)

  useEffect(() => {
    if (!startedAt) {
      return
    }
    const timer = window.setInterval(
      () => setElapsedSeconds(Math.max(0, Math.floor((Date.now() - Date.parse(startedAt)) / 1000))),
      1000,
    )
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
              <input
                type="checkbox"
                checked={checked}
                onChange={(event) => void onSetProgress(itemID, event.target.checked)}
                className="mt-1 h-4 w-4 accent-tactical-cyan"
              />
              <span className={checked ? 'line-through decoration-tactical-cyan/70' : ''}>{item}</span>
            </label>
          )
        })}
      </div>

      <div className="mt-4 grid gap-2 sm:grid-cols-3">
        <button
          type="button"
          onClick={() => setStartedAt(startedAt ? '' : new Date().toISOString())}
          className="rounded-xl border border-tactical-cyan/30 bg-tactical-cyan/10 px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-tactical-cyan"
        >
          {startedAt ? `Stop ${formatDuration(elapsedSeconds)}` : 'Start timer'}
        </button>
        <button
          type="button"
          disabled={!startedAt}
          onClick={() => {
            void onFinishSession({
              taskId: taskID(task),
              focus: task.focus,
              map: task.map ?? '',
              agent: task.agent ?? '',
              durationSeconds: elapsedSeconds,
              startedAt,
            }).then(() => {
              setStartedAt('')
              setElapsedSeconds(0)
            })
          }}
          className="rounded-xl border border-tactical-red/30 bg-tactical-red/10 px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-tactical-red disabled:cursor-not-allowed disabled:opacity-50"
        >
          Finish
        </button>
        <button
          type="button"
          onClick={() => void onResetProgress()}
          className="rounded-xl border border-white/10 bg-white/[0.03] px-3 py-2 text-xs font-black uppercase tracking-[0.14em] text-slate-300"
        >
          Reset
        </button>
      </div>

      {sessions.sessions.length > 0 && (
        <p className="mt-3 text-xs text-slate-500">
          Session gần nhất: {formatDuration(sessions.sessions[0].durationSeconds)} / {formatDateTime(sessions.sessions[0].finishedAt)}
        </p>
      )}
    </Panel>
  )
}

interface CompactBreakdownProps {
  title: string
  rows: BreakdownRow[]
  value: 'winRate' | 'kd'
}

const CompactBreakdown = ({ title, rows, value }: CompactBreakdownProps) => (
  <Panel>
    <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500">{title}</p>
    <div className="mt-3 space-y-2">
      {rows.length === 0 ? (
        <p className="text-sm text-slate-400">Chưa đủ dữ liệu.</p>
      ) : (
        rows.map((row) => {
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
        })
      )}
    </div>
  </Panel>
)
