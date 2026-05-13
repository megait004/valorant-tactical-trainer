# US-012 Fetch And Display Current Rank

## Status

implemented

## Lane

high-risk

## Product Contract

For a consenting current player, the desktop app can fetch latest competitive MMR/rank from Henrik, cache the result locally, persist rank snapshots in SQLite, and display the latest rank in the setup panel.

## Relevant Product Docs

- `docs/product/api-contracts.md`
- `docs/product/desktop-app.md`
- `docs/product/privacy-consent.md`

## Acceptance Criteria

- The provider adapter exposes MMR by PUUID through Go only.
- Rank refresh uses the throttled Basic provider client and stored API key fallback.
- Latest rank is saved as a local SQLite snapshot and can be loaded on app start.
- Reset local data removes rank snapshots.
- React displays latest tier, RR, elo, last-game MMR delta, region, and fetch time.

## Design Notes

- Commands: `RefreshRank` and `LatestRank` Wails bindings.
- Queries: latest rank by player PUUID.
- API: `MMRByPUUID` calls `/valorant/v2/by-puuid/mmr/{region}/{puuid}`.
- Tables: `rank_snapshots`.
- Domain rules: rank snapshot is a provider-derived fact, not an analysis result.
- UI surfaces: setup panel player card gains a rank card and refresh action.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | None beyond DTO mapping. |
| Integration | Provider mock MMR test and SQLite rank snapshot/reset tests. |
| E2E | Manual live Henrik smoke not run. |
| Platform | Wails build passes. |
| Release | Windows exe builds. |

## Harness Delta

No harness change needed.

## Evidence

- `go test ./...` passed.
- `pnpm validate` passed.
- `wails build` passed.
