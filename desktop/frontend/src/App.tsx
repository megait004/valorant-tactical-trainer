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
  const [consent, setConsent] = useState(false);
  const [lookupLoading, setLookupLoading] = useState(false);
  const [matchLoading, setMatchLoading] = useState(false);
  const [rankLoading, setRankLoading] = useState(false);
  const [reportLoading, setReportLoading] = useState(false);
  const [resetLoading, setResetLoading] = useState(false);
  const [settingsLoading, setSettingsLoading] = useState(false);
  const [assistantLoading, setAssistantLoading] = useState(false);
  const [status, setStatus] = useState('Go core waiting for binding smoke');

  useEffect(() => {
    const loadCurrentPlayer = async () => {
      try {
        const savedSettings = await GetSettings();
        setSettings(savedSettings);
        setLanguage(getLanguage(savedSettings.language));
        const player = await GetCurrentPlayer();
        if (player) {
          setCurrentPlayer(player);
          setStatus(`current player: ${player.name}#${player.tag}`);
          const savedMatches = await ListMatches(player.puuid);
          setMatches(savedMatches);
          const savedRank = await LatestRank(player.puuid);
          setRank(savedRank);
        }
      } catch (err) {
        setStatus('err: current player unavailable');
        console.error('err loading current player', err);
      }
    };

    void loadCurrentPlayer();
  }, []);

  const checkCore = async () => {
    setStatus('checking Go core...');

    try {
      const info = await AppInfo();
      setAppInfo(info);
      setStatus('data received');
    } catch (err) {
      setStatus('err: Go binding unavailable');
      console.error('err checking core', err);
    }
  };

  const saveSettings = async () => {
    setSettingsLoading(true);
    setStatus('saving settings...');

    try {
      const result = await SaveSettings({ apiKey });
      setSettings(result);
      setStatus(result.message);
      if (apiKey.trim() === '') {
        setApiKey('');
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: save settings failed');
      console.error('err saving settings', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const changeLanguage = async (nextLanguage: Language) => {
    setLanguage(nextLanguage);
    setStatus(nextLanguage === 'vi' ? 'đang lưu ngôn ngữ...' : 'saving language...');

    try {
      const result = await SaveLanguage({ language: nextLanguage });
      setSettings(result);
      setLanguage(getLanguage(result.language));
      setStatus(nextLanguage === 'vi' ? 'đã lưu ngôn ngữ' : 'language saved');
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: save language failed');
      console.error('err saving language', err);
    }
  };

  const clearExpiredCache = async () => {
    setSettingsLoading(true);
    setStatus('clearing expired cache...');

    try {
      const result = await ClearExpiredCache();
      const refreshed = await GetSettings();
      setSettings(refreshed);
      setStatus(`${result.message}: ${result.cleared} removed`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: clear cache failed');
      console.error('err clearing cache', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const exportLocalData = async () => {
    setSettingsLoading(true);
    setStatus('exporting local data...');

    try {
      const result = await ExportLocalData();
      setStatus(result.path ? `${result.message}: ${result.path}` : result.message);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: export failed');
      console.error('err exporting data', err);
    } finally {
      setSettingsLoading(false);
    }
  };

  const lookupPlayer = async () => {
    setLookupLoading(true);
    setLookupResult(null);
    setStatus('checking consent and provider...');

    try {
      const result = await LookupPlayer({ name, tag, region, consent, apiKey });
      setLookupResult(result);
      setCurrentPlayer(result.player);
      const savedMatches = await ListMatches(result.player.puuid);
      setMatches(savedMatches);
      const savedRank = await LatestRank(result.player.puuid);
      setRank(savedRank);
      setStatus(result.message || 'data received');
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: lookup failed');
      console.error('err lookup player', err);
    } finally {
      setLookupLoading(false);
    }
  };

  const refreshMatches = async () => {
    if (!currentPlayer) {
      setStatus('err: lookup player first');
      return;
    }

    setMatchLoading(true);
    setStatus('refreshing matches...');

    try {
      const result = await RefreshMatches({
        puuid: currentPlayer.puuid,
        region: currentPlayer.region || region,
        size: '10',
        apiKey,
      });
      setMatches(result.matches);
      setStatus(`${result.message}: ${result.matches.length} total, ${result.imported} touched`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: refresh matches failed');
      console.error('err refresh matches', err);
    } finally {
      setMatchLoading(false);
    }
  };

  const generateReport = async () => {
    if (!currentPlayer) {
      setStatus('err: lookup player first');
      return;
    }

    setReportLoading(true);
    setStatus('generating tactical report...');

    try {
      const nextReport = await GenerateReport(currentPlayer.puuid);
      setReport(nextReport);
      setStatus(`report generated: ${nextReport.findings.length} findings`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: report failed');
      console.error('err generating report', err);
    } finally {
      setReportLoading(false);
    }
  };

  const refreshRank = async () => {
    if (!currentPlayer) {
      setStatus('err: lookup player first');
      return;
    }

    setRankLoading(true);
    setStatus('refreshing rank...');

    try {
      const result = await RefreshRank({
        puuid: currentPlayer.puuid,
        region: currentPlayer.region || region,
        apiKey,
      });
      setRank(result.rank);
      setStatus(`${result.message}: ${result.rank.tierName || 'rank unknown'}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: refresh rank failed');
      console.error('err refresh rank', err);
    } finally {
      setRankLoading(false);
    }
  };

  const resetAllData = async () => {
    setResetLoading(true);
    setStatus('resetting local data...');

    try {
      const result = await ResetAllData();
      if (result.message === 'reset cancelled') {
        setStatus(result.message);
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
      setConsent(false);
      const refreshed = await GetSettings();
      setSettings(refreshed);
      setStatus(result.message);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: reset failed');
      console.error('err resetting data', err);
    } finally {
      setResetLoading(false);
    }
  };

  const queryAssistant = async () => {
    setAssistantLoading(true);
    setStatus('loading tactical assistant...');

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
      setStatus(`assistant ready: ${result.cards.length} cards, ${result.economyAdvice.plan}`);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: assistant failed');
      console.error('err querying assistant', err);
    } finally {
      setAssistantLoading(false);
    }
  };

  const toggleOverlay = async () => {
    const nextOverlay = !overlayEnabled;
    setStatus(nextOverlay ? 'enabling assistant overlay...' : 'disabling assistant overlay...');

    try {
      const result = await SetAssistantOverlay(nextOverlay);
      setOverlayEnabled(result.overlay);
      setStatus(result.message);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: overlay mode failed');
      console.error('err toggling overlay', err);
    }
  };

  const canLookup = consent && name.trim() !== '' && tag.trim() !== '' && !lookupLoading;
  const t = translations[language];

  return (
    <main className="min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(255,70,85,0.22),_transparent_34%),linear-gradient(135deg,_#07080d_0%,_#111827_55%,_#0e111a_100%)] text-slate-100">
      <section className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-6 py-8 sm:px-10 lg:px-12">
        <AppHeader language={language} onLanguageChange={changeLanguage} status={status} t={t} />
        <div className="grid flex-1 gap-6 py-8 xl:grid-cols-[1fr_0.9fr]">
          <SetupPanel
            appInfo={appInfo}
            apiKey={apiKey}
            canLookup={canLookup}
            consent={consent}
            currentPlayer={currentPlayer}
            lookupLoading={lookupLoading}
            lookupResult={lookupResult}
            matchLoading={matchLoading}
            name={name}
            onApiKeyChange={setApiKey}
            onCheckCore={checkCore}
            onConsentChange={setConsent}
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
