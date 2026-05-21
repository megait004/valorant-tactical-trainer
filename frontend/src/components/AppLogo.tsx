// Logo app — nguồn gốc: logoapp.png (desktop/logoapp.png).
// Build: build/appicon.png + build/windows/icon.ico được sinh từ file đó.
// UI: frontend/public/logo.png (256px).

type AppLogoProps = {
  /** Chiều rộng/cao ô vuông (px). */
  size?: number
  className?: string
}

export const AppLogo = ({ size = 48, className = '' }: AppLogoProps) => (
  <img
    src="/logo.png"
    alt="Valorant Tactical Trainer"
    width={size}
    height={size}
    className={`rounded-2xl object-cover shadow-lg shadow-black/40 ring-1 ring-white/10 ${className}`}
    draggable={false}
  />
)
