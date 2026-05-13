import { useEffect, useState } from 'react';
import { AppInfo } from '../wailsjs/go/main/App';
import { GenerateReport } from '../wailsjs/go/wailsiface/AnalysisService';
import { ListMatches, RefreshMatches } from '../wailsjs/go/wailsiface/MatchService';
import { GetCurrentPlayer, LookupPlayer } from '../wailsjs/go/wailsiface/PlayerService';
import { ResetAllData } from '../wailsjs/go/wailsiface/SettingsService';
import type { main, wailsiface } from '../wailsjs/go/models';
import { AppHeader } from './components/AppHeader';
import { MatchCachePanel } from './components/MatchCachePanel';
import { SetupPanel } from './components/SetupPanel';

const App = () => {
  const [appInfo, setAppInfo] = useState<main.AppInfo | null>(null);
  const [currentPlayer, setCurrentPlayer] = useState<wailsiface.PlayerDTO | null>(null);
  const [lookupResult, setLookupResult] = useState<wailsiface.LookupPlayerResult | null>(null);
  const [matches, setMatches] = useState<wailsiface.MatchDTO[]>([]);
  const [report, setReport] = useState<wailsiface.ReportDTO | null>(null);
  const [name, setName] = useState('');
  const [tag, setTag] = useState('');
  const [region, setRegion] = useState('ap');
  const [apiKey, setApiKey] = useState('');
  const [consent, setConsent] = useState(false);
  const [lookupLoading, setLookupLoading] = useState(false);
  const [matchLoading, setMatchLoading] = useState(false);
  const [reportLoading, setReportLoading] = useState(false);
  const [resetLoading, setResetLoading] = useState(false);
  const [status, setStatus] = useState('Go core waiting for binding smoke');

  useEffect(() => {
    const loadCurrentPlayer = async () => {
      try {
        const player = await GetCurrentPlayer();
        if (player) {
          setCurrentPlayer(player);
          setStatus(`current player: ${player.name}#${player.tag}`);
          const savedMatches = await ListMatches(player.puuid);
          setMatches(savedMatches);
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
      setReport(null);
      setName('');
      setTag('');
      setApiKey('');
      setConsent(false);
      setStatus(result.message);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setStatus(message || 'err: reset failed');
      console.error('err resetting data', err);
    } finally {
      setResetLoading(false);
    }
  };

  const canLookup = consent && name.trim() !== '' && tag.trim() !== '' && !lookupLoading;

  return (
    <main className="min-h-screen overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(255,70,85,0.22),_transparent_34%),linear-gradient(135deg,_#07080d_0%,_#111827_55%,_#0e111a_100%)] text-slate-100">
      <section className="mx-auto flex min-h-screen w-full max-w-7xl flex-col px-6 py-8 sm:px-10 lg:px-12">
        <AppHeader status={status} />
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
            onRegionChange={setRegion}
            onResetAllData={resetAllData}
            onTagChange={setTag}
            region={region}
            report={report}
            reportLoading={reportLoading}
            resetLoading={resetLoading}
            tag={tag}
          />
          <MatchCachePanel matches={matches} />
        </div>
      </section>
    </main>
  );
};

export default App;
