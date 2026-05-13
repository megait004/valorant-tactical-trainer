# Desktop App Contract

## Stack

- Wails v2 stable.
- Go backend/core.
- React and TypeScript frontend.
- Tailwind CSS UI.
- SQLite local storage.

## App Root

The Wails app lives in `desktop/` because the repository root also stores downloaded reference docs and project planning docs.

Important paths:

- `desktop/app.go`
- `desktop/main.go`
- `desktop/wails.json`
- `desktop/frontend/`
- future Go layers under `desktop/internal/`

## Boundary Rule

React must call Go through Wails-generated bindings.

React must not:

- Call Henrik API directly.
- Read or write SQLite directly.
- Own domain analysis logic.
- Store API key outside local Go-managed settings flow.

Go must own:

- API calls.
- Rate limits.
- Cache.
- Consent enforcement.
- Storage.
- Analysis.
- Recommendation generation.

## Screens

Setup screen:

- Consent text.
- Riot name/tag inputs.
- Region/affinity selection.
- Optional API key entry or link to settings.
- Lookup action.

Dashboard screen:

- Current player.
- Import status.
- Match count.
- Latest MMR when available.
- Latest rank/MMR card with tier, RR, elo, last-game delta, region, and fetch time.
- Top findings.
- Latest recommendations.
- Rate-limit status.

Match list screen:

- Imported matches.
- Filters by map, mode, agent, date.
- Open match detail.

Match detail screen:

- Match metadata.
- Teams and player stats.
- Round timeline.
- Economy summary.
- Plant and defuse events.

Tactical report screen:

- Summary.
- Findings.
- Evidence.
- Metric charts.
- Linked recommendations.

Training plan screen:

- Recommendation list.
- Priority and status.
- Drill details.
- Evidence link.

Settings screen:

- API key.
- Region.
- Cache info.
- Export data.
- Reset data.
- Consent/data explanation.

## Wails Services

Initial services:

- `PlayerService`
- `MatchService`
- `AnalysisService`
- `SettingsService`

Service methods should return UI-safe DTOs and errors that can be shown naturally.

## Native Desktop Behavior

- Use native confirmation dialog before destructive reset.
- Use local file dialogs for export/import when implemented.
- Build must pass `wails build` before release/demo.

## UI Direction

- Dark tactical dashboard.
- Valorant-inspired red accents without copying official assets too closely.
- Dense but readable desktop layout.
- Responsive enough for narrower windows.
