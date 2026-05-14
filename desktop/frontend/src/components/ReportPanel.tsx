import type { ReportPanelProps } from './types';

export const ReportPanel = ({ report, t }: ReportPanelProps) => (
  <section className="mt-8 rounded-[2rem] border border-white/10 bg-black/20 p-5">
    <div className="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
      <div>
        <p className="text-sm font-medium uppercase tracking-[0.24em] text-slate-400">{t.tacticalReport}</p>
        <h3 className="mt-2 text-2xl font-bold text-white">{report.summary}</h3>
      </div>
      <div className="grid grid-cols-3 gap-2 text-center text-xs text-slate-300">
        <span className="rounded-2xl bg-white/10 px-3 py-2">KDA {report.averageKda}</span>
        <span className="rounded-2xl bg-white/10 px-3 py-2">HS {report.headshotPercent}%</span>
        <span className="rounded-2xl bg-white/10 px-3 py-2">DMG {report.averageDamage}</span>
      </div>
    </div>

    <div className="mt-5 grid gap-4 lg:grid-cols-2">
      <div className="space-y-3">
        <h4 className="font-semibold text-white">{t.findings}</h4>
        {report.findings.map((finding) => (
          <article className="rounded-2xl border border-white/10 bg-white/[0.04] p-4" key={`${finding.type}-${finding.title}`}>
            <div className="flex items-center justify-between gap-3">
              <h5 className="font-semibold text-white">{finding.title}</h5>
              <span className="rounded-full bg-tactical-red/15 px-2 py-1 text-xs font-bold text-tactical-red">
                {finding.severity}
              </span>
            </div>
            <p className="mt-2 text-sm leading-6 text-slate-300">{finding.description}</p>
            <p className="mt-2 text-xs text-slate-500">confidence {finding.confidence}</p>
          </article>
        ))}
      </div>

      <div className="space-y-3">
        <h4 className="font-semibold text-white">{t.trainingRecommendations}</h4>
        {report.recommendations.map((recommendation) => (
          <article className="rounded-2xl border border-tactical-cyan/20 bg-tactical-cyan/10 p-4" key={`${recommendation.priority}-${recommendation.title}`}>
            <div className="flex items-center justify-between gap-3">
              <h5 className="font-semibold text-white">{recommendation.title}</h5>
              <span className="rounded-full bg-black/25 px-2 py-1 text-xs font-bold text-tactical-cyan">
                {recommendation.priority}
              </span>
            </div>
            <p className="mt-2 text-sm leading-6 text-slate-200">{recommendation.drill}</p>
            <p className="mt-2 text-xs text-slate-500">{recommendation.reason}</p>
          </article>
        ))}
      </div>
    </div>
  </section>
);
