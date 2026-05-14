# US-015 Add Virtual Tactical Assistant MVP

## Status

implemented

## Lane

high-risk

## Product Contract

The desktop app provides a safe Virtual Tactical Assistant MVP. Users manually choose map, agent, side, phase, credits, and previous-round outcome. The app returns local tactical cards and economy advice without reading VALORANT memory, injecting into the game, or calling external APIs.

## Relevant Product Docs

- `docs/product/desktop-app.md`
- `docs/product/privacy-consent.md`
- `docs/product/tactical-analysis.md`

## Acceptance Criteria

- React calls Go through Wails-generated bindings only.
- Tactical card data is seeded locally into SQLite.
- Assistant can return pre-match composition/default-strat cards.
- Assistant can return in-game lookup cards for crosshair, lineups, and default strats.
- Economy manager returns Eco, Light/Half Buy, Force Buy, or Full Buy advice based on credits and previous outcome.
- UI states clearly that MVP does not read memory or inject into VALORANT.
- Overlay/hotkey/live presence is not implemented in this MVP.

## Design Notes

- Domain: `internal/domain/assistant`.
- Storage table: `tactical_cards`.
- Wails service: `AssistantService.QueryAssistant`.
- UI: `VirtualAssistantPanel`.
- Seed content covers initial Ascent, Bind, Haven, Sova, Viper, and Brimstone examples.

## Anti-Cheat Safety

- No memory reading.
- No process injection.
- No screen automation.
- No hidden live game state detection.
- Local-only manual lookup in this MVP.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | Economy decision tests and seed card completeness tests. |
| Integration | SQLite seed/query test returns priority-ordered cards and agent-specific lineup. |
| E2E | Manual UI smoke not run. |
| Platform | Wails build passes. |
| Release | Windows exe builds. |

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.

## Follow-Up Stories

- Add compact assistant window / always-on-top mode with a visible quick close control.
- Add configurable hotkeys after UX and anti-cheat review.
- Add local Riot presence only if a documented safe local endpoint is used and user explicitly enables it.
- Expand tactical cards with owned/original images or short clips.
