# Story Backlog

This backlog tracks the initial Valorant Tactical Trainer buildout. Create detailed story packets when a story is selected for implementation or when a durable product decision needs a work surface.

## Candidate Epics

| Epic | Description | Status |
| --- | --- | --- |
| E01 Project Foundation | Create Wails v2 React Go app, base docs, skills, validation, and app shell. | in_progress |
| E02 Consent & Settings | Implement consent, player setup, API key, local settings, and data reset. | in_progress |
| E03 Valorant API Integration | Implement Henrik API adapter, rate limiter, cache, account lookup, match fetch, and MMR fetch. | planned |
| E04 Local Storage | Add SQLite migrations, repositories, import dedupe, reports, and data export/import. | in_progress |
| E05 Tactical Analysis Engine | Normalize match data and generate metrics, findings, and recommendations. | planned |
| E06 React Analytics UI | Build setup, dashboard, match list, match detail, report, training, and settings screens. | in_progress |
| E07 Validation & Release | Add tests, mock provider checks, desktop smoke, release build, and demo script. | planned |

## Initial Story Candidates

| Story | Epic | Title | Lane | Status |
| --- | --- | --- | --- | --- |
| US-001 | E01 | Initialize Wails React Go app | normal | implemented |
| US-002 | E01 | Configure project docs and skills | normal | completed |
| US-003 | E02 | Consent and player lookup | high-risk | implemented |
| US-004 | E03 | Fetch and cache match history | high-risk | implemented |
| US-005 | E05 | Generate first tactical report | normal | implemented |
| US-006 | E05 | Generate training recommendation list | normal | implemented |
| US-007 | E02 | Desktop settings and data reset | high-risk | implemented |
| US-008 | E04 | Add SQLite schema and repositories | high-risk | implemented |
| US-009 | E06 | Build dashboard and match list UI | normal | implemented |
| US-010 | E07 | Add validation and release build flow | normal | implemented |
| US-011 | E03 | Add provider throttle | high-risk | implemented |
