# Valorant Tactical Trainer - Project Plan

Date: 2026-05-13

## 1. Project Summary

Valorant Tactical Trainer is a desktop application for training and tactical analysis in VALORANT.

The app uses:

- Wails v2 stable for the desktop shell.
- Go as the core backend/application layer.
- React and TypeScript for the frontend.
- Tailwind CSS for styling.
- SQLite for local-first storage.
- Henrik unofficial VALORANT API as the external data provider.

The product helps players and coaches analyze match history, identify tactical weaknesses, and generate practical training recommendations based on match evidence.

## 2. Final Stack Decisions

Use Wails v2 stable for the desktop framework because it is stable enough for DATN/demo/release and supports Go backend plus web frontend in one desktop binary.

Use React, TypeScript, and Tailwind CSS for frontend UI. React owns presentation only and must not call Henrik API directly.

Use Go for API integration, rate limiting, cache, SQLite persistence, domain logic, tactical analysis, recommendations, and Wails-bound application services.

Use SQLite for local-first storage because the app does not need an external backend in MVP and analysis requires structured local queries.

Store project-specific skills in `docs/skills/` so agent workflows stay with product docs and can later be mirrored to tool-specific skill folders if needed.

## 3. Source Material

Local docs used for planning:

- `tailieuapivalorant.md`
- `go-valorant-api/`
- `wails/`
- `react/`
- `harness-experimental/docs/`

Key findings:

- Henrik API docs mention consent and rate limits.
- Basic key allows 30 requests per minute.
- Enhanced key allows 90 requests per minute.
- Analytics without user consent is not supported and can be banned.
- `go-valorant-api` targets Henrik API v3.0.2 and should be treated as reference only.
- `tailieuapivalorant.md` describes API v4.5.0, so the app should own its adapter.
- Wails v2 is stable; Wails v3 is alpha.

## 4. Product Scope

### MVP Goals

- User can enter Riot name and tag.
- User must explicitly consent before personal data is fetched.
- App resolves account to PUUID.
- App fetches match history and MMR data.
- App caches API responses locally.
- App stores normalized match data in SQLite.
- App analyzes match data into tactical insights.
- App generates training recommendations.
- App presents reports in a desktop React UI.
- App supports settings, API key, refresh, and local data reset.

### Non-MVP

- No public analytics platform.
- No Riot credential collection.
- No store checker.
- No account checker.
- No hidden-name reveal or TOS-breaking feature.
- No cloud backend.
- No ML/AI recommendation engine in first version.
- No large-scale scraping.

## 5. User Roles

Player needs to understand personal weaknesses, track improvement, and get practical recommendations.

Coach needs to review player match history, identify map/economy/round/agent patterns, and prepare training plans.

Analyst needs to compare tactical patterns, inspect round-level evidence, and summarize reports.

## 6. Core Workflows

### First Run

```text
Open app
  -> read privacy/consent explanation
  -> enter Riot name and tag
  -> optionally enter API key
  -> choose region/affinity
  -> confirm consent
  -> fetch account profile
  -> save player locally
  -> open dashboard
```

### Import Match History

```text
Dashboard
  -> refresh matches
  -> Go checks rate limit/cache
  -> fetch match list from Henrik API
  -> normalize provider response
  -> dedupe by match id
  -> save into SQLite
  -> show import summary
```

### Generate Tactical Report

```text
Stored matches
  -> analysis engine computes metrics
  -> rule engine detects weaknesses
  -> recommendations generated
  -> report saved locally
  -> React displays report with evidence
```

### Training Recommendation

```text
Tactical report
  -> list weaknesses
  -> map weakness to drill
  -> assign severity/confidence
  -> user marks drill as planned/done
```

### Data Reset

```text
Settings
  -> user clicks reset data
  -> native confirmation dialog
  -> delete local player/matches/reports/settings
  -> return to setup screen
```

## 7. Architecture

Target dependency flow:

```text
domain
  <- application
      <- infrastructure
          <- interface
              <- Wails bindings / React UI
```

Domain contains pure business concepts: player, account, match, round, map, agent, economy, tactical metric, weakness, and training recommendation.

Application contains use cases: register player consent, lookup player, refresh match history, generate tactical report, list reports, save training progress, and reset local data.

Infrastructure contains adapters: Henrik API client, rate limiter, API cache, SQLite repositories, settings repository, file export/import, and logger.

Interface contains Wails-bound services: `PlayerService`, `MatchService`, `AnalysisService`, and `SettingsService`.

Frontend contains React screens and calls only Wails-generated bindings.

## 8. Proposed Folder Structure

```text
.
├─ docs/
│  ├─ product/
│  ├─ decisions/
│  ├─ skills/
│  ├─ stories/
│  ├─ TEST_MATRIX.md
│  └─ VALORANT_TACTICAL_TRAINER_PLAN.md
├─ desktop/
│  ├─ internal/
│  ├─ frontend/
│  ├─ app.go
│  ├─ main.go
│  ├─ go.mod
│  └─ wails.json
├─ internal/ (reserved if a future root module is needed)
│  ├─ domain/
│  ├─ application/
│  ├─ infrastructure/
│  └─ interface/
├─ frontend/
│  └─ src/
├─ app.go
├─ main.go
├─ go.mod
└─ wails.json
```

## 9. API Strategy

Use Henrik unofficial VALORANT API through a custom adapter in `internal/infrastructure/valorantapi`.

Do not bind domain logic directly to `go-valorant-api` because the local wrapper references v3.0.2 while current docs reference v4.5.0.

Required MVP endpoints:

- Account by name/tag.
- Account by PUUID.
- Matches by PUUID.
- Lifetime matches by PUUID.
- MMR by PUUID.
- MMR history by PUUID.
- Content metadata.

Normalized app errors:

- `RateLimited`
- `ProviderUnavailable`
- `NotFound`
- `InvalidPlayer`
- `ConsentRequired`
- `UnauthorizedApiKey`
- `DecodeFailed`
- `UnknownProviderError`

Rate limit design:

- Track request timestamps.
- Track provider `used`, `remaining`, and `reset` data when available.
- Queue or reject refresh if limit would be exceeded.
- Cache repeated reads.
- Show rate state in settings/dashboard.

## 10. Consent And Privacy

Consent is a product requirement.

Rules:

- App must not fetch personal match/account data before consent.
- Consent record is stored locally.
- User can revoke consent by deleting local data.
- User can view what is stored.
- User can delete local data.
- App must not collect Riot credentials.
- App must not perform public analytics over non-consenting users.
- App must not reveal hidden names or do account-checking behavior.

## 11. SQLite Data Model

Initial tables:

- `players`
- `consents`
- `settings`
- `api_cache`
- `matches`
- `match_players`
- `match_teams`
- `match_rounds`
- `analysis_reports`
- `analysis_findings`
- `training_recommendations`

## 12. Tactical Analysis Domain

Player-level metrics:

- KDA.
- Kill/death ratio.
- Assist rate.
- Headshot percentage.
- Damage made/received.
- Agent performance.
- Map performance.
- Recent trend.

Round-level metrics:

- Round win/loss.
- Plant round outcomes.
- Defuse round outcomes.
- Plant site distribution.
- Economy effectiveness.
- Ability casts per round.
- AFK/friendly-fire/spawn behavior if available.

Every finding must include type, severity, confidence, title, description, evidence, and suggested action.

Every recommendation must include title, drill, priority, reason, linked finding, evidence, and status.

## 13. React UI Plan

Screens:

- Setup screen.
- Dashboard screen.
- Match list screen.
- Match detail screen.
- Tactical report screen.
- Training plan screen.
- Settings screen.

UI direction:

- Dark tactical dashboard.
- Valorant-inspired red accents without copying official assets too closely.
- Card layout with map/agent filters.
- Desktop optimized at 1280x720+ and still usable in narrow windows.

## 14. Wails Bindings

Initial Go services exposed to frontend:

- `PlayerService`
- `MatchService`
- `AnalysisService`
- `SettingsService`

Frontend must treat generated bindings as the app API client.

## 15. Skills To Add

Project skills live under `docs/skills/`:

- `wails-desktop-app`
- `valorant-api-client`
- `tactical-analysis-domain`
- `react-desktop-ui`
- `validation-release`

## 16. Product Docs To Create

- `docs/product/overview.md`
- `docs/product/privacy-consent.md`
- `docs/product/api-contracts.md`
- `docs/product/tactical-analysis.md`
- `docs/product/desktop-app.md`

## 17. Decisions To Record

- `0004-desktop-stack-wails-react-go.md`
- `0005-valorant-api-adapter.md`
- `0006-local-first-sqlite.md`

## 18. Epic Backlog

- E01 Project Foundation.
- E02 Consent & Settings.
- E03 Valorant API Integration.
- E04 Local Storage.
- E05 Tactical Analysis Engine.
- E06 React Analytics UI.
- E07 Validation & Release.

## 19. First Stories

- US-001 Initialize Wails React Go app.
- US-002 Configure project docs and skills.
- US-003 Consent and player lookup.

## 20. Validation Ladder

Expected checks after implementation begins:

```text
pnpm lint
pnpm typecheck
pnpm test
go test ./...
wails build
manual desktop smoke
```

Do not claim validation passed unless the command exists and was run.

## 21. Execution Order

1. Create `docs/VALORANT_TACTICAL_TRAINER_PLAN.md`.
2. Create product docs under `docs/product/`.
3. Create decision records under `docs/decisions/`.
4. Create skills under `docs/skills/`.
5. Update `docs/stories/backlog.md`.
6. Update `docs/TEST_MATRIX.md`.
7. Create first story files.
8. Scaffold Wails v2 React TypeScript app in `desktop/`.
9. Configure pnpm and Tailwind.
10. Add Go internal layers under `desktop/internal/`.
11. Add SQLite.
12. Add API adapter.
13. Add consent/player lookup.
14. Add match import/cache.
15. Add tactical analysis engine.
16. Add UI screens.
17. Add validation and release build.

## 22. Risks And Mitigation

API privacy/TOS risk is mitigated by consent gate, local-only data, no credentials, no public analytics, and data reset/export.

API version drift is mitigated by custom adapter, provider DTO boundary, mock tests, and raw payload retention for debugging.

Rate limit risk is mitigated by cache, request queue, rate status UI, reset-aware backoff, and batch refresh controls.

Desktop framework risk is mitigated by Wails v2 stable, thin Wails boundary, `wails build`, and avoiding v3-only APIs.

Analysis quality risk is mitigated by requiring evidence, confidence score, graceful missing-data handling, and rule-based MVP.

## 23. Definition Of Done

A story is done when:

- Product contract is implemented or blocker documented.
- Relevant product docs are updated.
- Story packet evidence is updated.
- Test matrix row is updated.
- Validation commands were run if they exist.
- Missing validation is explicitly documented.
- No provider API call bypasses consent/rate-limit/cache rules.
- Final summary states changed files and validation result.
