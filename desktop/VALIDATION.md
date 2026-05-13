# Desktop Validation

Run these checks before handing off app changes:

```bash
cd desktop/frontend
pnpm lint
pnpm typecheck
pnpm build

cd ..
go test ./...
wails build
```

Expected output:

- Frontend lint passes.
- TypeScript check passes.
- Vite production build completes.
- Go unit/integration tests pass.
- Wails builds `build/bin/valorant-tactical-trainer.exe` on Windows.

Live Henrik API smoke is manual and should only be run with a consenting account.
