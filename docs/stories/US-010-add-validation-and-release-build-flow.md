# US-010 Add Validation And Release Build Flow

## Status

implemented

## Lane

normal

## Product Contract

The project has a repeatable validation flow for frontend, Go, and Wails desktop build checks.

## Relevant Product Docs

- `docs/TEST_MATRIX.md`
- `docs/skills/validation-release/SKILL.md`
- `desktop/VALIDATION.md`

## Acceptance Criteria

- Frontend has a validation script.
- Desktop validation docs describe the required command sequence.
- Go tests include provider adapter mock coverage.
- Go tests include SQLite storage coverage.
- Wails build remains part of platform proof.
- Test matrix reflects stronger proof.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | analysis/player tests |
| Integration | mock provider and SQLite temp DB tests |
| E2E | manual not run |
| Platform | `wails build` |
| Release | validation doc exists |

## Evidence

- Added `pnpm validate` script in `desktop/frontend/package.json`.
- Added `desktop/VALIDATION.md`.
- Added mock Henrik API tests for account lookup, rate limit, and matches by PUUID.
- Added SQLite temp DB tests for player/consent, match dedupe, API cache expiry, report save, and reset.
- `pnpm lint`, `pnpm typecheck`, `pnpm build`, `go test ./...`, and `wails build` passed.
