import type { SettingsPanelProps } from './types';

export const SettingsPanel = ({
  apiKey,
  loading,
  settings,
  onApiKeyChange,
  onClearExpiredCache,
  onExportLocalData,
  onSaveSettings,
}: SettingsPanelProps) => (
  <aside className="rounded-[2rem] border border-white/10 bg-tactical-900/80 p-6 shadow-2xl shadow-black/30">
    <div className="flex items-start justify-between gap-4">
      <div>
        <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">Settings</p>
        <h2 className="mt-3 text-2xl font-bold text-white">Local data controls</h2>
      </div>
      <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-semibold text-slate-300">
        {settings?.apiKeyConfigured ? 'key saved' : 'no key'}
      </span>
    </div>

    <label className="mt-5 block space-y-2">
      <span className="text-sm font-semibold text-slate-300">API key</span>
      <input
        className="w-full rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-red"
        onChange={(event) => onApiKeyChange(event.target.value)}
        placeholder={settings?.apiKeyConfigured ? 'leave empty and save to clear key' : 'optional Henrik API key'}
        type="password"
        value={apiKey}
      />
    </label>

    <div className="mt-4 flex flex-wrap gap-3">
      <button
        className="rounded-full bg-tactical-red px-4 py-2 text-sm font-bold text-white transition hover:-translate-y-0.5 hover:bg-red-400 disabled:cursor-not-allowed disabled:bg-slate-700 disabled:text-slate-400"
        disabled={loading}
        onClick={onSaveSettings}
        type="button"
      >
        {loading ? 'Saving...' : 'Save key'}
      </button>
      <button
        className="rounded-full border border-white/15 px-4 py-2 text-sm font-bold text-slate-200 transition hover:-translate-y-0.5 hover:bg-white/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={loading}
        onClick={onClearExpiredCache}
        type="button"
      >
        Clear expired cache
      </button>
      <button
        className="rounded-full border border-emerald-300/40 px-4 py-2 text-sm font-bold text-emerald-200 transition hover:-translate-y-0.5 hover:bg-emerald-300/10 disabled:cursor-not-allowed disabled:border-slate-700 disabled:text-slate-500"
        disabled={loading}
        onClick={onExportLocalData}
        type="button"
      >
        Export JSON
      </button>
    </div>

    <div className="mt-5 grid grid-cols-2 gap-3 text-sm text-slate-300">
      <Stat label="players" value={settings?.players ?? 0} />
      <Stat label="matches" value={settings?.matches ?? 0} />
      <Stat label="rank snaps" value={settings?.rankSnapshots ?? 0} />
      <Stat label="reports" value={settings?.reports ?? 0} />
      <Stat label="cache" value={settings?.cacheEntries ?? 0} />
      <Stat label="expired" value={settings?.expiredCacheEntries ?? 0} />
    </div>

    <p className="mt-5 break-all rounded-2xl border border-white/10 bg-black/20 p-3 text-xs text-slate-500">
      SQLite: {settings?.dataPath || 'not loaded'}
    </p>
    <p className="mt-3 text-xs leading-5 text-slate-500">
      Export excludes the saved API key value and raw provider payloads.
    </p>
  </aside>
);

const Stat = ({ label, value }: { label: string; value: number }) => (
  <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-3">
    <p className="text-xs uppercase tracking-[0.16em] text-slate-500">{label}</p>
    <p className="mt-1 text-xl font-black text-white">{value}</p>
  </div>
);
