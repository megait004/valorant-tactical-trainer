# US-006 Generate Training Recommendation List

## Status

implemented

## Lane

normal

## Product Contract

The app generates training recommendations from report findings, with reason, priority, drill, evidence, and status.

## Relevant Product Docs

- `docs/product/tactical-analysis.md`
- `docs/product/desktop-app.md`

## Acceptance Criteria

- Recommendations are generated from deterministic analysis rules.
- Each recommendation includes title, drill, priority, reason, evidence, and status.
- Recommendations are persisted to SQLite.
- UI displays recommendations beside findings.
- No recommendation is generated without evidence.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | analysis domain tests cover recommendation generation |
| Integration | recommendation persistence compile/build |
| E2E | manual not run |
| Platform | `wails build` |
| Release | not required |

## Evidence

- Added rule-based recommendations for import-data, aim discipline, trade spacing, map notes, and sample expansion.
- Added SQLite `training_recommendations` table.
- UI displays recommendation cards.
- `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, and `wails build` passed.
