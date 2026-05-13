import type { main, wailsiface } from '../../wailsjs/go/models';

export type AppStatusProps = {
  status: string;
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
};

export type ReportPanelProps = {
  report: wailsiface.ReportDTO;
};
