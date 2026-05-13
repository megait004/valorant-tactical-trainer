# US-013 Add Local Settings Panel

## Status

implemented

## Lane

high-risk

## Product Contract

The desktop app exposes a local settings panel that lets users manage the locally stored Henrik API key, inspect SQLite/cache counts, clear expired cache entries, and see the local database path without React directly reading storage.

## Relevant Product Docs

- `docs/product/desktop-app.md`
- `docs/product/privacy-consent.md`
- `docs/product/api-contracts.md`

## Acceptance Criteria

- React uses Wails settings bindings only.
- Users can save or clear the optional API key through Go-managed local settings.
- Settings show data path, player count, match count, rank snapshot count, report count, cache count, and expired cache count.
- Users can clear expired cache entries without deleting valid local data.
- Reset local data keeps native confirmation behavior.

## Design Notes

- Commands: `GetSettings`, `SaveSettings`, `ClearExpiredCache`, existing `ResetAllData`.
- Queries: SQLite stats and API key presence.
- API: no provider endpoint changes.
- Tables: no new tables.
- Domain rules: frontend never receives the saved API key value, only `apiKeyConfigured`.
- UI surfaces: right-side Settings panel above match cache.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | None. |
| Integration | SQLite stats, delete setting, and clear expired cache tests. |
| E2E | Manual desktop smoke not run. |
| Platform | Wails build passes. |
| Release | Windows exe builds. |

## Harness Delta

No harness change needed.

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.
