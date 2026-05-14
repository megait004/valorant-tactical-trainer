# US-017 Add Match Filter And Detail

## Status

implemented

## Lane

normal

## Product Contract

The match cache panel lets users filter stored matches and inspect a selected match inline without fetching new provider data.

## Acceptance Criteria

- Users can search stored matches by map, agent, mode, queue, or region.
- Users can filter stored matches by map.
- Users can select a match from the cached list.
- Selected match detail shows map, agent, mode/queue, region, K/D/A, rounds, headshots, damage, and match id.
- The feature uses only local cached match DTOs.
- Labels use existing English/Vietnamese translation flow.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | Not required; frontend-only local filtering. |
| Integration | TypeScript compile/build. |
| E2E | Manual click-through not run. |
| Platform | Wails build passes. |

## Evidence

- `pnpm validate` passed.
- `go test ./...` passed.
- `wails build` passed.
