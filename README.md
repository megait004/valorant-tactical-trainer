# Valorant Tactical Trainer

Ứng dụng desktop hỗ trợ phân tích dữ liệu VALORANT, tạo báo cáo chiến thuật, gợi ý luyện tập và cung cấp Virtual Tactical Assistant an toàn theo dạng tra cứu thủ công.

## Công Nghệ

- Wails v2: đóng gói desktop app Windows.
- Go: backend, adapter API, SQLite, business logic, Wails services.
- React + TypeScript: giao diện người dùng.
- Tailwind CSS: styling.
- SQLite: lưu dữ liệu cục bộ.
- Henrik unofficial VALORANT API: lấy dữ liệu account, match history, rank/MMR khi người dùng consent.

## Tính Năng Chính

- Consent gate trước khi fetch dữ liệu cá nhân.
- Tra cứu Riot name/tag và lưu player hiện tại.
- Fetch/cache match history cục bộ.
- Fetch/cache current rank/MMR.
- Tạo tactical report từ dữ liệu match đã lưu.
- Gợi ý training recommendations có evidence.
- Virtual Tactical Assistant MVP:
  - chọn map/agent/side/phase thủ công.
  - xem tactical cards cho pre-match và in-game.
  - Economy Manager gợi ý Eco, Light/Half Buy, Force Buy, Full Buy.
  - không đọc memory game, không inject, không tự động đọc trạng thái trận.
- Settings panel:
  - lưu/xoá API key local.
  - đổi ngôn ngữ English/Tiếng Việt.
  - xem số lượng dữ liệu local.
  - clear expired cache.
  - export JSON privacy-aware.
  - reset toàn bộ dữ liệu local bằng native confirmation dialog.

## Chạy App

Yêu cầu máy dev đã có:

- Go.
- Wails v2 CLI.
- pnpm.
- WebView2 Runtime trên Windows.

Chạy dev mode:

```powershell
cd desktop
wails dev
```

Build release:

```powershell
cd desktop
wails build
```

File exe sau build:

```text
desktop/build/bin/valorant-tactical-trainer.exe
```

## Validation

Backend tests:

```powershell
cd desktop
go test ./...
```

Frontend lint/typecheck/build:

```powershell
cd desktop/frontend
pnpm validate
```

Release build:

```powershell
cd desktop
wails build
```

## Hướng Dẫn Sử Dụng

### 1. Đổi Ngôn Ngữ

- Mở app.
- Ở góc phải header, chọn `Language`.
- Chọn `English` hoặc `Tiếng Việt`.
- Lựa chọn được lưu trong SQLite local và sẽ được load lại khi mở app sau.

### 2. Cấu Hình API Key

- API key là tuỳ chọn.
- Nhập key vào ô `API key` trong Settings hoặc Setup.
- Bấm `Save key`.
- App chỉ hiển thị trạng thái đã có key, không hiển thị lại giá trị key đã lưu.
- Muốn xoá key: để trống ô API key rồi bấm save.

### 3. Tra Cứu Người Chơi

- Nhập Riot name.
- Nhập tag, ví dụ `VN2`.
- Chọn region fallback, ví dụ `AP`.
- Tick consent checkbox.
- Bấm `Lookup player`.
- App gọi Henrik API qua Go backend và lưu player/consent vào SQLite.

### 4. Refresh Match History

- Sau khi đã lookup player, bấm `Refresh matches`.
- App fetch match history qua Go backend.
- Dữ liệu match được dedupe và lưu trong bảng `matches`.
- Match cache hiển thị ở panel bên phải.

### 5. Refresh Rank/MMR

- Sau khi đã lookup player, bấm `Refresh rank`.
- App fetch current rank/MMR và lưu vào `rank_snapshots`.
- UI hiển thị tier, RR, elo, last-game MMR change và thời gian fetch.

### 6. Tạo Tactical Report

- Sau khi có match history, bấm `Generate report`.
- Go domain engine tính các metric như KDA, headshot percentage, average damage, top agent, top map.
- App tạo findings và training recommendations có evidence.

### 7. Dùng Virtual Tactical Assistant

- Chọn map: `Ascent`, `Bind`, `Haven`.
- Chọn agent hoặc `Any`.
- Chọn side: Attack, Defense, Both.
- Chọn phase: Pre-match hoặc In-game.
- Nhập credits hiện tại.
- Chọn kết quả round trước: Win hoặc Loss.
- Bấm `Get tactical cards`.
- App trả về tactical cards và Economy Manager advice.

Lưu ý an toàn:

- VTA MVP chỉ dùng dữ liệu local/manual.
- Không đọc bộ nhớ `VALORANT.exe`.
- Không inject overlay vào game.
- Không tự động phát hiện trạng thái trận.

### 8. Export Dữ Liệu Local

- Trong Settings, bấm `Export JSON`.
- Chọn nơi lưu file JSON bằng native save dialog.
- File export không chứa API key value và không chứa raw provider payloads.

### 9. Reset Dữ Liệu Local

- Bấm `Reset local data`.
- App hiện native confirmation dialog.
- Nếu xác nhận, app xoá player, consent, cache, matches, rank snapshots, reports, recommendations và settings.

## Dữ Liệu Lưu Ở Đâu

SQLite DB nằm trong thư mục config user của Windows:

```text
%AppData%/ValorantTacticalTrainer/trainer.db
```

Đường dẫn chính xác được hiển thị trong Settings panel.

## Quy Tắc Privacy Và Anti-Cheat

- Người dùng phải consent trước khi fetch dữ liệu cá nhân.
- Không lưu Riot credentials.
- Không xây public analytics database.
- Không reveal hidden names.
- Không bypass rate limit.
- Không đọc memory hoặc inject vào VALORANT.
- API key chỉ lưu local qua Go-managed settings.

## Tài Liệu Liên Quan

- `docs/CODE_STRUCTURE.md`: cấu trúc code chi tiết để báo cáo.
- `docs/TEST_MATRIX.md`: ma trận kiểm thử.
- `docs/product/`: product contracts.
- `docs/stories/`: user stories đã implement.
- `desktop/VALIDATION.md`: checklist validate app.
