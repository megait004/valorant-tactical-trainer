# US-018 Add Safe Assistant Overlay Mode

## Status

implemented

## Lane

high-risk

## Product Contract

The Virtual Tactical Assistant can switch the app window into a compact always-on-top assistant mode. This gives a practical overlay-like workflow without injecting into VALORANT or reading game memory.

## How It Works

- User opens the VTA panel.
- User manually selects map, agent, side, phase, credits, and previous outcome.
- User clicks `Get tactical cards`.
- App queries local `tactical_cards` in SQLite and computes economy advice.
- User clicks `Compact overlay`.
- Go `WindowService` calls Wails runtime APIs:
  - `WindowSetAlwaysOnTop(true)`
  - `WindowSetSize(520, 760)`
  - `WindowSetTitle("VTA Overlay - Valorant Tactical Trainer")`
- User clicks `Exit overlay` to restore normal size/title and disable always-on-top.

## Anti-Cheat Safety

- No process memory reading.
- No game process injection.
- No click/aim automation.
- No hidden live game state detection.
- No transparent game hook.
- Overlay mode is only a Wails window state change.

## Acceptance Criteria

- VTA panel includes a clear safety notice.
- VTA panel includes `Compact overlay` and `Exit overlay` actions.
- Overlay mode sets the app window always-on-top and compact size.
- Exit restores app window size/title and disables always-on-top.
- Wails binding compiles and build passes.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | Not required; Wails runtime wrapper. |
| Integration | Wails bindings generated for `WindowService`. |
| E2E | Manual click-through not run. |
| Platform | Wails build passes. |

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.
