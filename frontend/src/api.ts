import { fallbackMapCatalog } from './map-catalog-fallback'
import type {
  APIStatus,
  AnalysisReport,
  AssistantSessionState,
  AssistantTipResult,
  ChatState,
  DataSettings,
  MapCatalogEntry,
  MapPlan,
  LastReportResult,
  LiveAnalysisResult,
  PracticeProgressState,
  PracticeSessionInput,
  PracticeSessionState,
  RiotLoginResult,
  RiotPlayerInfo,
} from './types'

type WailsBridge = {
  wailsiface?: {
    AuthService?: {
      Login?: (riotID: string, tagLine: string, region: string) => Promise<RiotLoginResult>
      GetPlayerInfo?: () => Promise<RiotPlayerInfo | null>
      IsLoggedIn?: () => Promise<boolean>
      Logout?: () => Promise<void>
    }
    ChatService?: {
      IsAvailable?: () => Promise<boolean>
      GetState?: () => Promise<ChatState>
      SendMessage?: (message: string) => Promise<ChatState>
      Reset?: () => Promise<ChatState>
    }
    SettingsService?: {
      GetDataSettings?: () => Promise<DataSettings>
      SaveDataSettings?: (settings: DataSettings) => Promise<DataSettings>
      GetAPIStatus?: () => Promise<APIStatus>
    }
    AnalysisService?: {
      GenerateDemoReport?: () => Promise<AnalysisReport>
      GetLastReport?: () => Promise<LastReportResult>
      FetchLiveReport?: () => Promise<LiveAnalysisResult>
    }
    PracticeService?: {
      GetPracticeProgress?: () => Promise<PracticeProgressState>
      SetPracticeProgress?: (itemID: string, done: boolean) => Promise<PracticeProgressState>
      ResetPracticeProgress?: () => Promise<PracticeProgressState>
      GetPracticeSessions?: () => Promise<PracticeSessionState>
      FinishPracticeSession?: (input: PracticeSessionInput) => Promise<PracticeSessionState>
    }
    AssistantService?: {
      GetSessionState?: () => Promise<AssistantSessionState>
      StartSession?: () => Promise<AssistantSessionState>
      StopSession?: () => Promise<AssistantSessionState>
      RequestTip?: () => Promise<AssistantTipResult>
      MarkRoundStart?: () => Promise<AssistantTipResult>
      PollAutoTip?: () => Promise<AssistantTipResult>
    }
    TacticalService?: {
      ListMaps?: () => Promise<MapCatalogEntry[]>
      LoadMapPlan?: (mapId: string) => Promise<MapPlan>
      SaveMapPlan?: (plan: MapPlan) => Promise<MapPlan>
      DeleteMapPlan?: (mapId: string) => Promise<void>
    }
  }
}

declare global {
  interface Window {
    go?: WailsBridge
  }
}

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

export const getDataSettings = async () => window.go?.wailsiface?.SettingsService?.GetDataSettings?.() ?? fallbackDataSettings

export const saveDataSettings = async (settings: DataSettings) =>
  window.go?.wailsiface?.SettingsService?.SaveDataSettings?.(settings) ?? { ...settings, lastUpdatedAt: new Date().toISOString() }

export const getAPIStatus = async () => window.go?.wailsiface?.SettingsService?.GetAPIStatus?.() ?? fallbackAPIStatus

export const generateDemoReport = async () => window.go?.wailsiface?.AnalysisService?.GenerateDemoReport?.() ?? fallbackReport

export const getLastReport = async () => window.go?.wailsiface?.AnalysisService?.GetLastReport?.() ?? { hasReport: false, result: { report: fallbackReport, source: 'fallback', cached: false, fetchedAt: '', message: '' } }

export const fetchLiveReport = async () => {
  if (!window.go?.wailsiface?.AnalysisService?.FetchLiveReport) {
    throw new Error('Chưa có Wails bridge; chạy app bằng wails dev/build để fetch Henrik.')
  }
  return window.go.wailsiface.AnalysisService.FetchLiveReport()
}

export const getPracticeProgress = async () => window.go?.wailsiface?.PracticeService?.GetPracticeProgress?.() ?? { items: {}, updatedAt: '' }

export const setPracticeProgress = async (itemID: string, done: boolean) =>
  window.go?.wailsiface?.PracticeService?.SetPracticeProgress?.(itemID, done) ?? { items: { [itemID]: done }, updatedAt: new Date().toISOString() }

export const resetPracticeProgress = async () => window.go?.wailsiface?.PracticeService?.ResetPracticeProgress?.() ?? { items: {}, updatedAt: new Date().toISOString() }

export const getPracticeSessions = async () => window.go?.wailsiface?.PracticeService?.GetPracticeSessions?.() ?? { sessions: [], updatedAt: '' }

export const finishPracticeSession = async (input: PracticeSessionInput) =>
  window.go?.wailsiface?.PracticeService?.FinishPracticeSession?.(input) ?? { sessions: [{ id: `session-${Date.now()}`, finishedAt: new Date().toISOString(), ...input }], updatedAt: new Date().toISOString() }

const fallbackAssistantSession: AssistantSessionState = {
  active: false,
  startedAt: '',
  roundCount: 0,
  tipsShown: 0,
  lastAlertAt: '',
  message: 'Chưa có Wails bridge — chạy wails dev để dùng Live Assistant.',
  queueSize: 0,
}

const emptyTip: AssistantTipResult = {
  hasTip: false,
  alert: { id: '', title: '', message: '', severity: 'low', source: '' },
  state: fallbackAssistantSession,
}

export const getAssistantSession = async () =>
  window.go?.wailsiface?.AssistantService?.GetSessionState?.() ?? fallbackAssistantSession

export const startAssistantSession = async () =>
  window.go?.wailsiface?.AssistantService?.StartSession?.() ?? {
    ...fallbackAssistantSession,
    active: true,
    message: 'Demo: chạy wails dev để rule engine thật.',
    currentAlert: {
      id: 'demo',
      title: 'Demo Live Assistant',
      message: 'Utility trước peek — fetch report ở Coach để gợi ý cá nhân.',
      severity: 'medium',
      source: 'demo',
    },
    queueSize: 1,
    tipsShown: 1,
    startedAt: new Date().toISOString(),
  }

export const stopAssistantSession = async () =>
  window.go?.wailsiface?.AssistantService?.StopSession?.() ?? { ...fallbackAssistantSession, message: 'Đã tắt (demo).' }

export const requestAssistantTip = async () =>
  window.go?.wailsiface?.AssistantService?.RequestTip?.() ?? emptyTip

export const markAssistantRoundStart = async () =>
  window.go?.wailsiface?.AssistantService?.MarkRoundStart?.() ?? emptyTip

export const pollAssistantAutoTip = async () =>
  window.go?.wailsiface?.AssistantService?.PollAutoTip?.() ?? emptyTip

const defaultMapPlan = (mapId: string): MapPlan => ({
  mapId,
  title: 'Kế hoạch trước trận',
  side: 'attack',
  notes: '',
  markers: [],
  lines: [],
  updatedAt: '',
})

export const listTacticalMaps = async () =>
  window.go?.wailsiface?.TacticalService?.ListMaps?.() ?? fallbackMapCatalog

export const loadMapPlan = async (mapId: string) =>
  window.go?.wailsiface?.TacticalService?.LoadMapPlan?.(mapId) ?? defaultMapPlan(mapId)

export const saveMapPlan = async (plan: MapPlan) =>
  window.go?.wailsiface?.TacticalService?.SaveMapPlan?.(plan) ?? { ...plan, updatedAt: new Date().toISOString() }

export const deleteMapPlan = async (mapId: string) => {
  if (window.go?.wailsiface?.TacticalService?.DeleteMapPlan) {
    await window.go.wailsiface.TacticalService.DeleteMapPlan(mapId)
  }
}

// --- Riot Auth (Account-V1, key load từ .env trong Go core) ---

export const riotLogin = async (
  riotID: string,
  tagLine: string,
  region: string,
): Promise<RiotLoginResult> => {
  if (!window.go?.wailsiface?.AuthService?.Login) {
    return { success: false, error: 'Wails bridge chưa sẵn sàng — chạy bằng wails dev/build' }
  }
  return window.go.wailsiface.AuthService.Login(riotID, tagLine, region)
}

export const riotGetPlayerInfo = async (): Promise<RiotPlayerInfo | null> =>
  window.go?.wailsiface?.AuthService?.GetPlayerInfo?.() ?? null

export const riotIsLoggedIn = async (): Promise<boolean> =>
  window.go?.wailsiface?.AuthService?.IsLoggedIn?.() ?? false

export const riotLogout = async (): Promise<void> => {
  if (window.go?.wailsiface?.AuthService?.Logout) {
    await window.go.wailsiface.AuthService.Logout()
  }
}

// --- Chat (LLM coach) ---

const fallbackChatState: ChatState = {
  available: false,
  message: 'Wails bridge chưa sẵn sàng — chạy wails dev để dùng bot AI.',
  history: [],
}

export const chatIsAvailable = async (): Promise<boolean> =>
  window.go?.wailsiface?.ChatService?.IsAvailable?.() ?? false

export const chatGetState = async (): Promise<ChatState> =>
  window.go?.wailsiface?.ChatService?.GetState?.() ?? fallbackChatState

export const chatSendMessage = async (message: string): Promise<ChatState> => {
  if (!window.go?.wailsiface?.ChatService?.SendMessage) {
    throw new Error('Bot AI chưa sẵn sàng — chạy bằng wails dev/build')
  }
  return window.go.wailsiface.ChatService.SendMessage(message)
}

export const chatReset = async (): Promise<ChatState> =>
  window.go?.wailsiface?.ChatService?.Reset?.() ?? fallbackChatState
