# 0005 Valorant API Adapter

Date: 2026-05-13

## Status

Accepted

## Context

The downloaded `go-valorant-api` wrapper targets Henrik API v3.0.2. The local API documentation references unofficial-valorant-api v4.5.0 and includes important consent/rate-limit constraints.

The app needs to enforce consent, handle rate limits, cache responses, normalize errors, and protect domain logic from provider drift.

## Decision

Create a custom Go adapter around Henrik API under `internal/infrastructure/valorantapi`.

Use `go-valorant-api` as reference material only unless a later story proves a specific dependency is safe and current.

## Alternatives Considered

1. Use `go-valorant-api` directly: faster initial integration, but locks app to older API assumptions.
2. Call Henrik API directly from React: simpler UI prototype, but breaks consent/rate-limit/security boundaries.
3. Add a separate backend server: unnecessary for local-first MVP.

## Consequences

Positive:

- Provider version drift is isolated.
- Rate limit and cache behavior are centralized.
- Consent can be enforced before provider calls.
- Tests can mock provider behavior.

Tradeoffs:

- More initial code.
- Endpoint DTOs must be maintained by the project.
- Integration tests are required to keep provider behavior safe.

## Follow-Up

- Define provider DTOs only inside infrastructure.
- Add mock API tests for account lookup, match fetch, rate limit, and provider errors.
