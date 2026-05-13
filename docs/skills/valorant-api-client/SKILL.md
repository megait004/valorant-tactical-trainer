---
name: valorant-api-client
description: Use when changing Henrik API integration, provider DTOs, endpoints, cache, rate limit, errors, or consent-gated fetch behavior.
---

# Valorant API Client

## Use When

- Adding or changing a Henrik API endpoint.
- Changing provider DTOs or response parsing.
- Handling 429/rate-limit behavior.
- Changing API cache behavior.
- Changing API key settings.
- Fetching account, match, MMR, or content metadata.

## Rules

- Require explicit user consent before personal data fetches.
- Assume Basic limit of 30 requests per minute unless settings indicate otherwise.
- Respect Enhanced limit of 90 requests per minute when configured.
- Handle 429 with reset-aware behavior.
- Never store Riot credentials.
- Never implement store checker, hidden-name reveal, or account checker behavior.
- Keep provider DTOs inside infrastructure.
- Normalize provider data before it enters application/domain code.
- Normalize provider errors into app errors.
- Use cache for repeated read calls.

## Error Contract

Map provider behavior to:

- `ConsentRequired`
- `RateLimited`
- `ProviderUnavailable`
- `NotFound`
- `InvalidPlayer`
- `UnauthorizedApiKey`
- `DecodeFailed`
- `UnknownProviderError`

## Validation

- Add or update mock provider tests for new endpoint behavior.
- Test rate-limit behavior without calling the real provider.
- Do not rely on live API tests as the only proof.
