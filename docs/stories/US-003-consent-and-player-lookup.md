# US-003 Consent And Player Lookup

## Status

implemented

## Lane

high-risk

## Product Contract

The app must require explicit consent before fetching personal VALORANT account data from Henrik API.

## Relevant Product Docs

- `docs/product/privacy-consent.md`
- `docs/product/api-contracts.md`
- `docs/product/desktop-app.md`
- `docs/decisions/0005-valorant-api-adapter.md`
- `docs/decisions/0006-local-first-sqlite.md`
- `docs/skills/valorant-api-client/SKILL.md`
- `docs/skills/wails-desktop-app/SKILL.md`
- `docs/skills/react-desktop-ui/SKILL.md`

## Acceptance Criteria

- User cannot call lookup without checking consent.
- Consent copy explains provider, data fetched, local storage, and reset path.
- Account lookup resolves name/tag to account data and PUUID.
- Consent and player profile are stored locally.
- Provider errors are normalized and shown safely.
- No Riot credentials are requested.
- React does not call Henrik API directly.

## Design Notes

- Commands: `ConfirmConsent`, `LookupPlayer`, `GetCurrentPlayer`.
- Queries: current player, consent status.
- API: Henrik account by name/tag through Go adapter.
- Tables: `players`, `consents`, `settings` if API key/region is included.
- Domain rules: consent gate must run before provider call.
- UI surfaces: setup screen and initial dashboard transition.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | consent gate rules |
| Integration | mock provider account lookup and error normalization |
| E2E | setup consent -> lookup flow |
| Platform | desktop launch and lookup smoke |
| Release | not required until release story |

## Harness Delta

High-risk story because it touches privacy, external provider behavior, API contracts, and local data.

## Evidence

- Added Go player domain normalization and consent version.
- Added Henrik account lookup adapter under `desktop/internal/infrastructure/valorantapi`.
- Added SQLite local store under `desktop/internal/infrastructure/storage` with `players`, `consents`, and `settings` tables.
- Added Wails `PlayerService` with `GetCurrentPlayer` and `LookupPlayer`.
- React setup UI requires consent before enabling lookup.
- React calls only Wails bindings for player lookup.
- `wails generate module`: passed.
- `pnpm typecheck`: passed.
- `pnpm lint`: passed.
- `go test ./...`: passed, no test files yet.
- `pnpm build`: passed.
- `wails build`: passed and created `desktop/build/bin/valorant-tactical-trainer.exe`.

Manual live provider lookup and desktop click-through smoke were not run in this turn.
