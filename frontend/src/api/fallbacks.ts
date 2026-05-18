// Fallback values dùng khi Wails bridge chưa sẵn sàng (vd dev mode browser
// thuần, hoặc Go core chưa bind xong). Cũng dùng làm seed demo cho UI.

import type {
  APIStatus,
  AnalysisReport,
  AssistantSessionState,
  AssistantTipResult,
  ChatState,
  DataSettings,
  MapPlan,
} from '../types'

export const fallbackReport: AnalysisReport = {
  player: { name: 'giaphue', tagline: 'DATN', region: 'ap', primaryRole: 'Kiểm Soát (Controller)' },
  metrics: {
    matches: 5,
    rounds: 105,
    kd: 0.89,
    kda: 1.28,
    headshotPercent: 17.4,
    firstBloodRate: 0.07,
    firstDeathRate: 0.21,
    winRate: 0.4,
    weakestMap: 'Ascent',
    weakestMapWinRate: 0,
    weakestMapSample: 3,
    primaryRoleObserved: 'Kiểm Soát (Controller)',
  },
  findings: [
    {
      id: 'finding-first-death-non-duelist',
      title: 'First Death cao so với vai trò không phải Đối Đầu',
      severity: 'medium',
      confidence: 'high',
      detail: 'Tỉ lệ chết đầu round là 21% trong khi vai trò chính là Kiểm Soát (Controller).',
      evidence: [{ metric: 'first_death_rate', value: 0.21, sampleSize: 105, comparisonBaseline: 0.18 }],
    },
  ],
  mapBreakdown: [
    { name: 'Ascent', matches: 3, rounds: 61, kd: 0.67, winRate: 0, headshotPercent: 15.4 },
    { name: 'Bind', matches: 1, rounds: 23, kd: 1.5, winRate: 1, headshotPercent: 22.8 },
    { name: 'Haven', matches: 1, rounds: 21, kd: 1.13, winRate: 1, headshotPercent: 19.1 },
  ],
  agentBreakdown: [
    { name: 'Omen', matches: 3, rounds: 63, kd: 0.81, winRate: 0.33, headshotPercent: 17.1 },
    { name: 'Brimstone', matches: 1, rounds: 19, kd: 0.59, winRate: 0, headshotPercent: 14.2 },
    { name: 'Cypher', matches: 1, rounds: 23, kd: 1.5, winRate: 1, headshotPercent: 22.8 },
  ],
  practicePlan: [
    {
      day: 1,
      focus: 'Map gap protocol',
      map: 'Ascent',
      agent: 'Omen',
      duration: '25 phút',
      checklist: ['Viết 2 default route attack', 'Viết 2 setup defense', 'Chạy custom walkthrough timing utility'],
      evidence: 'Ascent win rate 0% trên 3 trận.',
    },
  ],
  recommendations: [
    {
      id: 'rec-survive-contact',
      findingId: 'finding-first-death-non-duelist',
      title: 'Giảm chết sớm bằng utility-before-peek drill',
      reason: 'Kiểm Soát (Controller) cần sống tới mid-late round để giữ smoke và call rotate.',
      drill: 'Custom 10 round defense: trước mỗi lần peek phải dùng smoke/info hoặc gọi teammate trade.',
      cadence: '15 phút/ngày trong 5 ngày',
    },
  ],
}

export const fallbackDataSettings: DataSettings = {
  consentPersonalData: false,
  riotName: '',
  riotTag: '',
  region: 'ap',
  apiKey: '',
  apiKeyHeader: 'Authorization',
  rateLimitTier: 'basic',
  matchCount: 5,
  cacheTTLMinutes: 30,
  lastUpdatedAt: '',
}

export const fallbackAPIStatus: APIStatus = {
  baseURL: 'https://api.henrikdev.xyz/valorant',
  consentGranted: false,
  canFetchPersonalData: false,
  rateLimitPerMinute: 30,
  cacheTTLMinutes: 30,
  safeMode: true,
  message: 'Chưa có Wails bridge hoặc chưa bật consent.',
  nextStep: 'Mở Settings, nhập consent/Riot ID/API key rồi fetch report.',
}

export const fallbackAssistantSession: AssistantSessionState = {
  active: false,
  startedAt: '',
  roundCount: 0,
  tipsShown: 0,
  lastAlertAt: '',
  message: 'Chưa có Wails bridge — chạy wails dev để dùng Live Assistant.',
  queueSize: 0,
}

export const emptyAssistantTip: AssistantTipResult = {
  hasTip: false,
  alert: { id: '', title: '', message: '', severity: 'low', source: '' },
  state: fallbackAssistantSession,
}

export const fallbackChatState: ChatState = {
  available: false,
  message: 'Wails bridge chưa sẵn sàng — chạy wails dev để dùng bot AI.',
  history: [],
}

export const defaultMapPlan = (mapId: string): MapPlan => ({
  mapId,
  title: 'Kế hoạch trước trận',
  side: 'attack',
  notes: '',
  markers: [],
  lines: [],
  updatedAt: '',
})
