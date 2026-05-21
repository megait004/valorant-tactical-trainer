import {
  ScreenGetAll,
  WindowSetAlwaysOnTop,
  WindowSetMaxSize,
  WindowSetMinSize,
  WindowSetPosition,
  WindowSetSize,
} from '../../../wailsjs/runtime/runtime'

const overlayWidth = 420
const overlayHeight = 200
const mainWidth = 1280
const mainHeight = 820

export const enableOverlayMode = async () => {
  WindowSetAlwaysOnTop(true)
  WindowSetMinSize(320, 160)
  WindowSetMaxSize(overlayWidth, overlayHeight)
  WindowSetSize(overlayWidth, overlayHeight)

  try {
    const screens = await ScreenGetAll()
    const screen = screens[0]
    if (screen) {
      const x = Math.max(0, screen.width - overlayWidth - 24)
      const y = Math.max(0, screen.height - overlayHeight - 80)
      WindowSetPosition(x, y)
    }
  } catch {
    WindowSetPosition(40, 40)
  }
}

export const disableOverlayMode = async () => {
  WindowSetAlwaysOnTop(false)
  WindowSetMaxSize(0, 0)
  WindowSetMinSize(960, 640)
  WindowSetSize(mainWidth, mainHeight)
  WindowSetPosition(80, 60)
}
