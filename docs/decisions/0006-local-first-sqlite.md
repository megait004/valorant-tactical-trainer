# 0006 Local First SQLite

Date: 2026-05-13

## Status

Accepted

## Context

The app is a desktop tool that stores consented player data, imported matches, reports, recommendations, settings, and API cache locally. Match analysis benefits from structured queries and durable storage.

## Decision

Use SQLite as the primary local storage engine.

## Alternatives Considered

1. JSON files: simpler, but weak for query-heavy match analysis and migrations.
2. Embedded key-value store: good for cache/settings, weaker for relational match data.
3. Cloud database: unnecessary for local-first MVP and increases privacy risk.

## Consequences

Positive:

- No external backend required.
- Good fit for structured match/report queries.
- Supports export/import and deterministic tests.
- Keeps user data local.

Tradeoffs:

- Requires migrations.
- Requires repository layer and test fixtures.
- Reset/delete behavior must be carefully validated.

## Follow-Up

- Add migrations in the storage story.
- Add repository tests for consent, matches, reports, and reset behavior.
