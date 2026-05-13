# 0004 Desktop Stack Wails React Go

Date: 2026-05-13

## Status

Accepted

## Context

The product is a desktop tactical analysis app that needs local storage, provider API integration, and a modern UI. The selected stack must support a Go core app, React UI, and a desktop release suitable for DATN/demo work.

Wails has two active versions: v2 stable and v3 alpha. The project needs stability more than alpha-only APIs.

## Decision

Use Wails v2 stable with Go backend/core and React TypeScript frontend.

Use Tailwind CSS for UI styling.

React communicates with Go through Wails-generated bindings.

## Alternatives Considered

1. Wails v3 alpha: better newer APIs, but higher release risk.
2. Electron: mature desktop ecosystem, but heavier and less aligned with Go core.
3. Pure web app: easier deployment, but misses local-first desktop requirement.

## Consequences

Positive:

- Single desktop app with Go and web UI.
- Go can own API/storage/analysis cleanly.
- React can focus on presentation.
- Wails v2 is stable enough for demo/release.

Tradeoffs:

- Multi-window and newer v3 APIs are not MVP assumptions.
- Build/release must validate Wails-specific tooling.
- Generated bindings must be refreshed after exported Go service changes.

## Follow-Up

- Scaffold Wails v2 React TypeScript app in US-001.
- Add validation around `wails build` before release.
