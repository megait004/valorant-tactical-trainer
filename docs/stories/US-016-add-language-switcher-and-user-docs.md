# US-016 Add Language Switcher And User Docs

## Status

implemented

## Lane

normal

## Product Contract

The app lets users switch core UI labels between English and Vietnamese. The selected language is persisted locally in SQLite. The repository also includes user-facing instructions and a code structure document for graduation thesis reporting.

## Acceptance Criteria

- Header exposes a language selector with English and Tiếng Việt.
- Language is loaded from local settings on startup.
- Language changes are saved through Go/Wails, not browser storage.
- Core labels in setup, assistant, settings, match cache, and report panels use translation keys.
- README explains installation, validation, core flows, privacy, and VTA safety.
- Code structure doc explains architecture and folders for DATN reporting.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | Existing Go tests pass. |
| Integration | Wails bindings compile with `SaveLanguage`. |
| E2E | Manual language switch smoke not run. |
| Platform | Wails build passes. |

## Evidence

- `go test ./...` passed before docs.
- `pnpm validate` passed before docs.
- `go test ./...` passed after docs.
- `pnpm validate` passed after docs.
- `wails build` passed after docs.
