import { useEffect, useState } from 'react'
import type { ValorantAgent } from '../../types'

const AGENT_LOCALE = 'vi-VN'
const API_URL = `https://valorant-api.com/v1/agents?isPlayableCharacter=true&language=${AGENT_LOCALE}`

// Role display order (vi-VN from Valorant API)
export const ROLE_ORDER = ['Đối đầu', 'Khởi tranh', 'Kiểm soát', 'Hộ vệ']

// Role colors — Vietnamese + English keys for saved map markers
export const ROLE_COLORS: Record<string, string> = {
  'Đối đầu': '#ff4655',
  'Khởi tranh': '#fbbf24',
  'Kiểm soát': '#42e8f3',
  'Hộ vệ': '#4ade80',
  Duelist: '#ff4655',
  Initiator: '#fbbf24',
  Controller: '#42e8f3',
  Sentinel: '#4ade80',
}

let agentCache: ValorantAgent[] | null = null
let agentCacheLocale: string | null = null

export const useAgents = () => {
  const [agents, setAgents] = useState<ValorantAgent[]>(agentCache ?? [])
  const [loading, setLoading] = useState(!agentCache)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (agentCache && agentCacheLocale === AGENT_LOCALE) {
      setAgents(agentCache)
      setLoading(false)
      return
    }

    const fetchAgents = async () => {
      try {
        const res = await fetch(API_URL)
        const json = (await res.json()) as { data: ValorantAgent[] }
        const sorted = [...json.data].sort((a, b) => a.displayName.localeCompare(b.displayName))
        agentCache = sorted
        agentCacheLocale = AGENT_LOCALE
        setAgents(sorted)
      } catch {
        setError('Không thể tải danh sách agent. Kiểm tra kết nối mạng.')
      } finally {
        setLoading(false)
      }
    }

    void fetchAgents()
  }, [])

  const byRole = agents.reduce<Record<string, ValorantAgent[]>>((acc, agent) => {
    const role = agent.role?.displayName ?? 'Unknown'
    if (!acc[role]) acc[role] = []
    acc[role].push(agent)
    return acc
  }, {})

  return { agents, byRole, loading, error }
}
