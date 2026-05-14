import { ReportPanel } from './ReportPanel';
import type { SetupPanelProps } from './types';

export const SetupPanel = ({
  appInfo,
  apiKey,
  canLookup,
  currentPlayer,
  lookupLoading,
  lookupResult,
  matchLoading,
  name,
  rank,
  rankLoading,
  region,
  report,
  reportLoading,
  resetLoading,
  tag,
  t,
  onApiKeyChange,
  onCheckCore,
  onGenerateReport,
  onLookupPlayer,
  onNameChange,
  onRefreshMatches,
  onRefreshRank,
  onRegionChange,
  onResetAllData,
  onTagChange,
}: SetupPanelProps) => (
  <section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-7 shadow-2xl shadow-black/30 backdrop-blur">
    <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">{t.consentGate}</p>
    <h2 className="mt-4 max-w-3xl text-3xl font-bold leading-tight text-white md:text-5xl">{t.consentTitle}</h2>
    <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300">{t.consentDescription}</p>

    <div className="mt-8 grid gap-4 md:grid-cols-2">
      <label className="space-y-2">
        <span className="text-sm font-semibold text-slate-300">{t.riotName}</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
          onChange={(event) => onNameChange(event.target.value)}
          placeholder={t.playerNamePlaceholder}
          type="text"
          value={name}
        />
      </label>
      <label className="space-y-2">
        <span className="text-sm font-semibold text-slate-300">Tag</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
          onChange={(event) => onTagChange(event.target.value)}
          placeholder="VN2"
          type="text"
          value={tag}
        />
      </label>
      <label className="space-y-2">
        <span className="text-sm font-semibold text-slate-300">{t.regionFallback}</span>
        <select
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition focus:border-tactical-red"
          onChange={(event) => onRegionChange(event.target.value)}
          value={region}
        >
          <option value="ap">AP</option>
          <option value="eu">EU</option>
          <option value="na">NA</option>
          <option value="kr">KR</option>
          <option value="latam">LATAM</option>
          <option value="br">BR</option>
        </select>
      </label>
      <label className="space-y-2">
        <span className="text-sm font-semibold text-slate-300">{t.apiKeyOptional}</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
          onChange={(event) => onApiKeyChange(event.target.value)}
          placeholder={t.apiKeySetupPlaceholder}
          type="password"
          value={apiKey}
        />
      </label>
    </div>

    <div className="mt-6 flex flex-wrap gap-3">
      <button
        className="rounded-full bg-tactical-red px-5 py-3 text-sm font-bold text-white shadow-lg shadow-red-950/40 transition hover:-translate-y-0.5 hover:bg-red-400 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400 disabled:shadow-none"
        disabled={!canLookup}
        onClick={onLookupPlayer}
        type="button"
      >
        {lookupLoading ? t.lookupLoading : t.lookupPlayer}
      </button>
      <button
        className="rounded-full bg-tactical-red px-5 py-3 text-sm font-bold text-white shadow-lg shadow-red-950/40 transition hover:-translate-y-0.5 hover:bg-red-400"
        onClick={onCheckCore}
        type="button"
      >
        {t.checkCore}
      </button>
      <button
        className="rounded-full border border-tactical-cyan/40 px-5 py-3 text-sm font-bold text-tactical-cyan transition hover:-translate-y-0.5 hover:bg-tactical-cyan/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={!currentPlayer || matchLoading}
        onClick={onRefreshMatches}
        type="button"
      >
        {matchLoading ? t.refreshing : t.refreshMatches}
      </button>
      <button
        className="rounded-full border border-emerald-300/40 px-5 py-3 text-sm font-bold text-emerald-200 transition hover:-translate-y-0.5 hover:bg-emerald-300/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={!currentPlayer || rankLoading}
        onClick={onRefreshRank}
        type="button"
      >
        {rankLoading ? t.refreshingRank : t.refreshRank}
      </button>
      <button
        className="rounded-full border border-white/15 px-5 py-3 text-sm font-bold text-slate-200 transition hover:-translate-y-0.5 hover:bg-white/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={!currentPlayer || reportLoading}
        onClick={onGenerateReport}
        type="button"
      >
        {reportLoading ? t.analyzing : t.generateReport}
      </button>
      <button
        className="rounded-full border border-red-400/40 px-5 py-3 text-sm font-bold text-red-200 transition hover:-translate-y-0.5 hover:bg-red-500/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={resetLoading}
        onClick={onResetAllData}
        type="button"
      >
        {resetLoading ? t.resetting : t.resetLocalData}
      </button>
    </div>

    {(currentPlayer || lookupResult || appInfo) && (
      <div className="mt-8 space-y-4">
        {currentPlayer && (
          <div className="rounded-2xl border border-tactical-cyan/30 bg-tactical-cyan/10 p-5">
            <p className="text-sm font-semibold text-tactical-cyan">{t.currentPlayer}</p>
            <h3 className="mt-2 text-2xl font-bold text-white">
              {currentPlayer.name}#{currentPlayer.tag}
            </h3>
            <p className="mt-2 text-sm text-slate-300">
              {t.region} {currentPlayer.region || region.toUpperCase()} · {t.level}{' '}
              {currentPlayer.accountLevel || t.unknown}
            </p>
            <p className="mt-2 break-all text-xs text-slate-500">PUUID: {currentPlayer.puuid}</p>
          </div>
        )}
        {rank && (
          <div className="rounded-2xl border border-emerald-300/30 bg-emerald-300/10 p-5">
            <p className="text-sm font-semibold text-emerald-200">{t.latestRank}</p>
            <h3 className="mt-2 text-2xl font-bold text-white">{rank.tierName || t.unratedUnknown}</h3>
            <p className="mt-2 text-sm text-slate-300">
              {rank.rankingInTier} RR · elo {rank.elo || t.unknown} · {t.lastGame}{' '}
              {rank.mmrChangeToLast > 0 ? '+' : ''}
              {rank.mmrChangeToLast}
            </p>
            <p className="mt-2 text-xs text-slate-500">
              {t.region} {rank.region || region.toUpperCase()} · {t.fetched} {rank.fetchedAt || t.unknown}
            </p>
          </div>
        )}
        {lookupResult && (
          <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-5 text-sm text-slate-300">
            {t.provider}: {lookupResult.provider} · {t.consentVersion}: {lookupResult.consentVersion}
          </div>
        )}
        {appInfo && (
          <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-5">
            <p className="text-sm font-semibold text-tactical-cyan">{appInfo.status}</p>
            <h3 className="mt-2 text-xl font-bold text-white">{appInfo.name}</h3>
            <div className="mt-4 flex flex-wrap gap-2">
              {appInfo.stack.map((item) => (
                <span className="rounded-full bg-black/30 px-3 py-1 text-xs font-semibold text-slate-200" key={item}>
                  {item}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    )}

    {report && <ReportPanel report={report} t={t} />}
  </section>
);
