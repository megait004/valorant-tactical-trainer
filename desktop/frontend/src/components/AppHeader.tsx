import type { AppStatusProps } from './types';

export const AppHeader = ({ status }: AppStatusProps) => (
  <header className="flex flex-col gap-4 border-b border-white/10 pb-6 md:flex-row md:items-center md:justify-between">
    <div>
      <p className="text-xs font-semibold uppercase tracking-[0.36em] text-tactical-red">Valorant Tactical Trainer</p>
      <h1 className="mt-3 text-4xl font-black tracking-tight text-white md:text-6xl">Desktop tactics lab</h1>
    </div>
    <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-slate-300 shadow-2xl shadow-black/20">
      <span className="mr-2 inline-flex h-2.5 w-2.5 rounded-full bg-tactical-cyan shadow-[0_0_18px_#49f5d4]" />
      {status}
    </div>
  </header>
);
