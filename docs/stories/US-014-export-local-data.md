# US-014 Export Local Data

## Status

implemented

## Lane

high-risk

## Product Contract

The desktop app can export local user data to a JSON file through a native save dialog. The export is privacy-aware: it includes useful local analysis data and stats, but excludes the stored API key value and raw provider payloads.

## Relevant Product Docs

- `docs/product/desktop-app.md`
- `docs/product/privacy-consent.md`
- `docs/product/api-contracts.md`

## Acceptance Criteria

- Export is triggered through Wails/Go, not direct frontend filesystem access.
- Native save dialog is used.
- Export JSON includes metadata, data stats, players, consents, matches, rank snapshots, reports, findings, and recommendations.
- Export JSON excludes API key value and raw provider payloads.
- Settings UI exposes an export action and explains what is excluded.

## Design Notes

- Commands: `ExportLocalData` on `SettingsService`.
- Queries: read-only table snapshots from SQLite.
- API: no provider call.
- Tables: no schema changes.
- Domain rules: export is local-only and never calls Henrik.
- UI surfaces: Settings panel export button.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | None. |
| Integration | SQLite export JSON test verifies valid JSON and no API key/raw payload leak. |
| E2E | Manual native save dialog smoke not run. |
| Platform | Wails build passes. |
| Release | Windows exe builds. |

## Harness Delta

No harness change needed.

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.
