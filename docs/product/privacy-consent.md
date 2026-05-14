# Privacy And Consent

## Product Rule

The app must require explicit user consent before fetching personal VALORANT account or match data from Henrik API.

Consent is required because the provider documentation states that analytic services without user consent are not supported and can be banned.

## Allowed Data Use

- Fetch account data for a consenting player.
- Fetch match and MMR data for a consenting player.
- Store fetched data locally on the user's machine.
- Generate local tactical reports and training recommendations.
- Let the user delete or export local data.

## Forbidden Data Use

- Do not collect Riot credentials.
- Do not build a public analytics database of non-consenting players.
- Do not implement public store checker behavior.
- Do not reveal hidden names.
- Do not implement account checker behavior for resale or abuse.
- Do not bypass provider rate limits.

## Consent UX Requirements

Before the first API call, the setup screen must show:

- Which provider is used: Henrik unofficial VALORANT API.
- What data is fetched: account profile, match history, MMR/history when available.
- Where data is stored: local SQLite database on this machine.
- How data is used: tactical analysis and training recommendations.
- How to remove data: settings reset/delete local data.

The lookup action must be disabled until consent is checked.

## Consent Record

Store a local consent record with:

- player id or pending player identity.
- Riot name and tag.
- PUUID after lookup.
- provider name.
- consent version.
- consent timestamp.
- revoked timestamp when data is reset.

## Data Reset

Reset local data must:

- Use a native confirmation dialog.
- Delete player, consent, cached provider payloads, matches, reports, and recommendations.
- Return the UI to setup state.
- Never call the provider while resetting.

## API Key Handling

API key is optional for MVP.

If provided, it is stored locally only and used by the Go API adapter.

The frontend must not send the API key to any endpoint except through Wails-bound local Go services.

Settings UI may show whether an API key is configured, but must not display the stored API key value after saving.

Local data export must exclude saved API key values. Exported match/rank/report data should not include raw provider payloads unless a future explicit debug export story adds that behavior.

Virtual Tactical Assistant MVP must remain local/manual. It must not read VALORANT process memory, inject into the game, automate input, or infer hidden live match state without an explicit safe API and user opt-in.
