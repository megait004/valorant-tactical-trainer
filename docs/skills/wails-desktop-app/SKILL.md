---
name: wails-desktop-app
description: Use when changing Wails app shell, Go bindings, native dialogs, app lifecycle, desktop build, or the React/Go bridge.
---

# Wails Desktop App

## Use When

- Adding or changing Wails services.
- Changing exported Go methods used by React.
- Working on desktop lifecycle, native dialogs, or app settings.
- Running or fixing Wails dev/build behavior.
- Changing frontend/backend boundary contracts.

## Rules

- Use Wails v2 stable.
- Go owns app logic, API calls, storage, rate limits, consent, analysis, and recommendations.
- React calls Go only through Wails-generated bindings.
- Do not call Henrik API directly from React.
- Keep Wails interface services thin.
- Regenerate bindings after exported Go service changes.
- Return UI-safe DTOs from Wails services.
- Use native confirmation dialogs before destructive reset actions.

## Validation

- Run Go tests for changed backend code when tests exist.
- Run frontend typecheck when binding consumers change.
- Run `wails build` before release/demo when command exists.
- Document missing validation in story evidence.
