# US-001 Initialize Wails React Go App

## Status

implemented

## Lane

normal

## Product Contract

The project has a runnable Wails v2 desktop app in `desktop/` with React TypeScript frontend and Go backend callable through Wails bindings.

## Relevant Product Docs

- `docs/product/overview.md`
- `docs/product/desktop-app.md`
- `docs/decisions/0004-desktop-stack-wails-react-go.md`
- `docs/skills/wails-desktop-app/SKILL.md`
- `docs/skills/react-desktop-ui/SKILL.md`

## Acceptance Criteria

- Wails app opens a desktop window.
- React renders an initial app shell.
- A Go method can be called from React through generated bindings.
- Frontend uses pnpm.
- Frontend uses React, TypeScript, and Tailwind CSS.
- Basic validation commands are documented.

## Design Notes

- Commands: initialize Wails v2 React TypeScript template, add validation scripts.
- Queries: none for MVP shell.
- API: Wails binding smoke method only.
- Tables: none.
- Domain rules: none yet.
- UI surfaces: initial shell with navigation placeholders.

## Validation

| Layer | Expected proof |
| --- | --- |
| Unit | none for initial shell |
| Integration | Go binding smoke if feasible |
| E2E | manual app launch smoke |
| Platform | `wails build` |
| Release | not required until release story |

## Harness Delta

Uses docs and skills created in US-002.

## Evidence

- Scaffolded Wails v2 React TypeScript app in `desktop/`.
- Replaced template greet binding with `AppInfo` binding smoke.
- Switched Wails frontend commands to `pnpm`.
- Added Tailwind CSS setup.
- Added frontend `lint` and `typecheck` scripts.
- `pnpm typecheck`: passed.
- `pnpm lint`: passed.
- `go test ./...`: passed, no test files yet.
- `pnpm build`: passed.
- `wails build`: passed and created `desktop/build/bin/valorant-tactical-trainer.exe`.
