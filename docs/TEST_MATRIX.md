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
| US-003 | Consent required before account lookup | player normalization tests | provider adapter mock tests and storage temp DB tests | manual not run | `wails build` passed | implemented | `pnpm validate`, `go test ./...`, `wails build` passed |
| US-004 | Match history can be fetched, cached, and deduped | thin domain proof only | mock provider matches test and SQLite match/cache tests | manual not run | `wails build` passed | implemented | `pnpm validate`, `go test ./...`, `wails build` passed |
| US-005 | Tactical report generated from stored matches | analysis domain tests | report persistence compile/build | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, `wails build` passed |
| US-006 | Training recommendations are generated with evidence | analysis domain tests | recommendation persistence compile/build | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm typecheck`, `pnpm lint`, `pnpm build`, `wails build` passed |
| US-007 | Settings and local data reset work safely | no unit test yet | SQLite reset temp DB test | manual dialog smoke not run | `wails build` passed | implemented | native Wails confirmation added; `pnpm validate`, `go test ./...`, `wails build` passed |
| US-008 | SQLite schema persists players, consent, matches, reports, and cache | analysis/player tests | SQLite temp DB tests for player, cache, matches, reports, reset | no | `wails build` passed | implemented | SQLite tables and tests added |
| US-009 | Dashboard and match list render imported data | no | TypeScript compile/build | manual not run | `wails build` passed | implemented | `pnpm typecheck`, `pnpm lint`, `go test ./...`, `pnpm build`, `wails build` passed |
| US-010 | Validation and release build flow is documented and runnable | analysis/player tests | mock provider and SQLite temp DB tests | manual not run | `wails build` passed | implemented | `desktop/VALIDATION.md`, `pnpm validate`, `go test ./...`, `wails build` passed |
| US-011 | Live provider calls are throttled for Basic provider limits | rate limiter sequential and concurrent tests | provider mock tests still pass | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-012 | Latest MMR/rank is fetched, cached, persisted, displayed, and reset | no | provider mock MMR test and SQLite rank snapshot/reset tests | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-013 | Local settings panel manages API key presence, data stats, and expired cache cleanup | no | SQLite stats/delete setting/clear expired cache tests | manual not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-014 | Local data can be exported to privacy-aware JSON | no | SQLite export JSON test excludes API key and raw payloads | manual save dialog not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-015 | Virtual Tactical Assistant provides safe local tactical cards and economy advice | economy tests and seed-card completeness tests | SQLite tactical card query test | manual UI smoke not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-016 | Core UI supports English/Vietnamese and docs explain usage/code structure | existing tests | Wails bindings compile for `SaveLanguage` | manual language smoke not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-017 | Cached matches can be searched, filtered, selected, and inspected inline | no | TypeScript compile/build | manual UI smoke not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |
| US-018 | VTA supports compact always-on-top overlay-like mode safely | no | Wails bindings compile for `WindowService` | manual overlay click-through not run | `wails build` passed | implemented | `go test ./...`, `pnpm validate`, `wails build` passed |

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
