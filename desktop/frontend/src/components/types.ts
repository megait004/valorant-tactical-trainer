import type { main, wailsiface } from '../../wailsjs/go/models';
import type { Language, Translation } from '../i18n';

export type AppStatusProps = {
  language: Language;
  onLanguageChange: (value: Language) => void;
  status: string;
  t: Translation;
};

export type SetupPanelProps = {
  appInfo: main.AppInfo | null;
  apiKey: string;
  canLookup: boolean;
  consent: boolean;
  currentPlayer: wailsiface.PlayerDTO | null;
  lookupLoading: boolean;
  lookupResult: wailsiface.LookupPlayerResult | null;
  matchLoading: boolean;
  name: string;
  region: string;
  rank: wailsiface.RankDTO | null;
  rankLoading: boolean;
  report: wailsiface.ReportDTO | null;
  reportLoading: boolean;
  resetLoading: boolean;
  tag: string;
  t: Translation;
  onApiKeyChange: (value: string) => void;
  onCheckCore: () => void;
  onConsentChange: (value: boolean) => void;
  onGenerateReport: () => void;
  onLookupPlayer: () => void;
  onNameChange: (value: string) => void;
  onRefreshMatches: () => void;
  onRefreshRank: () => void;
  onRegionChange: (value: string) => void;
  onResetAllData: () => void;
  onTagChange: (value: string) => void;
};

export type MatchCachePanelProps = {
  matches: wailsiface.MatchDTO[];
  t: Translation;
};

export type SettingsPanelProps = {
  apiKey: string;
  loading: boolean;
  settings: wailsiface.SettingsDTO | null;
  t: Translation;
  onApiKeyChange: (value: string) => void;
  onClearExpiredCache: () => void;
  onExportLocalData: () => void;
  onSaveSettings: () => void;
};

export type VirtualAssistantPanelProps = {
  agent: string;
  assistantLoading: boolean;
  credits: number;
  mapName: string;
  overlayEnabled: boolean;
  phase: string;
  previousOutcome: string;
  result: wailsiface.AssistantResultDTO | null;
  side: string;
  t: Translation;
  onAgentChange: (value: string) => void;
  onCreditsChange: (value: number) => void;
  onMapNameChange: (value: string) => void;
  onPhaseChange: (value: string) => void;
  onPreviousOutcomeChange: (value: string) => void;
  onQueryAssistant: () => void;
  onSideChange: (value: string) => void;
  onToggleOverlay: () => void;
};

export type ReportPanelProps = {
  report: wailsiface.ReportDTO;
  t: Translation;
};
