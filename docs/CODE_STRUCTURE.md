# Code Structure Report

Tài liệu này mô tả cấu trúc mã nguồn của Valorant Tactical Trainer để dùng trong báo cáo đồ án tốt nghiệp.

## Tổng Quan Kiến Trúc

Ứng dụng dùng mô hình desktop local-first:

```text
React UI -> Wails bindings -> Go services -> Domain logic / Infrastructure -> SQLite / Henrik API
```

Nguyên tắc chính:

- React chỉ gọi Go thông qua Wails-generated bindings.
- Go backend sở hữu provider API, consent, cache, SQLite, phân tích, recommendations, settings và export.
- Dữ liệu được lưu cục bộ trên máy người dùng.
- Không có backend server riêng.
- Không đọc bộ nhớ game hoặc inject vào VALORANT.

## Thư Mục Gốc

```text
DATN/
├── desktop/                  # Ứng dụng Wails chính
├── docs/                     # Tài liệu sản phẩm, stories, test matrix
├── README.md                 # Hướng dẫn sử dụng app
└── .gitignore                # Loại trừ node_modules, build, DB local, reference lớn
```

## Desktop App

```text
desktop/
├── app.go                    # Binding smoke AppInfo
├── main.go                   # Khởi tạo Wails, storage, services
├── wails.json                # Cấu hình Wails, frontend command dùng pnpm
├── go.mod                    # Go module dependencies
├── VALIDATION.md             # Checklist validate
├── build/                    # Output build, bị ignore
├── frontend/                 # React + TypeScript UI
└── internal/                 # Go backend layers
```

## Wails Entry Point

`desktop/main.go` làm nhiệm vụ:

- mở SQLite store.
- tạo Wails services.
- bind services vào Wails runtime.
- cấu hình window title/size/assets.
- gọi `settingsService.Startup(ctx)` để service dùng native dialogs.

Services đang bind:

- `App`
- `PlayerService`
- `MatchService`
- `RankService`
- `AnalysisService`
- `AssistantService`
- `SettingsService`

## Go Internal Layers

```text
desktop/internal/
├── domain/                   # Business/domain logic thuần
├── infrastructure/           # SQLite + provider API adapter
└── interface/                # Wails service DTO/methods
```

## Domain Layer

Domain layer không phụ thuộc React, Wails, SQLite hoặc provider DTO.

```text
desktop/internal/domain/
├── analysis/
│   ├── analysis.go           # Tactical report + recommendations rules
│   └── analysis_test.go
├── assistant/
│   ├── assistant.go          # VTA tactical cards + economy rules
│   └── assistant_test.go
├── match/
│   └── match.go              # Match summary entity
├── player/
│   ├── player.go             # Account, consent, name/tag normalization
│   └── player_test.go
└── rank/
    └── rank.go               # Rank/MMR snapshot entity
```

### `domain/player`

Chứa:

- `Account`: thông tin player sau lookup.
- `Consent`: record consent local.
- normalize Riot name/tag.
- consent version.

### `domain/match`

Chứa `Summary`, đại diện dữ liệu match đã được normalize để phân tích và hiển thị.

### `domain/rank`

Chứa `Snapshot`, đại diện một lần fetch rank/MMR.

### `domain/analysis`

Chứa:

- `Report`
- `Finding`
- `Recommendation`
- rule-based analysis engine.

Các metric MVP:

- KDA.
- Headshot percentage.
- Average damage.
- Top agent.
- Top map.

### `domain/assistant`

Chứa Virtual Tactical Assistant MVP:

- `Query`: map/agent/side/phase/credits/previous outcome.
- `TacticalCard`: thẻ chiến thuật local.
- `EconomyAdvice`: gợi ý mua đồ.
- `RecommendEconomy`: rule engine kinh tế.
- `SeedCards`: seed tactical cards ban đầu.

## Infrastructure Layer

```text
desktop/internal/infrastructure/
├── storage/
│   ├── store.go              # SQLite migrations/repositories/export/settings/cache
│   └── store_test.go
└── valorantapi/
    ├── client.go             # Henrik API adapter + rate limiter
    └── client_test.go
```

### `infrastructure/storage`

`Store` quản lý SQLite.

Các bảng chính:

- `players`: thông tin player.
- `consents`: consent records.
- `settings`: API key, language, current player.
- `api_cache`: cache provider response có TTL.
- `matches`: match summaries.
- `rank_snapshots`: lịch sử rank/MMR fetch.
- `analysis_reports`: report metadata.
- `analysis_findings`: findings theo report.
- `training_recommendations`: recommendations theo report.
- `tactical_cards`: dữ liệu VTA local.

Các nhóm hàm chính:

- player/consent: `SavePlayerWithConsent`, `CurrentPlayer`.
- settings: `SaveSetting`, `DeleteSetting`, `Setting`.
- cache: `SaveAPICache`, `APICache`, `ClearExpiredAPICache`.
- stats/export: `Stats`, `ExportSnapshot`, `ExportJSON`.
- matches: `SaveMatches`, `MatchesForPlayer`.
- rank: `SaveRankSnapshot`, `LatestRankSnapshot`.
- reports: `SaveReport`.
- assistant: `SeedTacticalCards`, `TacticalCards`.
- reset: `ResetAll`.

### `infrastructure/valorantapi`

Go adapter cho Henrik API.

Chức năng:

- lookup account.
- fetch matches by PUUID.
- fetch MMR/rank by PUUID.
- inject test HTTP client/base URL.
- Basic provider rate limit 2 giây/request.

Frontend không gọi Henrik API trực tiếp.

## Wails Interface Layer

```text
desktop/internal/interface/wails/
├── analysis_service.go       # GenerateReport
├── assistant_service.go      # QueryAssistant
├── match_service.go          # RefreshMatches/ListMatches
├── player_service.go         # LookupPlayer/GetCurrentPlayer
├── rank_service.go           # RefreshRank/LatestRank
└── settings_service.go       # Settings/reset/cache/export/language
```

Layer này làm nhiệm vụ:

- nhận input từ React.
- validate input cơ bản.
- gọi domain/infrastructure.
- trả DTO an toàn cho UI.
- không trả API key đã lưu về frontend.

## Frontend Structure

```text
desktop/frontend/
├── package.json              # pnpm scripts, validate
├── src/
│   ├── App.tsx               # State orchestration + gọi Wails bindings
│   ├── i18n.ts               # English/Tiếng Việt labels
│   ├── main.tsx              # React entry
│   └── components/
│       ├── AppHeader.tsx
│       ├── MatchCachePanel.tsx
│       ├── ReportPanel.tsx
│       ├── SettingsPanel.tsx
│       ├── SetupPanel.tsx
│       ├── VirtualAssistantPanel.tsx
│       └── types.ts
└── wailsjs/                  # Generated Wails bindings
```

### `App.tsx`

Quản lý state chính:

- current player.
- matches.
- rank.
- tactical report.
- assistant query/result.
- settings.
- language.
- loading states.
- status message.

Gọi Wails bindings:

- `AppInfo`
- `LookupPlayer`
- `GetCurrentPlayer`
- `RefreshMatches`
- `ListMatches`
- `RefreshRank`
- `LatestRank`
- `GenerateReport`
- `QueryAssistant`
- `GetSettings`
- `SaveSettings`
- `SaveLanguage`
- `ClearExpiredCache`
- `ExportLocalData`
- `ResetAllData`

### `i18n.ts`

Chứa translation object cho:

- English.
- Tiếng Việt.

Ngôn ngữ được lưu vào SQLite qua `SaveLanguage` và load lại qua `GetSettings`.

### Components

- `AppHeader`: title, status, language selector.
- `SetupPanel`: consent, player lookup, refresh matches/rank, report, reset.
- `VirtualAssistantPanel`: VTA input + tactical cards + economy advice.
- `SettingsPanel`: API key, cache cleanup, export, stats.
- `MatchCachePanel`: danh sách match đã lưu.
- `ReportPanel`: findings và training recommendations.
- `types.ts`: props types dùng chung.

## Data Flow Chính

### Player Lookup

```text
User input + consent -> React -> PlayerService.LookupPlayer -> Henrik API adapter -> SQLite players/consents/settings -> React DTO
```

### Match Refresh

```text
React -> MatchService.RefreshMatches -> Henrik matches endpoint -> normalize -> SQLite matches/cache -> React list
```

### Rank Refresh

```text
React -> RankService.RefreshRank -> Henrik MMR endpoint -> SQLite rank_snapshots -> React rank card
```

### Tactical Report

```text
React -> AnalysisService.GenerateReport -> SQLite matches -> domain analysis -> SQLite reports/findings/recommendations -> React report panel
```

### Virtual Tactical Assistant

```text
Manual query -> React -> AssistantService.QueryAssistant -> SQLite tactical_cards + economy rules -> React VTA panel
```

### Export Data

```text
React -> SettingsService.ExportLocalData -> native save dialog -> Store.ExportJSON -> write JSON file
```

Export loại trừ:

- API key value.
- raw provider payloads.

## Validation Strategy

Các lệnh validate chính:

```text
go test ./...
pnpm validate
wails build
```

Test hiện có:

- player normalization.
- provider mock account/matches/MMR.
- rate limiter.
- SQLite player/cache/matches/rank/report/reset/export/settings.
- analysis domain rules.
- assistant economy rules và seed cards.

## Giới Hạn Hiện Tại

- Chưa có overlay always-on-top riêng.
- Chưa có hotkey global.
- Chưa có manual live smoke với account thật trong tài liệu này.
- VTA chưa có ảnh/video lineups do cần asset tự tạo hoặc có quyền sử dụng.
- Match detail/round timeline/filter nâng cao vẫn là hướng phát triển tiếp.

## Hướng Phát Triển

- Match detail và filters theo map/agent/date.
- Compact VTA window/always-on-top mode có nút tắt nhanh.
- Import JSON local data.
- Thêm tactical cards bằng asset tự thiết kế.
- Demo script và installer packaging cho nghiệm thu.
