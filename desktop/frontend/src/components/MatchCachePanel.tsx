import { useMemo, useState } from 'react';
import type { MatchCachePanelProps } from './types';

export const MatchCachePanel = ({ matches, t }: MatchCachePanelProps) => {
  const [search, setSearch] = useState('');
  const [mapFilter, setMapFilter] = useState('');
  const [selectedMatchId, setSelectedMatchId] = useState('');

  const maps = useMemo(
    () => Array.from(new Set(matches.map((match) => match.mapName).filter(Boolean))).sort(),
    [matches],
  );
  const filteredMatches = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase();
    return matches.filter((match) => {
      const sameMap = mapFilter === '' || match.mapName === mapFilter;
      const haystack = `${match.mapName} ${match.agent} ${match.mode} ${match.queue} ${match.region}`.toLowerCase();
      return sameMap && (normalizedSearch === '' || haystack.includes(normalizedSearch));
    });
  }, [mapFilter, matches, search]);
  const selectedMatch =
    filteredMatches.find((match) => match.matchId === selectedMatchId) ?? filteredMatches[0] ?? matches[0] ?? null;
  const plannedModules = [
    t.moduleConsentLookup,
    t.moduleProviderAdapter,
    t.moduleLocalStorage,
    t.moduleTacticalReports,
    t.moduleTraining,
  ];

  return (
    <aside className="rounded-[2rem] border border-white/10 bg-tactical-900/80 p-6 shadow-2xl shadow-black/30">
      <div className="flex items-center justify-between gap-4">
        <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">{t.matchCache}</p>
        <span className="rounded-full bg-white/10 px-3 py-1 text-xs font-semibold text-slate-300">
          {filteredMatches.length}/{matches.length} {t.stored}
        </span>
      </div>
      {matches.length > 0 ? (
        <>
          <div className="mt-5 grid gap-3 sm:grid-cols-[1fr_0.7fr]">
            <input
              className="rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none transition placeholder:text-slate-600 focus:border-tactical-cyan"
              onChange={(event) => setSearch(event.target.value)}
              placeholder={t.searchMatches}
              type="search"
              value={search}
            />
            <select
              className="rounded-2xl border border-white/10 bg-black/30 px-4 py-3 text-sm text-white outline-none transition focus:border-tactical-cyan"
              onChange={(event) => setMapFilter(event.target.value)}
              value={mapFilter}
            >
              <option value="">{t.allMaps}</option>
              {maps.map((map) => (
                <option key={map} value={map}>
                  {map}
                </option>
              ))}
            </select>
          </div>

          {selectedMatch && (
            <section className="mt-5 rounded-3xl border border-tactical-cyan/20 bg-tactical-cyan/10 p-4">
              <p className="text-xs font-black uppercase tracking-[0.2em] text-tactical-cyan">{t.selectedMatch}</p>
              <h3 className="mt-2 text-xl font-black text-white">{selectedMatch.mapName || t.unknownMap}</h3>
              <p className="mt-1 text-sm text-slate-300">
                {selectedMatch.agent || t.unknownAgent} · {selectedMatch.mode || selectedMatch.queue || t.unknownMode} ·{' '}
                {selectedMatch.region || t.unknownRegion}
              </p>
              <div className="mt-4 grid grid-cols-2 gap-2 text-sm text-slate-200 sm:grid-cols-4">
                <Detail label="K/D/A" value={`${selectedMatch.kills}/${selectedMatch.deaths}/${selectedMatch.assists}`} />
                <Detail label={t.rounds} value={selectedMatch.roundsPlayed || 0} />
                <Detail label="HS" value={selectedMatch.headshots || 0} />
                <Detail label={t.damageShort} value={selectedMatch.damageMade || 0} />
              </div>
              <p className="mt-3 break-all text-xs text-slate-500">{selectedMatch.matchId}</p>
            </section>
          )}

          <div className="mt-6 max-h-[28rem] space-y-4 overflow-y-auto pr-1">
            {filteredMatches.map((match) => (
              <button
                className={`w-full rounded-2xl border p-4 text-left transition hover:-translate-y-0.5 ${
                  match.matchId === selectedMatch?.matchId
                    ? 'border-tactical-cyan/40 bg-tactical-cyan/10'
                    : 'border-white/10 bg-white/[0.04] hover:bg-white/[0.08]'
                }`}
                key={match.matchId}
                onClick={() => setSelectedMatchId(match.matchId)}
                type="button"
              >
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <h3 className="font-semibold text-white">{match.mapName || t.unknownMap}</h3>
                    <p className="mt-1 text-sm text-slate-400">
                      {match.agent || t.unknownAgent} · {match.mode || match.queue || t.unknownMode}
                    </p>
                  </div>
                  <span className="rounded-full bg-tactical-red/15 px-3 py-1 text-xs font-black text-tactical-red">
                    {match.kills}/{match.deaths}/{match.assists}
                  </span>
                </div>
                <div className="mt-4 grid grid-cols-3 gap-2 text-xs text-slate-400">
                  <span>{match.roundsPlayed || 0} {t.rounds}</span>
                  <span>{match.headshots || 0} HS</span>
                  <span>{match.damageMade || 0} {t.damageShort}</span>
                </div>
              </button>
            ))}
          </div>
        </>
      ) : (
        <div className="mt-6 space-y-4">
          {plannedModules.map((module, index) => (
            <div className="flex gap-4 rounded-2xl border border-white/10 bg-white/[0.04] p-4" key={module}>
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-tactical-red/15 text-sm font-black text-tactical-red">
                {index + 1}
              </span>
              <div>
                <h3 className="font-semibold text-white">{module}</h3>
                <p className="mt-1 text-sm text-slate-400">{t.plannedVerticalSlice}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </aside>
  );
};

const Detail = ({ label, value }: { label: string; value: number | string }) => (
  <div className="rounded-2xl border border-white/10 bg-black/20 p-3">
    <p className="text-[0.65rem] font-black uppercase tracking-[0.16em] text-slate-500">{label}</p>
    <p className="mt-1 text-lg font-black text-white">{value}</p>
  </div>
);
