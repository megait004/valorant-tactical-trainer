# US-009 Build Dashboard And Match List UI

## Status

implemented

## Lane

normal

## Product Contract

The desktop UI separates setup, status, match cache, report, and recommendation presentation into maintainable React components while preserving Wails-only data access.

## Relevant Product Docs

- `docs/product/desktop-app.md`
- `docs/product/tactical-analysis.md`
- `docs/skills/react-desktop-ui/SKILL.md`

## Acceptance Criteria

- `App.tsx` owns Wails orchestration and state only.
- Setup/consent controls are in a dedicated component.
- Header/status is in a dedicated component.
- Match cache list is in a dedicated component.
- Tactical report and recommendations are in a dedicated component.
- React still calls only Wails bindings.
- Validation commands pass.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | none |
| Integration | TypeScript compile |
| E2E | manual not run |
| Platform | `wails build` |
| Release | not required |

## Evidence

- Added `src/components/AppHeader.tsx`.
- Added `src/components/SetupPanel.tsx`.
- Added `src/components/MatchCachePanel.tsx`.
- Added `src/components/ReportPanel.tsx`.
- Added shared component prop types in `src/components/types.ts`.
- Reduced `App.tsx` to state orchestration and Wails calls.
- `pnpm typecheck`, `pnpm lint`, `go test ./...`, `pnpm build`, and `wails build` passed.
