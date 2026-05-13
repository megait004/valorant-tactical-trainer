---
name: validation-release
description: Use before finishing a story, preparing a demo, or creating a release build.
---

# Validation And Release

## Use When

- Finishing a story.
- Preparing demo evidence.
- Preparing a release build.
- Updating `docs/TEST_MATRIX.md`.
- Checking whether implementation is ready to hand off.

## Validation Ladder

Expected commands after implementation begins:

```text
pnpm lint
pnpm typecheck
pnpm test
go test ./...
wails build
manual desktop smoke
```

## Rules

- Do not claim validation passed unless the command exists and was run.
- If a command does not exist, document it as missing.
- For API work, test with a mock provider first.
- For analysis work, use deterministic unit tests.
- For desktop work, include manual smoke evidence when automated E2E does not exist.
- Update `docs/TEST_MATRIX.md` before closing the story.

## Done Checklist

- Product contract updated if behavior changed.
- Story evidence updated.
- Test matrix updated.
- Validation output recorded or missing command documented.
- No consent/rate-limit/cache bypass introduced.
