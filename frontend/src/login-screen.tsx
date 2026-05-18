import { useEffect, useState } from 'react'
import { getDataSettings, riotLogin } from './api'
import type { RiotPlayerInfo } from './types'

type LoginScreenProps = {
  onLoginSuccess: (player: RiotPlayerInfo) => void
}

const REGIONS = [
  { id: 'ap', label: 'AP — Asia Pacific (VN, SEA)' },
  { id: 'kr', label: 'KR — Korea' },
  { id: 'na', label: 'NA — North America' },
  { id: 'br', label: 'BR — Brazil' },
  { id: 'latam', label: 'LATAM — Latin America' },
  { id: 'eu', label: 'EU — Europe' },
]

const LoginScreen = ({ onLoginSuccess }: LoginScreenProps) => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [region, setRegion] = useState('ap')
  const [showRegion, setShowRegion] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    void getDataSettings().then((settings) => {
      if (settings.riotName) setUsername(settings.riotName)
      if (settings.riotTag) setPassword(settings.riotTag)
      if (settings.region) setRegion(settings.region)
    })
  }, [])

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!username.trim() || !password.trim()) {
      setError('Nhập username và password')
      return
    }
    setLoading(true)
    setError('')
    try {
      // username = Riot ID (gameName), password = Tag — verify qua Account-V1.
      // RGAPI key load từ file .env trong Go core (RIOT_API_KEY).
      const result = await riotLogin(username.trim(), password.trim().replace(/^#/, ''), region)
      if (result.error) {
        setError(result.error)
        return
      }
      if (result.success && result.playerInfo) {
        // Backend authSink đã tự lưu Riot ID + PUUID + xóa report cache cũ.
        // Frontend chỉ cần forward kết quả lên App.
        onLoginSuccess(result.playerInfo)
      } else {
        setError('Login thất bại — thử lại')
      }
    } catch (err) {
      setError(`err: ${err instanceof Error ? err.message : String(err)}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-tactical-bg px-4 py-8">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <span className="inline-block rounded-full border border-tactical-red/40 bg-tactical-red/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-tactical-red">
            Valorant Coach
          </span>
          <h1 className="mt-4 text-3xl font-black tracking-tight text-slate-100">
            Đăng nhập Riot
          </h1>
          <p className="mt-2 text-sm text-slate-400">
            Dùng tài khoản Riot Games để tiếp tục
          </p>
        </div>

        <div className="rounded-3xl border border-tactical-line bg-gradient-to-br from-tactical-panel via-[#151825] to-[#0b0d14] p-6 shadow-2xl shadow-black/30">
          <form onSubmit={(e) => void handleLogin(e)} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Username
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="riot_username"
                autoComplete="username"
                disabled={loading}
                className="rounded-xl border border-tactical-line bg-tactical-bg px-4 py-2.5 text-sm text-slate-100 placeholder-slate-600 outline-none transition focus:border-tactical-red/60 focus:ring-1 focus:ring-tactical-red/30 disabled:opacity-50"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Password
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                autoComplete="current-password"
                disabled={loading}
                className="rounded-xl border border-tactical-line bg-tactical-bg px-4 py-2.5 text-sm text-slate-100 placeholder-slate-600 outline-none transition focus:border-tactical-red/60 focus:ring-1 focus:ring-tactical-red/30 disabled:opacity-50"
              />
            </div>

            {showRegion ? (
              <div className="flex flex-col gap-1.5">
                <label className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                  Region
                </label>
                <select
                  value={region}
                  onChange={(e) => setRegion(e.target.value)}
                  disabled={loading}
                  className="rounded-xl border border-tactical-line bg-tactical-bg px-4 py-2.5 text-sm text-slate-100 outline-none transition focus:border-tactical-red/60 focus:ring-1 focus:ring-tactical-red/30 disabled:opacity-50"
                >
                  {REGIONS.map((r) => (
                    <option key={r.id} value={r.id}>
                      {r.label}
                    </option>
                  ))}
                </select>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => setShowRegion(true)}
                className="self-start text-[11px] text-slate-500 transition hover:text-slate-300"
              >
                Region: <span className="font-semibold uppercase">{region}</span> · đổi
              </button>
            )}

            {error && (
              <p className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-2.5 text-xs leading-5 text-red-400">
                {error}
              </p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="mt-1 rounded-xl bg-tactical-red px-4 py-2.5 text-sm font-semibold text-white shadow-lg shadow-tactical-red/20 transition hover:bg-tactical-red/80 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
            </button>
          </form>
        </div>

        <p className="mt-4 text-center text-xs text-slate-600">
          Thông tin đăng nhập chỉ dùng để xác thực với Riot — không lưu trữ.
        </p>
      </div>
    </main>
  )
}

export default LoginScreen
