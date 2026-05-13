import { ReportPanel } from './ReportPanel';
import type { SetupPanelProps } from './types';

export const SetupPanel = ({
  appInfo,
  apiKey,
  canLookup,
  consent,
  currentPlayer,
  lookupLoading,
  lookupResult,
  matchLoading,
  name,
  region,
  report,
  reportLoading,
  resetLoading,
  tag,
  onApiKeyChange,
  onCheckCore,
  onConsentChange,
  onGenerateReport,
  onLookupPlayer,
  onNameChange,
  onRefreshMatches,
  onRegionChange,
  onResetAllData,
  onTagChange,
}: SetupPanelProps) => (
  <section className="rounded-[2rem] border border-white/10 bg-white/[0.06] p-7 shadow-2xl shadow-black/30 backdrop-blur">
    <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">Consent gate</p>
    <h2 className="mt-4 max-w-3xl text-3xl font-bold leading-tight text-white md:text-5xl">
      Start with explicit player consent before any Valorant data fetch.
    </h2>
    <p className="mt-5 max-w-2xl text-base leading-7 text-slate-300">
      The app uses Henrik unofficial VALORANT API, stores imported data locally in SQLite, and only fetches account data
      after the consent checkbox is confirmed.
    </p>

    <div className="mt-8 grid gap-4 md:grid-cols-2">
      <label className="space-y-2">
        <span className="text-sm font-semibold text-slate-300">Riot name</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
          onChange={(event) => onNameChange(event.target.value)}
          placeholder="ten nguoi choi"
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
        <span className="text-sm font-semibold text-slate-300">Region fallback</span>
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
        <span className="text-sm font-semibold text-slate-300">API key optional</span>
        <input
          className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
          onChange={(event) => onApiKeyChange(event.target.value)}
          placeholder="de trong neu chua co key"
          type="password"
          value={apiKey}
        />
      </label>
    </div>

    <label className="mt-6 flex gap-3 rounded-2xl border border-tactical-red/30 bg-tactical-red/10 p-4 text-sm leading-6 text-slate-200">
      <input
        checked={consent}
        className="mt-1 h-4 w-4 accent-tactical-red"
        onChange={(event) => onConsentChange(event.target.checked)}
        type="checkbox"
      />
      <span>
        I consent to fetch this player's account data from Henrik unofficial VALORANT API and store it locally on this
        machine for tactical analysis. I understand no Riot credentials are required.
      </span>
    </label>

    <div className="mt-6 flex flex-wrap gap-3">
      <button
        className="rounded-full bg-tactical-red px-5 py-3 text-sm font-bold text-white shadow-lg shadow-red-950/40 transition hover:-translate-y-0.5 hover:bg-red-400 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400 disabled:shadow-none"
        disabled={!canLookup}
        onClick={onLookupPlayer}
        type="button"
      >
        {lookupLoading ? 'Looking up...' : 'Lookup player'}
      </button>
      <button
        className="rounded-full bg-tactical-red px-5 py-3 text-sm font-bold text-white shadow-lg shadow-red-950/40 transition hover:-translate-y-0.5 hover:bg-red-400"
        onClick={onCheckCore}
        type="button"
      >
        Check Go core
      </button>
      <button
        className="rounded-full border border-tactical-cyan/40 px-5 py-3 text-sm font-bold text-tactical-cyan transition hover:-translate-y-0.5 hover:bg-tactical-cyan/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={!currentPlayer || matchLoading}
        onClick={onRefreshMatches}
        type="button"
      >
        {matchLoading ? 'Refreshing...' : 'Refresh matches'}
      </button>
      <button
        className="rounded-full border border-white/15 px-5 py-3 text-sm font-bold text-slate-200 transition hover:-translate-y-0.5 hover:bg-white/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={!currentPlayer || reportLoading}
        onClick={onGenerateReport}
        type="button"
      >
        {reportLoading ? 'Analyzing...' : 'Generate report'}
      </button>
      <button
        className="rounded-full border border-red-400/40 px-5 py-3 text-sm font-bold text-red-200 transition hover:-translate-y-0.5 hover:bg-red-500/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={resetLoading}
        onClick={onResetAllData}
        type="button"
      >
        {resetLoading ? 'Resetting...' : 'Reset local data'}
      </button>
    </div>

    {(currentPlayer || lookupResult || appInfo) && (
      <div className="mt-8 space-y-4">
        {currentPlayer && (
          <div className="rounded-2xl border border-tactical-cyan/30 bg-tactical-cyan/10 p-5">
            <p className="text-sm font-semibold text-tactical-cyan">current player</p>
            <h3 className="mt-2 text-2xl font-bold text-white">
              {currentPlayer.name}#{currentPlayer.tag}
            </h3>
            <p className="mt-2 text-sm text-slate-300">
              Region {currentPlayer.region || region.toUpperCase()} · Level {currentPlayer.accountLevel || 'unknown'}
            </p>
            <p className="mt-2 break-all text-xs text-slate-500">PUUID: {currentPlayer.puuid}</p>
          </div>
        )}
        {lookupResult && (
          <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-5 text-sm text-slate-300">
            Provider: {lookupResult.provider} · Consent version: {lookupResult.consentVersion}
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

    {report && <ReportPanel report={report} />}
  </section>
);
