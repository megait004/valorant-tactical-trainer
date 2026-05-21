// Analysis Report: demo, last cached report, live fetch.
// FetchLiveReport bắt buộc cần Wails bridge — không có fallback (sẽ throw).

import type { AnalysisReport, LastReportResult, LiveAnalysisResult } from '../types'
import { bridge } from './bridge'
import { fallbackReport } from './fallbacks'

export const generateDemoReport = async (): Promise<AnalysisReport> =>
  bridge()?.AnalysisService?.GenerateDemoReport?.() ?? fallbackReport

export const getLastReport = async (): Promise<LastReportResult> =>
  bridge()?.AnalysisService?.GetLastReport?.() ?? {
    hasReport: false,
    result: { report: fallbackReport, source: 'fallback', cached: false, fetchedAt: '', message: '' },
  }

export const fetchLiveReport = async (): Promise<LiveAnalysisResult> => {
  const call = bridge()?.AnalysisService?.FetchLiveReport
  if (!call) {
    throw new Error('Chưa có Wails bridge; chạy app bằng wails dev/build để fetch Henrik.')
  }
  return call()
}
