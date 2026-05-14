import { useEffect, useState } from 'react';
import { AppInfo } from '../wailsjs/go/main/App';
import { GenerateReport } from '../wailsjs/go/wailsiface/AnalysisService';
import { QueryAssistant } from '../wailsjs/go/wailsiface/AssistantService';
import { ListMatches, RefreshMatches } from '../wailsjs/go/wailsiface/MatchService';
import { GetCurrentPlayer, LookupPlayer } from '../wailsjs/go/wailsiface/PlayerService';
import { LatestRank, RefreshRank } from '../wailsjs/go/wailsiface/RankService';
import {
  ClearExpiredCache,
  ExportLocalData,
  GetSettings,
  ResetAllData,
  SaveLanguage,
  SaveSettings,
} from '../wailsjs/go/wailsiface/SettingsService';
import { SetAssistantOverlay } from '../wailsjs/go/wailsiface/WindowService';
import type { main, wailsiface } from '../wailsjs/go/models';
import { AppHeader } from './components/AppHeader';
import { MatchCachePanel } from './components/MatchCachePanel';
import { SetupPanel } from './components/SetupPanel';
import { SettingsPanel } from './components/SettingsPanel';
import { VirtualAssistantPanel } from './components/VirtualAssistantPanel';
import { getLanguage, translations, type Language } from './i18n';

const App = () => {
  const [appInfo, setAppInfo] = useState<main.AppInfo | null>(null);
  const [currentPlayer, setCurrentPlayer] = useState<wailsiface.PlayerDTO | null>(null);
  const [lookupResult, setLookupResult] = useState<wailsiface.LookupPlayerResult | null>(null);
  const [matches, setMatches] = useState<wailsiface.MatchDTO[]>([]);
  const [rank, setRank] = useState<wailsiface.RankDTO | null>(null);
  const [report, setReport] = useState<wailsiface.ReportDTO | null>(null);
  const [settings, setSettings] = useState<wailsiface.SettingsDTO | null>(null);
  const [assistantResult, setAssistantResult] = useState<wailsiface.AssistantResultDTO | null>(null);
  const [overlayEnabled, setOverlayEnabled] = useState(false);
  const [language, setLanguage] = useState<Language>('en');
  const [name, setName] = useState('');
  const [tag, setTag] = useState('');
  const [region, setRegion] = useState('ap');
  const [apiKey, setApiKey] = useState('');
  const [assistantMap, setAssistantMap] = useState('Ascent');
  const [assistantAgent, setAssistantAgent] = useState('');
  const [assistantSide, setAssistantSide] = useState('attack');
  const [assistantPhase, setAssistantPhase] = useState('prematch');
  const [assistantCredits, setAssistantCredits] = useState(3900);
  const [assistantOutcome, setAssistantOutcome] = useState('win');
  const [lookupLoading, setLookupLoading] = useState(false);
  const [matchLoading, setMatchLoading] = useState(false);
  const [rankLoading, setRankLoading] = useState(false);
  const [reportLoading, setReportLoading] = useState(false);
  const [resetLoading, setResetLoading] = useState(false);
  const [settingsLoading, setSettingsLoading] = useState(false);
  const [assistantLoading, setAssistantLoading] = useState(false);
  const t = translations[language];
  const [status, setStatus] = useState<string>(translations.en.coreWaiting);

  useEffect(() => {
    const loadCurrentPlayer = async () => {
      try {
        const savedSettings = await GetSettings();
        setSettings(savedSettings);
        setLanguage(getLanguage(savedSettings.language));
        const player = await GetCurrentPlayer();
        if (player) {
          setCurrentPlayer(player);
          setStatus(`${translations[getLanguage(savedSettings.language)].currentPlayer}: ${player.name}#${player.tag}`);
          const savedMatches = await ListMatches(player.puuid);
          setMatches(savedMatches);
          const savedRank = await LatestRank(player.puuid);
          setRank(savedRank);
        }
      } catch (err) {
        setStatus(t.currentPlayerUnavailable);
        console.error('err loading current player', err);
      }
    };

    void loadCurrentPlayer();
  }, []);

  const checkCore = async () => {
    setStatus(t.checkingCore);

    try {
      const info = await AppInfo();
      setAppInfo(info);
      setStatus(t.dataReceived);
    } catch (err) {
      setStatus(t.goBindingUnavailable);
      console.error('err checking core', err);
    }
  };

  const saveSettings = async () => {
    setSettingsLoading(true);
    setStatus(t.savingSettings);

    try {
      const result = await SaveSettings({ apiKey });
      setSettings(result);
      setStatus(t.settingsSaved);
      if (apiKey.trim() === '') {
        setApiKey('');
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.saveSettingsFailed);
      console.error('err saving settings', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const changeLanguage = async (nextLanguage: Language) => {
    setLanguage(nextLanguage);
    setStatus(translations[nextLanguage].savingLanguage);

    try {
      const result = await SaveLanguage({ language: nextLanguage });
      setSettings(result);
      setLanguage(getLanguage(result.language));
      setStatus(translations[nextLanguage].languageSaved);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || translations[nextLanguage].saveLanguageFailed);
      console.error('err saving language', err);
    }
  };

  const clearExpiredCache = async () => {
    setSettingsLoading(true);
    setStatus(t.clearingExpiredCache);

    try {
      const result = await ClearExpiredCache();
      const refreshed = await GetSettings();
      setSettings(refreshed);
      setStatus(`${t.clearExpiredCache}: ${result.cleared} ${t.removed}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.clearCacheFailed);
      console.error('err clearing cache', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const exportLocalData = async () => {
    setSettingsLoading(true);
    setStatus(t.exportingLocalData);

    try {
      const result = await ExportLocalData();
      setStatus(result.path ? `${t.exportJson}: ${result.path}` : t.resetCancelled);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.exportFailed);
      console.error('err exporting data', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const lookupPlayer = async () => {
    setLookupLoading(true);
    setLookupResult(null);
    setStatus(t.checkingConsentProvider);

    try {
      const result = await LookupPlayer({ name, tag, region, consent: true, apiKey });
      setLookupResult(result);
      setCurrentPlayer(result.player);
      const savedMatches = await ListMatches(result.player.puuid);
      setMatches(savedMatches);
      const savedRank = await LatestRank(result.player.puuid);
      setRank(savedRank);
      setStatus(result.message || t.dataReceived);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.lookupFailed);
      console.error('err lookup player', err);
    } finally {
      setLookupLoading(false);
    }
  };

  const refreshMatches = async () => {
    if (!currentPlayer) {
      setStatus(t.lookupFirst);
      return;
    }

    setMatchLoading(true);
    setStatus(t.refreshingMatchesStatus);

    try {
      const result = await RefreshMatches({
        puuid: currentPlayer.puuid,
        region: currentPlayer.region || region,
        size: '10',
        apiKey,
      });
      setMatches(result.matches);
      setStatus(`${t.refreshMatches}: ${result.matches.length} ${t.matches}, ${result.imported} ${t.stored}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.refreshMatchesFailed);
      console.error('err refresh matches', err);
    } finally {
      setMatchLoading(false);
    }
  };

  const generateReport = async () => {
    if (!currentPlayer) {
      setStatus(t.lookupFirst);
      return;
    }

    setReportLoading(true);
    setStatus(t.generatingReport);

    try {
      const nextReport = await GenerateReport(currentPlayer.puuid);
      setReport(nextReport);
      setStatus(`${t.reportGenerated}: ${nextReport.findings.length} ${t.findings}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.reportFailed);
      console.error('err generating report', err);
    } finally {
      setReportLoading(false);
    }
  };

  const refreshRank = async () => {
    if (!currentPlayer) {
      setStatus(t.lookupFirst);
      return;
    }

    setRankLoading(true);
    setStatus(t.refreshingRankStatus);

    try {
      const result = await RefreshRank({
        puuid: currentPlayer.puuid,
        region: currentPlayer.region || region,
        apiKey,
      });
      setRank(result.rank);
      setStatus(`${t.refreshRank}: ${result.rank.tierName || t.rankUnknown}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.refreshRankFailed);
      console.error('err refresh rank', err);
    } finally {
      setRankLoading(false);
    }
  };

  const resetAllData = async () => {
    setResetLoading(true);
    setStatus(t.resettingLocalData);

    try {
      const result = await ResetAllData();
      if (result.message === 'reset cancelled') {
        setStatus(t.resetCancelled);
        return;
      }
      setCurrentPlayer(null);
      setLookupResult(null);
      setMatches([]);
      setRank(null);
      setReport(null);
      setName('');
      setTag('');
      setApiKey('');
      const refreshed = await GetSettings();
      setSettings(refreshed);
      setStatus(t.resetLocalData);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.resetFailed);
      console.error('err resetting data', err);
    } finally {
      setResetLoading(false);
    }
  };

  const queryAssistant = async () => {
    setAssistantLoading(true);
    setStatus(t.assistantLoading);

    try {
      const result = await QueryAssistant({
        agent: assistantAgent,
        credits: assistantCredits,
        mapName: assistantMap,
        phase: assistantPhase,
        previousOutcome: assistantOutcome,
        side: assistantSide,
      });
      setAssistantResult(result);
      setStatus(`${t.assistantReady}: ${result.cards.length} cards, ${result.economyAdvice.plan}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.assistantFailed);
      console.error('err querying assistant', err);
    } finally {
      setAssistantLoading(false);
    }
  };

  const toggleOverlay = async () => {
    const nextOverlay = !overlayEnabled;
    setStatus(nextOverlay ? t.enablingOverlay : t.disablingOverlay);

    try {
      const result = await SetAssistantOverlay(nextOverlay);
      setOverlayEnabled(result.overlay);
      setStatus(result.overlay ? t.compactOverlay : t.exitOverlay);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || t.overlayFailed);
      console.error('err toggling overlay', err);
    }
  };

  const canLookup = name.trim() !== '' && tag.trim() !== '' && !lookupLoading;
  return (
    <main className="min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(255,70,85,0.22),_transparent_34%),linear-gradient(135deg,_#07080d_0%,_#111827_55%,_#0e111a_100%)] text-slate-100">
      <section className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-6 py-8 sm:px-10 lg:px-12">
        <AppHeader language={language} onLanguageChange={changeLanguage} status={status} t={t} />
        <div className="grid flex-1 gap-6 py-8 xl:grid-cols-[1fr_0.9fr]">
          <SetupPanel
            appInfo={appInfo}
            apiKey={apiKey}
            canLookup={canLookup}
            currentPlayer={currentPlayer}
            lookupLoading={lookupLoading}
            lookupResult={lookupResult}
            matchLoading={matchLoading}
            name={name}
            onApiKeyChange={setApiKey}
            onCheckCore={checkCore}
            onGenerateReport={generateReport}
            onLookupPlayer={lookupPlayer}
            onNameChange={setName}
            onRefreshMatches={refreshMatches}
            onRefreshRank={refreshRank}
            onRegionChange={setRegion}
            onResetAllData={resetAllData}
            onTagChange={setTag}
            region={region}
            rank={rank}
            rankLoading={rankLoading}
            report={report}
            reportLoading={reportLoading}
            resetLoading={resetLoading}
            tag={tag}
            t={t}
          />
          <div className="space-y-6">
            <VirtualAssistantPanel
              agent={assistantAgent}
              assistantLoading={assistantLoading}
              credits={assistantCredits}
              mapName={assistantMap}
              overlayEnabled={overlayEnabled}
              onAgentChange={setAssistantAgent}
              onCreditsChange={setAssistantCredits}
              onMapNameChange={setAssistantMap}
              onPhaseChange={setAssistantPhase}
              onPreviousOutcomeChange={setAssistantOutcome}
              onQueryAssistant={queryAssistant}
              onSideChange={setAssistantSide}
              onToggleOverlay={toggleOverlay}
              phase={assistantPhase}
              previousOutcome={assistantOutcome}
              result={assistantResult}
              side={assistantSide}
              t={t}
            />
            <SettingsPanel
              apiKey={apiKey}
              loading={settingsLoading}
              onApiKeyChange={setApiKey}
              onClearExpiredCache={clearExpiredCache}
              onExportLocalData={exportLocalData}
              onSaveSettings={saveSettings}
              settings={settings}
              t={t}
            />
            <MatchCachePanel matches={matches} t={t} />
          </div>
        </div>
      </section>
    </main>
  );
};

export default App;
