---
name: react-desktop-ui
description: Use when changing React screens, state, charts, forms, desktop UX, or Wails binding consumers.
---

# React Desktop UI

## Use When

- Building setup, dashboard, match list, match detail, report, training, or settings screens.
- Changing frontend state or Wails binding consumers.
- Adding charts, tables, filters, or forms.
- Handling desktop UI loading/error/empty states.

## Rules

- Use React and TypeScript.
- Use Tailwind CSS, not plain CSS.
- Use arrow functions, destructuring, and template literals.
- Do not call Henrik API directly from React.
- Treat Wails-generated bindings as the frontend API client.
- Handle loading, empty, error, and rate-limited states.
- Keep screens usable on desktop and narrower windows.
- Preserve dark tactical visual language with restrained Valorant-inspired red accents.

## State Guidance

- Keep provider/API state in Go where possible.
- Keep UI-only state in React.
- Avoid overusing memoization unless measurement or existing patterns justify it.
- Use concurrent React patterns only when they improve user-perceived responsiveness.

## Validation

- Run `pnpm lint` when available.
- Run `pnpm typecheck` when available.
- Run UI tests when available.
- Manually smoke desktop flows when Wails app exists.
