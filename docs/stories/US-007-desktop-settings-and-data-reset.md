# US-007 Desktop Settings And Data Reset

## Status

implemented

## Lane

high-risk

## Product Contract

The app can reset all local data only after native desktop confirmation.

## Relevant Product Docs

- `docs/product/privacy-consent.md`
- `docs/product/desktop-app.md`
- `docs/decisions/0006-local-first-sqlite.md`

## Acceptance Criteria

- Reset action uses a native Wails confirmation dialog.
- Reset deletes player, consent, API cache, matches, reports, recommendations, and settings.
- UI state clears after reset.
- Reset does not call Henrik API.
- Cancel leaves local data untouched.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | no unit test yet |
| Integration | reset SQL compile/build |
| E2E | manual dialog smoke not run |
| Platform | `wails build` |
| Release | not required |

## Evidence

- Added `SettingsService.ResetAllData`.
- Added native Wails question dialog before destructive reset.
- Added `Store.ResetAll` deleting all local app tables.
- UI reset button clears local frontend state after confirmed reset.
- `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, and `wails build` passed.
