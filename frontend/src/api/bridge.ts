// Wails bridge type declaration. Mỗi service backend (Go) được Bind ra
// window.go.wailsiface.<ServiceName>.<Method>. File này chỉ là type contract;
// implement gọi nằm trong từng service file (auth.ts, analysis.ts, ...).

import type {
  APIStatus,
  AnalysisReport,
  AssistantSessionState,
  AssistantTipResult,
  ChatState,
  DataSettings,
  LastReportResult,
  LiveAnalysisResult,
  MapCatalogEntry,
  MapPlan,
  PracticeProgressState,
  PracticeSessionInput,
  PracticeSessionState,
  RiotLoginResult,
  RiotPlayerInfo,
} from '../types'

export type WailsBridge = {
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

// Helper để mỗi service gọi window.go.wailsiface ngắn gọn.
export const bridge = () => window.go?.wailsiface
