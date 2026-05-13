# US-004 Fetch And Cache Match History

## Status

implemented

## Lane

high-risk

## Product Contract

The app can fetch recent match history for a consenting current player through the Go Henrik API adapter, cache provider responses, persist match summaries locally, and list stored matches in the desktop UI.

## Relevant Product Docs

- `docs/product/api-contracts.md`
- `docs/product/privacy-consent.md`
- `docs/product/desktop-app.md`
- `docs/product/tactical-analysis.md`
- `docs/decisions/0005-valorant-api-adapter.md`
- `docs/decisions/0006-local-first-sqlite.md`
- `docs/skills/valorant-api-client/SKILL.md`
- `docs/skills/wails-desktop-app/SKILL.md`
- `docs/skills/react-desktop-ui/SKILL.md`

## Acceptance Criteria

- Match refresh requires a current player PUUID.
- React calls Wails `MatchService`, not Henrik API directly.
- Go adapter fetches matches by PUUID from Henrik API v3 endpoint shape.
- Provider response is cached locally with a short TTL.
- Match summaries are persisted in SQLite and deduped by match id plus player PUUID.
- UI can list stored matches and show basic KDA/map/agent data.
- Validation commands pass.

## Design Notes

- Commands: `RefreshMatches`.
- Queries: `ListMatches`.
- API: Henrik `v3/by-puuid/matches/{region}/{puuid}?size=...`.
- Tables: `api_cache`, `matches`.
- Domain rules: match summary owns player PUUID and raw payload remains outside UI DTO.
- UI surfaces: refresh matches button and match cache panel.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | player normalization tests added while working this slice |
| Integration | provider adapter and SQLite paths compile; live provider not called in validation |
| E2E | manual desktop match refresh not run |
| Platform | `wails build` |
| Release | not required until release story |

## Harness Delta

Added story evidence and updated test matrix. Unit proof is still thin; future API/cache work should add mock provider tests.

## Evidence

- Added match domain summary model.
- Added `MatchesByPUUID` to Henrik API adapter.
- Added SQLite `api_cache` and `matches` tables.
- Added storage save/list/cache methods.
- Added Wails `MatchService` with `RefreshMatches` and `ListMatches`.
- React UI can refresh and display stored match summaries.
- `go test ./...`: passed, including new player domain tests.
- `pnpm typecheck`: passed.
- `pnpm lint`: passed.
- `pnpm build`: passed.
- `wails build`: passed and created `desktop/build/bin/valorant-tactical-trainer.exe`.

Manual live provider refresh was not run in this turn.
