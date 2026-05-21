// Barrel export — import từ '@/api' hoặc '../api' cho ngắn.
//
// Mỗi service backend (Go) có 1 file riêng:
//
//	bridge.ts    — Wails bridge type contract + helper `bridge()`
//	fallbacks.ts — Giá trị mặc định khi bridge chưa sẵn sàng
//	auth.ts      — Riot Account-V1 login/logout
//	settings.ts  — Data & Consent + Henrik API status
//	analysis.ts  — Demo/last/live report
//	practice.ts  — Progress checklist + session timer
//	assistant.ts — Live Assistant tip engine
//	tactical.ts  — Map catalog + map plan
//	chat.ts      — Bot AI chat (LLM coach)

export * from './auth'
export * from './settings'
export * from './analysis'
export * from './practice'
export * from './assistant'
export * from './tactical'
export * from './chat'
export { fallbackReport } from './fallbacks'
