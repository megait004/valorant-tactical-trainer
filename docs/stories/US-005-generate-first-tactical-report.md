# US-005 Generate First Tactical Report

## Status

implemented

## Lane

normal

## Product Contract

The app can generate an explainable tactical report from stored match summaries.

## Relevant Product Docs

- `docs/product/tactical-analysis.md`
- `docs/product/desktop-app.md`

## Acceptance Criteria

- Report generation reads stored matches for the current player.
- Report includes match count, KDA, headshot percent, average damage, top agent, and top map.
- Report creates findings with severity, confidence, description, and evidence.
- Empty match sets produce a safe recommendation to import matches.
- Report data is persisted to SQLite.
- UI displays report summary and findings.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | analysis domain tests |
| Integration | report persistence compile/build |
| E2E | manual not run |
| Platform | `wails build` |
| Release | not required |

## Evidence

- Added `desktop/internal/domain/analysis` with deterministic report rules.
- Added analysis unit tests.
- Added `AnalysisService.GenerateReport` Wails binding.
- Added SQLite report/finding/recommendation tables.
- UI can generate and display tactical report.
- `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, and `wails build` passed.
