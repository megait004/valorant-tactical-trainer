# API Contracts

## Provider

The external provider is Henrik unofficial VALORANT API.

The app uses a custom Go adapter instead of coupling domain logic directly to `go-valorant-api`.

Reason:

- Local `go-valorant-api` references API v3.0.2.
- Downloaded provider docs reference v4.5.0.
- The app needs consent, cache, rate-limit, and error normalization around provider calls.

## Rate Limits

Provider docs describe:

- Basic key: 30 requests per minute.
- Enhanced key: 90 requests per minute.

The app must assume the lower Basic limit unless settings indicate otherwise.

The default desktop provider client spaces live calls by 2 seconds, matching the Basic 30 requests per minute limit. Tests can use an unthrottled client or a shorter injected limiter.

## Adapter Boundary

Provider responses enter the app only through `internal/infrastructure/valorantapi`.

The adapter must:

- Build provider URLs.
- Attach API key when configured.
- Decode provider responses.
- Normalize provider errors.
- Update rate-limit state.
- Cache successful responses when configured.
- Convert provider DTOs into app-facing DTOs.

Domain and React code must not depend on provider DTOs.

## MVP Endpoints

- Account by name/tag.
- Account by PUUID.
- Matches by PUUID.
- Lifetime matches by PUUID.
- MMR by PUUID.
- MMR history by PUUID.
- Content metadata.

## Error Model

Normalize provider failures into these app errors:

| Error | Meaning |
| --- | --- |
| `ConsentRequired` | Request attempted before consent. |
| `RateLimited` | Provider returned or predicted rate limit. |
| `ProviderUnavailable` | Provider cannot be reached or returns temporary failure. |
| `NotFound` | Player or requested resource not found. |
| `InvalidPlayer` | Name/tag/PUUID input is invalid. |
| `UnauthorizedApiKey` | API key is missing, invalid, or rejected. |
| `DecodeFailed` | Provider payload could not be decoded. |
| `UnknownProviderError` | Provider returned an unclassified error. |

## Cache Policy

Cache key should include:

- endpoint.
- normalized params.
- player PUUID when relevant.
- provider version if known.

Cache rules:

- Account profile can be refreshed manually.
- Match details are immutable enough to cache long term.
- Match history pages can be cached with refresh controls.
- MMR history can be cached but should allow manual refresh.

## Frontend Contract

React calls only Wails bindings.

React receives UI-safe DTOs and never receives raw provider payload unless a debugging feature is explicitly added later.
