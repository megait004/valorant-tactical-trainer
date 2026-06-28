# Valorant Tactical Trainer

Ung dung desktop ho tro nguoi choi Valorant xem bao cao tran dau, phan tich loi choi, goi y luyen tap va tao tactical plan. Du an dung Wails cho desktop app, Go cho backend va React/Vite cho frontend.

## Tinh nang chinh

- Dang nhap bang Riot ID va lay du lieu tran dau tu Riot/Henrik API.
- Phan tich chi so, diem manh/yeu va de xuat bai tap luyen tap.
- Chat/coach AI tuy chon khi cau hinh `LLM_API_KEY`.
- Tao va luu tactical plan, practice progress tren may local.

## Chay local

Yeu cau: Go theo `go.mod`, Node.js, pnpm va Wails CLI.

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@latest
cd frontend
pnpm install
cd ..
wails dev
```

Tao file `.env` o thu muc goc neu can dung API:

```env
RIOT_API_KEY=RGAPI-...
HENRIK_API_KEY=...
LLM_API_KEY=...
LLM_MODEL=gemini-2.5-flash-lite
```

## Build

```powershell
wails build -clean -nsis
```

File build/installer nam trong `build/bin`.


