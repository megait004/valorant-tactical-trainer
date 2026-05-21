// Riot Account-V1 authentication API. Backend đọc RIOT_API_KEY từ .env;
// frontend chỉ pass Riot ID + tag + region.

import type { RiotLoginResult, RiotPlayerInfo } from '../types'
import { bridge } from './bridge'

export const riotLogin = async (
  riotID: string,
  tagLine: string,
  region: string,
): Promise<RiotLoginResult> => {
  const login = bridge()?.AuthService?.Login
  if (!login) {
    return { success: false, error: 'Wails bridge chưa sẵn sàng — chạy bằng wails dev/build' }
  }
  return login(riotID, tagLine, region)
}

export const riotGetPlayerInfo = async (): Promise<RiotPlayerInfo | null> =>
  bridge()?.AuthService?.GetPlayerInfo?.() ?? null

export const riotIsLoggedIn = async (): Promise<boolean> =>
  bridge()?.AuthService?.IsLoggedIn?.() ?? false

export const riotLogout = async (): Promise<void> => {
  const logout = bridge()?.AuthService?.Logout
  if (logout) {
    await logout()
  }
}
