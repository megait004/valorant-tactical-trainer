import type { VirtualAssistantPanelProps } from './types';

const maps = ['Ascent', 'Bind', 'Haven'];
const agents = ['', 'Sova', 'Viper', 'Brimstone'];

export const VirtualAssistantPanel = ({
  agent,
  assistantLoading,
  credits,
  mapName,
  onAgentChange,
  onCreditsChange,
  onMapNameChange,
  onPhaseChange,
  onPreviousOutcomeChange,
  onQueryAssistant,
  onSideChange,
  phase,
  previousOutcome,
  result,
  side,
  t,
}: VirtualAssistantPanelProps) => (
  <aside className="rounded-[2rem] border border-cyan-300/20 bg-cyan-950/20 p-6 shadow-2xl shadow-cyan-950/30 backdrop-blur">
    <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
      <div>
        <p className="text-xs font-black uppercase tracking-[0.35em] text-cyan-200/80">{t.virtualAssistant}</p>
        <h2 className="mt-2 text-2xl font-black text-white">{t.assistantSubtitle}</h2>
      </div>
      <span className="rounded-full border border-emerald-300/30 bg-emerald-300/10 px-3 py-1 text-xs font-bold text-emerald-200">
        {t.assistantBadge}
      </span>
    </div>

    <div className="mt-6 grid gap-3 sm:grid-cols-2">
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Map
        <select
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          onChange={(event) => onMapNameChange(event.target.value)}
          value={mapName}
        >
          {maps.map((value) => (
            <option key={value} value={value}>
              {value}
            </option>
          ))}
        </select>
      </label>
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Agent
        <select
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          onChange={(event) => onAgentChange(event.target.value)}
          value={agent}
        >
          {agents.map((value) => (
            <option key={value || 'any'} value={value}>
              {value || t.anyAgent}
            </option>
          ))}
        </select>
      </label>
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Side
        <select
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          onChange={(event) => onSideChange(event.target.value)}
          value={side}
        >
          <option value="attack">Attack</option>
          <option value="defense">Defense</option>
          <option value="both">Both</option>
        </select>
      </label>
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Phase
        <select
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          onChange={(event) => onPhaseChange(event.target.value)}
          value={phase}
        >
          <option value="prematch">Pre-match</option>
          <option value="ingame">In-game</option>
        </select>
      </label>
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Credits
        <input
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          min="0"
          onChange={(event) => onCreditsChange(Number(event.target.value))}
          type="number"
          value={credits}
        />
      </label>
      <label className="text-xs font-bold uppercase tracking-[0.18em] text-slate-400">
        Previous round
        <select
          className="mt-2 w-full rounded-2xl border border-white/10 bg-slate-950/60 px-4 py-3 text-sm text-white outline-none focus:border-cyan-300/70"
          onChange={(event) => onPreviousOutcomeChange(event.target.value)}
          value={previousOutcome}
        >
          <option value="win">Win</option>
          <option value="loss">Loss</option>
        </select>
      </label>
    </div>

    <button
      className="mt-5 w-full rounded-full bg-cyan-300 px-5 py-3 text-sm font-black text-slate-950 transition hover:-translate-y-0.5 hover:bg-white disabled:cursor-not-allowed disabled:bg-slate-600 disabled:text-slate-300"
      disabled={assistantLoading}
      onClick={onQueryAssistant}
      type="button"
    >
      {assistantLoading ? t.assistantLoading : t.assistantButton}
    </button>

    {result && (
      <div className="mt-6 space-y-4">
        <section className="rounded-3xl border border-amber-300/20 bg-amber-300/10 p-4">
          <p className="text-xs font-black uppercase tracking-[0.25em] text-amber-200">{t.economyManager}</p>
          <h3 className="mt-2 text-xl font-black text-white">{result.economyAdvice.plan}</h3>
          <p className="mt-2 text-sm text-amber-50/90">{result.economyAdvice.summary}</p>
          <p className="mt-2 text-xs leading-5 text-amber-100/70">{result.economyAdvice.reminder}</p>
        </section>

        <div className="grid gap-3">
          {result.cards.map((card) => (
            <article key={card.id} className="rounded-3xl border border-white/10 bg-white/[0.04] p-4">
              <div className="flex flex-wrap items-center gap-2 text-[0.68rem] font-black uppercase tracking-[0.18em] text-cyan-200/80">
                <span>{card.category}</span>
                <span className="text-slate-600">/</span>
                <span>{card.mapName}</span>
                {card.agent && <span>{card.agent}</span>}
              </div>
              <h3 className="mt-2 text-lg font-black text-white">{card.title}</h3>
              <p className="mt-2 text-sm leading-6 text-slate-300">{card.summary}</p>
              <p className="mt-3 rounded-2xl border border-cyan-300/15 bg-cyan-300/10 p-3 text-sm font-semibold leading-6 text-cyan-50">
                {card.action}
              </p>
              <p className="mt-2 text-xs text-slate-500">{card.safetyNotes}</p>
            </article>
          ))}
        </div>

        <div className="rounded-3xl border border-emerald-300/15 bg-emerald-300/10 p-4 text-xs leading-5 text-emerald-100/80">
          {result.safetyNotes.map((note) => (
            <p key={note}>{note}</p>
          ))}
        </div>
      </div>
    )}
  </aside>
);
