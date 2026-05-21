import { useState } from 'react'
import type { FC } from 'react'
import type { ValorantAgent, ValorantAbility } from '../../types'
import { ROLE_COLORS, ROLE_ORDER, useAgents } from './useAgents'

const SLOT_KEYS: Record<string, string> = {
  Ability1: 'E',
  Ability2: 'Q',
  Grenade: 'C',
  Ultimate: 'X',
}

const SLOT_LABEL: Record<string, string> = {
  Ability1: 'Kỹ năng E',
  Ability2: 'Kỹ năng Q',
  Grenade: 'Kỹ năng C',
  Ultimate: 'Ultimate X',
}

// ─── Agent Detail ──────────────────────────────────────────────────────────────

const AbilityCard: FC<{ ability: ValorantAbility }> = ({ ability }) => {
  const [expanded, setExpanded] = useState(false)
  const key = SLOT_KEYS[ability.slot] ?? ability.slot
  const slotLabel = SLOT_LABEL[ability.slot] ?? ability.slot
  const isUlt = ability.slot === 'Ultimate'

  return (
    <button
      type="button"
      onClick={() => setExpanded((v) => !v)}
      className={`w-full rounded-xl border text-left transition ${
        isUlt
          ? 'border-tactical-red/40 bg-tactical-red/10 hover:bg-tactical-red/20'
          : 'border-tactical-line bg-black/30 hover:bg-white/5'
      }`}
    >
      <div className="flex items-center gap-3 p-3">
        <div
          className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-lg ${
            isUlt ? 'bg-tactical-red/20' : 'bg-white/5'
          }`}
        >
          {ability.displayIcon ? (
            <img
              src={ability.displayIcon}
              alt={ability.displayName}
              className="h-7 w-7 object-contain"
              onError={(e) => {
                e.currentTarget.style.display = 'none'
              }}
            />
          ) : (
            <span className={`text-sm font-black ${isUlt ? 'text-tactical-red' : 'text-slate-400'}`}>{key}</span>
          )}
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span
              className={`rounded px-1.5 py-0.5 text-[9px] font-black uppercase tracking-widest ${
                isUlt ? 'bg-tactical-red/20 text-tactical-red' : 'bg-white/10 text-slate-400'
              }`}
            >
              {key}
            </span>
            <span className="truncate text-sm font-bold text-white">{ability.displayName}</span>
          </div>
          <p className="mt-0.5 text-[10px] text-slate-500">{slotLabel}</p>
        </div>
        <span className="shrink-0 text-slate-600">{expanded ? '▲' : '▼'}</span>
      </div>
      {expanded && (
        <p className="border-t border-white/5 px-3 pb-3 pt-2 text-xs leading-relaxed text-slate-300">
          {ability.description}
        </p>
      )}
    </button>
  )
}

const AgentDetail: FC<{ agent: ValorantAgent; onClose: () => void }> = ({ agent, onClose }) => {
  const role = agent.role?.displayName ?? 'Unknown'
  const roleColor = ROLE_COLORS[role] ?? '#42e8f3'
  const gradientColors = agent.backgroundGradientColors
  const bgGrad =
    gradientColors.length >= 2
      ? `linear-gradient(135deg, #${gradientColors[0]?.slice(0, 6) ?? '0f1923'}, #${gradientColors[1]?.slice(0, 6) ?? '0f1923'})`
      : undefined

  const abilityOrder = ['Ability2', 'Grenade', 'Ability1', 'Ultimate']
  const sortedAbilities = [...agent.abilities].sort(
    (a, b) => abilityOrder.indexOf(a.slot) - abilityOrder.indexOf(b.slot),
  )

  return (
    <div className="flex h-full flex-col overflow-hidden rounded-3xl border border-tactical-line bg-tactical-panel/90 shadow-2xl shadow-black/30">
      {/* Hero section */}
      <div
        className="relative flex items-end overflow-hidden px-5 pb-4 pt-5"
        style={{ background: bgGrad, minHeight: 180 }}
      >
        <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
        {(agent.fullPortrait ?? agent.bustPortrait) && (
          <img
            src={agent.fullPortrait ?? agent.bustPortrait ?? ''}
            alt={agent.displayName}
            className="pointer-events-none absolute bottom-0 right-0 h-48 object-contain drop-shadow-2xl"
            draggable={false}
          />
        )}
        <div className="relative z-10 flex-1">
          <div className="flex items-center gap-2">
            {agent.role?.displayIcon && (
              <img
                src={agent.role.displayIcon}
                alt={role}
                className="h-5 w-5 object-contain opacity-90"
                style={{ filter: `drop-shadow(0 0 4px ${roleColor})` }}
              />
            )}
            <span className="text-xs font-black uppercase tracking-widest" style={{ color: roleColor }}>
              {role}
            </span>
          </div>
          <h2 className="mt-1 text-3xl font-black tracking-tight text-white">{agent.displayName}</h2>
        </div>
        <button
          type="button"
          onClick={onClose}
          className="relative z-10 ml-2 shrink-0 rounded-xl border border-white/10 bg-black/40 px-3 py-1.5 text-xs font-bold text-slate-300 hover:bg-black/60"
        >
          ✕
        </button>
      </div>

      {/* Bio */}
      <p className="border-b border-tactical-line px-5 py-3 text-xs leading-relaxed text-slate-400">
        {agent.description}
      </p>

      {/* Abilities */}
      <div className="flex-1 overflow-y-auto px-5 py-4">
        <p className="mb-3 text-xs font-black uppercase tracking-[0.18em] text-slate-500">Kỹ năng</p>
        <div className="space-y-2">
          {sortedAbilities.map((ability) => (
            <AbilityCard key={ability.slot} ability={ability} />
          ))}
        </div>
      </div>
    </div>
  )
}

// ─── Agent Card ────────────────────────────────────────────────────────────────

const AgentCard: FC<{
  agent: ValorantAgent
  selected: boolean
  onClick: () => void
}> = ({ agent, selected, onClick }) => {
  const role = agent.role?.displayName ?? 'Unknown'
  const roleColor = ROLE_COLORS[role] ?? '#42e8f3'

  return (
    <button
      type="button"
      onClick={onClick}
      title={agent.displayName}
      className={`group relative flex flex-col items-center gap-1 rounded-2xl border p-2 text-center transition ${
        selected
          ? 'border-tactical-cyan/60 bg-tactical-cyan/10 shadow-lg shadow-tactical-cyan/10'
          : 'border-tactical-line bg-black/20 hover:border-white/20 hover:bg-white/[0.04]'
      }`}
    >
      <div className="relative h-16 w-full overflow-hidden rounded-xl bg-black/40">
        {agent.displayIcon ? (
          <img
            src={agent.displayIcon}
            alt={agent.displayName}
            className="h-full w-full object-cover transition group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-2xl text-slate-600">?</div>
        )}
        <div
          className="absolute bottom-0 left-0 right-0 h-1 rounded-b-xl opacity-80"
          style={{ backgroundColor: roleColor }}
        />
      </div>
      <span
        className={`w-full truncate text-[10px] font-black uppercase leading-tight transition ${
          selected ? 'text-tactical-cyan' : 'text-slate-200 group-hover:text-white'
        }`}
      >
        {agent.displayName}
      </span>
    </button>
  )
}

// ─── Role Tab ─────────────────────────────────────────────────────────────────

const RoleTab: FC<{
  role: string
  roleIcon?: string
  count: number
  active: boolean
  onClick: () => void
}> = ({ role, roleIcon, count, active, onClick }) => {
  const color = ROLE_COLORS[role] ?? '#42e8f3'
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex items-center gap-2 rounded-xl border px-3 py-2 text-xs font-black uppercase tracking-wider transition ${
        active
          ? 'border-white/20 text-white shadow-md'
          : 'border-transparent text-slate-500 hover:border-white/10 hover:text-slate-300'
      }`}
      style={active ? { backgroundColor: `${color}22`, borderColor: `${color}55`, color } : {}}
    >
      {roleIcon && (
        <img
          src={roleIcon}
          alt={role}
          className="h-4 w-4 object-contain"
          style={active ? { filter: `drop-shadow(0 0 4px ${color})` } : { opacity: 0.6 }}
        />
      )}
      {role}
      <span
        className="ml-0.5 rounded-full px-1.5 py-0.5 text-[9px]"
        style={active ? { backgroundColor: `${color}30`, color } : { backgroundColor: '#ffffff10', color: '#64748b' }}
      >
        {count}
      </span>
    </button>
  )
}

// ─── Agent Browser Panel ───────────────────────────────────────────────────────

const AgentBrowserPanel: FC = () => {
  const { agents, byRole, loading, error } = useAgents()
  const [activeRole, setActiveRole] = useState<string | null>(null)
  const [selectedAgent, setSelectedAgent] = useState<ValorantAgent | null>(null)

  const roles = ROLE_ORDER.filter((r) => byRole[r]?.length)
  const roleIcons = Object.fromEntries(
    agents
      .filter((a) => a.role?.displayName)
      .map((a) => [a.role.displayName, a.role.displayIcon]),
  ) as Record<string, string>

  const displayedAgents = activeRole ? (byRole[activeRole] ?? []) : agents

  if (loading) {
    return (
      <div className="flex min-h-[300px] flex-col items-center justify-center gap-4 rounded-3xl border border-tactical-line bg-tactical-panel/90 p-8">
        <div className="h-10 w-10 animate-spin rounded-full border-2 border-tactical-cyan border-t-transparent" />
        <p className="text-sm text-slate-400">Đang tải agent (tiếng Việt) từ Valorant API…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex min-h-[200px] flex-col items-center justify-center gap-3 rounded-3xl border border-tactical-red/30 bg-tactical-panel/90 p-8">
        <p className="text-sm font-bold text-tactical-red">{error}</p>
        <p className="text-xs text-slate-500">Cần kết nối internet để tải dữ liệu agent lần đầu.</p>
      </div>
    )
  }

  return (
    <section className="grid gap-5 xl:grid-cols-[1fr_380px]">
      {/* Left: browse */}
      <div className="space-y-5">
        {/* Header */}
        <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">
          <p className="text-xs font-black uppercase tracking-[0.2em] text-tactical-cyan">Agent Browser</p>
          <p className="mt-1 text-2xl font-black text-white">Danh sách Agent Valorant</p>
          <p className="mt-1 text-sm text-slate-400">
            {agents.length} agent · theo vai trò · click để xem chi tiết kỹ năng
          </p>

          {/* Role filter */}
          <div className="mt-4 flex flex-wrap gap-2">
            <button
              type="button"
              onClick={() => setActiveRole(null)}
              className={`rounded-xl border px-3 py-2 text-xs font-black uppercase tracking-wider transition ${
                activeRole
                  ? 'border-transparent text-slate-500 hover:border-white/10 hover:text-slate-300'
                  : 'border-white/20 bg-white/10 text-white'
              }`}
            >
              Tất cả ({agents.length})
            </button>
            {roles.map((role) => (
              <RoleTab
                key={role}
                role={role}
                roleIcon={roleIcons[role]}
                count={byRole[role]?.length ?? 0}
                active={activeRole === role}
                onClick={() => setActiveRole(activeRole === role ? null : role)}
              />
            ))}
          </div>
        </div>

        {/* Agent grid — by role if showing all */}
        {activeRole ? (
          <div className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20">
            <div className="flex items-center gap-2 mb-4">
              {roleIcons[activeRole] && (
                <img src={roleIcons[activeRole]} alt={activeRole} className="h-5 w-5 object-contain" />
              )}
              <p className="font-black uppercase tracking-widest text-sm" style={{ color: ROLE_COLORS[activeRole] ?? '#42e8f3' }}>
                {activeRole}
              </p>
              <span className="text-xs text-slate-500">— {displayedAgents.length} agent</span>
            </div>
            <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6">
              {displayedAgents.map((agent) => (
                <AgentCard
                  key={agent.uuid}
                  agent={agent}
                  selected={selectedAgent?.uuid === agent.uuid}
                  onClick={() => setSelectedAgent(agent.uuid === selectedAgent?.uuid ? null : agent)}
                />
              ))}
            </div>
          </div>
        ) : (
          roles.map((role) => {
            const roleAgents = byRole[role] ?? []
            if (!roleAgents.length) return null
            const roleColor = ROLE_COLORS[role] ?? '#42e8f3'
            return (
              <div
                key={role}
                className="rounded-3xl border border-tactical-line bg-tactical-panel/90 p-5 shadow-xl shadow-black/20"
              >
                <div className="mb-4 flex items-center gap-2">
                  {roleIcons[role] && (
                    <img
                      src={roleIcons[role]}
                      alt={role}
                      className="h-5 w-5 object-contain"
                      style={{ filter: `drop-shadow(0 0 4px ${roleColor})` }}
                    />
                  )}
                  <p className="font-black uppercase tracking-widest text-sm" style={{ color: roleColor }}>
                    {role}
                  </p>
                  <span className="text-xs text-slate-600">— {roleAgents.length} agent</span>
                  <button
                    type="button"
                    onClick={() => setActiveRole(role)}
                    className="ml-auto text-[10px] text-slate-600 hover:text-slate-400"
                  >
                    Lọc →
                  </button>
                </div>
                <p className="mb-3 text-xs text-slate-500">{roleAgents[0]?.role?.description}</p>
                <div className="grid grid-cols-3 gap-2 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8">
                  {roleAgents.map((agent) => (
                    <AgentCard
                      key={agent.uuid}
                      agent={agent}
                      selected={selectedAgent?.uuid === agent.uuid}
                      onClick={() => setSelectedAgent(agent.uuid === selectedAgent?.uuid ? null : agent)}
                    />
                  ))}
                </div>
              </div>
            )
          })
        )}
      </div>

      {/* Right: detail */}
      <aside className="sticky top-4 self-start">
        {selectedAgent ? (
          <AgentDetail agent={selectedAgent} onClose={() => setSelectedAgent(null)} />
        ) : (
          <div className="flex min-h-[300px] flex-col items-center justify-center gap-3 rounded-3xl border border-tactical-line bg-tactical-panel/90 p-8 text-center shadow-xl shadow-black/20">
            <div className="h-12 w-12 rounded-full border border-tactical-line bg-black/30 opacity-30" />
            <p className="text-sm font-bold text-slate-400">Chọn một agent để xem chi tiết kỹ năng</p>
            <p className="text-xs text-slate-600">Click vào bất kỳ agent nào ở bên trái</p>
          </div>
        )}
      </aside>
    </section>
  )
}

export default AgentBrowserPanel
