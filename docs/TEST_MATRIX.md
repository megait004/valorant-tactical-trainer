# Test Matrix

This file maps product behavior to proof.

## Status Values

| Status | Meaning |
| --- | --- |
| planned | Accepted as intended behavior, not implemented. |
| in_progress | Actively being built. |
| implemented | Implemented and proof exists. |
| changed | Contract changed after earlier implementation. |
| retired | No longer part of the product contract. |

## Matrix

| Story | Contract | Unit | Integration | E2E | Platform | Status | Evidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| US-001 | Runnable Wails React Go shell in `desktop/` | no | binding smoke | no | passed | implemented | `pnpm typecheck`, `pnpm lint`, `go test ./...`, `pnpm build`, `wails build` passed |
| US-002 | Product docs, decisions, skills, backlog, and test matrix exist | no | no | no | no | implemented | docs created |
| US-003 | Consent required before account lookup | player normalization tests | provider adapter/storage implemented, no mock tests yet | manual not run | `wails build` passed | implemented | `pnpm typecheck`, `pnpm lint`, `go test ./...`, `pnpm build`, `wails build` passed |
| US-004 | Match history can be fetched, cached, and deduped | thin domain proof only | provider adapter/storage compile, no mock tests yet | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, `wails build` passed |
| US-005 | Tactical report generated from stored matches | analysis domain tests | report persistence compile/build | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, `wails build` passed |
| US-006 | Training recommendations are generated with evidence | analysis domain tests | recommendation persistence compile/build | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, `wails build` passed |
| US-007 | Settings and local data reset work safely | no unit test yet | reset SQL compile/build | manual dialog smoke not run | `wails build` passed | implemented | native Wails confirmation added; validation commands passed |
| US-008 | SQLite schema persists players, consent, matches, reports, and cache | analysis/player tests | migrations compile/build | no | `wails build` passed | implemented | SQLite tables added for player, consent, cache, matches, reports, findings, recommendations |
| US-009 | Dashboard and match list render imported data | no | TypeScript compile/build | manual not run | `wails build` passed | implemented | `pnpm typecheck`, `pnpm lint`, `go test ./...`, `pnpm build`, `wails build` passed |
| US-010 | Validation and release build flow is documented and runnable | no | no | no | planned | planned | none |

## Evidence Rules

- Unit proof covers pure domain and application rules.
- Integration proof covers API adapter, storage, provider mock, and service contracts.
- E2E proof covers user-visible flows.
- Platform proof covers Wails desktop behavior and release builds.
- A story can be implemented without every proof column if the story packet explains why.

## Expected Validation Ladder

```text
pnpm lint
pnpm typecheck
pnpm test
go test ./...
wails build
manual desktop smoke
```

Do not mark proof as passed until commands exist and have been run.
