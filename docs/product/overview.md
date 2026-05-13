# Product Overview

## Vision

Valorant Tactical Trainer is a local-first desktop app that helps VALORANT players and coaches turn match history into tactical training actions.

The app should answer three questions:

- What patterns are hurting performance?
- What evidence supports each finding?
- What should the player train next?

## Users

Player:

- Reviews personal performance.
- Tracks recent match trends.
- Uses recommendations as training tasks.

Coach:

- Reviews a player's match history with consent.
- Finds map, economy, round, and agent weaknesses.
- Converts findings into practice focus.

Analyst:

- Inspects round and match evidence.
- Compares tactical patterns.
- Produces summaries for review.

## MVP Scope

- Consent-gated player lookup by Riot name/tag.
- Henrik API integration through Go adapter.
- Local SQLite persistence.
- Match history import and dedupe.
- MMR/history import when available.
- Tactical report generation from stored matches.
- Rule-based recommendations with evidence.
- Desktop React UI via Wails bindings.
- Settings for API key, region, cache, and local reset.

## Out Of Scope

- Public analytics over non-consenting users.
- Riot credential collection.
- Store checker features.
- Hidden-name reveal.
- Account checker behavior.
- Cloud backend.
- ML/AI recommendation engine in MVP.
- Large-scale scraping.

## Core Workflows

First run:

```text
Open app -> consent -> enter name/tag -> lookup account -> save player -> dashboard
```

Refresh data:

```text
Dashboard -> refresh -> rate limit/cache check -> fetch API -> normalize -> persist -> show summary
```

Analyze:

```text
Stored matches -> metrics -> findings -> recommendations -> report UI
```

Train:

```text
Report -> recommendation -> drill -> mark planned/done
```

Reset:

```text
Settings -> confirm native dialog -> delete local data -> setup screen
```

## Success Criteria

- A user can complete setup without reading technical docs.
- Personal data is never fetched before consent.
- Imported match data is cached and deduped.
- At least three useful tactical findings can be generated from sample match data.
- Each recommendation includes evidence and a practical drill.
- The app builds as a Wails v2 desktop app.
