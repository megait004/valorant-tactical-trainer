# US-011 Add Provider Throttle

## Status

implemented

## Lane

high-risk

## Product Contract

The app spaces live Henrik provider calls so normal desktop use does not burst past the Basic provider limit. Provider access remains isolated behind the Go adapter and React continues to call only Wails bindings.

## Relevant Product Docs

- `docs/product/api-contracts.md`
- `docs/product/privacy-consent.md`
- `docs/product/desktop-app.md`

## Acceptance Criteria

- Account lookup and match refresh use the throttled provider client.
- The default throttle assumes the Basic provider limit of 30 requests per minute.
- Tests prove sequential calls wait and concurrent calls reserve separate slots.
- Existing provider mock, storage, frontend validation, and Wails release build still pass.

## Design Notes

- Commands: Wails services continue exposing `LookupPlayer` and `RefreshMatches`.
- Queries: unchanged.
- API: `NewBasicClient` wraps `NewClient` with a shared `RateLimiter`.
- Tables: unchanged.
- Domain rules: unchanged.
- UI surfaces: existing status/error banner displays provider feedback.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | Rate limiter waits between sequential calls and reserves concurrent slots. |
| Integration | Provider mock tests still pass. |
| E2E | Manual live provider smoke still not run. |
| Platform | Wails build passes. |
| Release | `wails build` produces the Windows exe. |

## Harness Delta

No harness change needed.

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.
