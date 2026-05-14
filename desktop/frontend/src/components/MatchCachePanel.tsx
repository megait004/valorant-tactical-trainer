import type { MatchCachePanelProps } from './types';

const plannedModules = [
  'Consent-gated player lookup',
  'Henrik API adapter with rate limit guard',
  'SQLite local match storage',
  'Evidence-based tactical reports',
  'Training recommendations and drill tracking',
];

export const MatchCachePanel = ({ matches, t }: MatchCachePanelProps) => (
  <aside className="rounded-[2rem] border border-white/10 bg-tactical-900/80 p-6 shadow-2xl shadow-black/30">
    <div className="flex items-center justify-between gap-4">
      <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">{t.matchCache}</p>
      <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-semibold text-slate-300">{matches.length} stored</span>
    </div>
    {matches.length > 0 ? (
      <div className="mt-6 max-h-[34rem] space-y-4 overflow-y-auto pr-1">
        {matches.map((match) => (
          <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-4" key={match.matchId}>
            <div className="flex items-start justify-between gap-4">
              <div>
                <h3 className="font-semibold text-white">{match.mapName || 'Unknown map'}</h3>
                <p className="mt-1 text-sm text-slate-400">
                  {match.agent || 'Unknown agent'} · {match.mode || match.queue || 'mode unknown'}
                </p>
              </div>
              <span className="rounded-full bg-tactical-red/15 px-3 py-1 text-xs font-black text-tactical-red">
                {match.kills}/{match.deaths}/{match.assists}
              </span>
            </div>
            <div className="mt-4 grid grid-cols-3 gap-2 text-xs text-slate-400">
              <span>{match.roundsPlayed || 0} rounds</span>
              <span>{match.headshots || 0} HS</span>
              <span>{match.damageMade || 0} dmg</span>
            </div>
            <p className="mt-3 truncate text-xs text-slate-600">{match.matchId}</p>
          </div>
        ))}
      </div>
    ) : (
      <div className="mt-6 space-y-4">
        {plannedModules.map((module, index) => (
          <div className="flex gap-4 rounded-2xl border border-white/10 bg-white/[0.04] p-4" key={module}>
            <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-tactical-red/15 text-sm font-black text-tactical-red">
              {index + 1}
            </span>
            <div>
              <h3 className="font-semibold text-white">{module}</h3>
              <p className="mt-1 text-sm text-slate-400">planned vertical slice</p>
            </div>
          </div>
        ))}
      </div>
    )}
  </aside>
);
