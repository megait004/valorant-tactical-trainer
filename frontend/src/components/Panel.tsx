// Layout primitives dùng chung trong toàn app — Panel (card), MetricPill/Card,
// EmptyState. Không có business logic; tách riêng để các feature/component
// khác dễ tái sử dụng và đảm bảo design system nhất quán.

import type { ReactNode } from 'react'

export const Panel = ({ children }: { children: ReactNode }) => (
  <section className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">
    {children}
  </section>
)

export const MetricPill = ({ label, value }: { label: string; value: string }) => (
  <div className="rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3">
    <div className="text-[11px] uppercase tracking-[0.2em] text-slate-500">{label}</div>
    <div className="mt-1 font-bold text-white">{value}</div>
  </div>
)

export const MetricCard = ({ label, value }: { label: string; value: string }) => (
  <div className="rounded-2xl border border-tactical-line bg-black/20 p-4">
    <p className="text-xs uppercase tracking-[0.18em] text-slate-500">{label}</p>
    <p className="mt-2 text-2xl font-black text-white">{value}</p>
  </div>
)

export const EmptyState = ({ title }: { title: string }) => (
  <div className="rounded-3xl border border-tactical-line bg-tactical-panel p-6 text-slate-300">{title}</div>
)
