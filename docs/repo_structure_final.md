# Monorepo Structure (Final) — Svelte 5 (SvelteKit) + Go Fiber

> Mục tiêu: **đơn giản – dễ hiểu – dễ quản lý**, nhưng vẫn đủ cho toàn bộ requirements:
- Editor có **Save button**, live preview realtime local
- Publish cache + CDN cache
- Custom domain: **1 domain = 1 page** (root `/`)
- Theme marketplace: preset **free/pro**, preset immutable, custom theme = patch
- Modes: light/dark/compact
- Password page Pro-only, remember 7 days
- Link style theo group (muốn link style riêng → tạo group riêng)
- GDPR tối thiểu: delete account/pages/assets

---

## 0) Tổng quan repo (Monorepo)
Repo chia thành 3 phần chính:
- `apps/web`: SvelteKit (Editor + Public renderer)
- `apps/api`: Go Fiber API
- `packages/shared`: types/schemas/contracts dùng chung (AI coding-friendly)

Ngoài ra:
- `docs`: tất cả spec/PRD/schema/migration
- `infra`: ghi chú cấu hình CDN/Cache
- `scripts`: tiện ích dev/deploy

```
linkbio/
├─ apps/
│  ├─ web/
│  └─ api/
├─ packages/
│  └─ shared/
├─ docs/
├─ infra/
├─ scripts/
└─ README.md
```

---

## 1) apps/web (SvelteKit)

### 1.1 Phân tách route (để public nhẹ, editor không kéo bundle nặng)
- `(public)`: PublicRenderer (render từ `compiled_json`)
- `(app)`: Dashboard + Editor + Appearance + Settings (auth-guarded)

```
apps/web/
├─ src/
│  ├─ routes/
│  │  ├─ (public)/
│  │  │  ├─ +layout.svelte
│  │  │  └─ [[...path]]/
│  │  │     ├─ +page.svelte
│  │  │     └─ +page.server.ts        # optional: proxy fetch/headers/etag
│  │  │
│  │  └─ (app)/
│  │     ├─ +layout.svelte            # app shell + auth guard
│  │     ├─ dashboard/
│  │     │  └─ +page.svelte
│  │     ├─ pages/
│  │     │  ├─ +page.svelte
│  │     │  └─ [pageId]/
│  │     │     ├─ edit/
│  │     │     │  ├─ +page.svelte
│  │     │     │  ├─ EditorShell.svelte
│  │     │     │  ├─ LeftPanel.svelte
│  │     │     │  ├─ PreviewPanel.svelte
│  │     │     │  └─ RightPanel.svelte
│  │     │     ├─ appearance/
│  │     │     │  └─ +page.svelte      # generate controls from meta.contract
│  │     │     └─ settings/
│  │     │        └─ +page.svelte      # header/cover/password (Pro)
│  │     ├─ themes/
│  │     │  ├─ +page.svelte            # marketplace free/pro
│  │     │  └─ [themeId]/
│  │     │     └─ +page.svelte
│  │     ├─ domains/
│  │     │  └─ +page.svelte            # add/verify/attach domain (Pro)
│  │     └─ billing/
│  │        └─ +page.svelte            # plan status
│  │
│  ├─ lib/
│  │  ├─ api/
│  │  │  ├─ client.ts                  # fetch wrapper + baseURL
│  │  │  ├─ auth.ts
│  │  │  ├─ pages.ts                   # getDraft/save/publish
│  │  │  ├─ themes.ts                  # presets/custom patch/apply
│  │  │  ├─ domains.ts
│  │  │  └─ assets.ts
│  │  │
│  │  ├─ renderer/                     # shared renderer (draft & public)
│  │  │  ├─ index.ts
│  │  │  ├─ ThemeResolver.ts
│  │  │  ├─ BlockRenderer.svelte
│  │  │  ├─ Header.svelte              # avatar/name/bio/social/verified/cover
│  │  │  ├─ blocks/
│  │  │  │  ├─ LinkGroupBlock.svelte
│  │  │  │  ├─ TextBlock.svelte
│  │  │  │  ├─ ImageBlock.svelte
│  │  │  │  ├─ ProductBlock.svelte
│  │  │  │  ├─ EmbedBlock.svelte
│  │  │  │  ├─ SocialRowBlock.svelte
│  │  │  │  ├─ SpacerBlock.svelte
│  │  │  │  └─ FormBlock.svelte
│  │  │  └─ primitives/                # button/card/text primitives theo theme
│  │  │
│  │  ├─ stores/
│  │  │  ├─ editorStore.svelte.ts      # Svelte 5 $state draft store
│  │  │  ├─ themeStore.svelte.ts
│  │  │  └─ userStore.svelte.ts
│  │  │
│  │  ├─ dnd/
│  │  │  ├─ sortKey.ts                 # fractional/lexo sort_key generator
│  │  │  └─ dndHelpers.ts
│  │  │
│  │  ├─ contracts/
│  │  │  └─ controls.ts                # UI generator for meta.contract.controls
│  │  │
│  │  └─ utils/
│  │     ├─ zod.ts
│  │     ├─ format.ts
│  │     └─ guards.ts
│  │
│  ├─ hooks.server.ts                  # auth guard (app routes), cookies
│  └─ app.css
│
├─ static/
├─ package.json
├─ svelte.config.js
└─ README.md
```

**Ghi chú quan trọng**
- `(public)` route chỉ dùng renderer + fetch compiled_json → bundle public cực nhẹ.
- `(app)` route chứa editor và appearance → code split tự nhiên.

---

## 2) apps/api (Go Fiber)

### 2.1 Tổ chức theo “module” (ít file nhưng rõ ràng)
Mỗi module giữ 3 lớp:
- `handler.go`: HTTP layer (Fiber routes)
- `service.go`: business logic
- `repo.go`: DB queries

```
apps/api/
├─ cmd/
│  └─ api/
│     └─ main.go                      # config, db, router, middleware
│
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ db/
│  │  ├─ pool.go                      # pgxpool init
│  │  └─ queries/                     # (optional) sql files / generated
│  ├─ middleware/
│  │  ├─ auth.go                      # cookie session auth (/api/*)
│  │  ├─ request_id.go
│  │  └─ recover.go
│  │
│  ├─ router/
│  │  └─ router.go                    # register all routes
│  │
│  ├─ modules/
│  │  ├─ auth/
│  │  │  ├─ handler.go                # register/login/oauth callbacks
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  ├─ billing/
│  │  │  ├─ handler.go                # plan status (webhook later)
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  ├─ themes/
│  │  │  ├─ handler.go                # presets, custom patch, apply
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  ├─ pages/
│  │  │  ├─ handler.go                # draft load/save, publish
│  │  │  ├─ service.go
│  │  │  ├─ compiler.go               # compile publish_cache (core performance)
│  │  │  └─ repo.go
│  │  ├─ domains/
│  │  │  ├─ handler.go                # add/attach domain (Pro)
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  ├─ assets/
│  │  │  ├─ handler.go                # upload-url, complete
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  ├─ access/
│  │  │  ├─ handler.go                # password verify, remember 7 days
│  │  │  ├─ service.go
│  │  │  └─ repo.go
│  │  └─ public/
│  │     ├─ handler.go                # GET /r (Host -> page -> compiled_json)
│  │     └─ service.go
│  │
│  └─ shared/
│     ├─ errors.go
│     ├─ response.go
│     └─ validate.go
│
├─ migrations/
│  └─ 0001_init.sql                   # link_in_bio_init_migration_updated.sql
│
├─ pkg/
│  ├─ hashing/                        # theme patch normalization + hash
│  └─ sortkey/                        # (optional) sort_key helpers
│
├─ go.mod
└─ README.md
```

---

## 3) packages/shared (Types + Schemas + Contracts dùng chung)

> Mục tiêu: FE/BE thống nhất payload, theme contract, block types → AI coding không bị lệch.

```
packages/shared/
├─ contracts/
│  ├─ api.ts                          # request/response types
│  ├─ page.ts                         # draft payload + compiled_json shape
│  ├─ blocks.ts                       # block type enums
│  └─ theme.ts                        # theme contract types (meta.contract.controls)
│
├─ schemas/
│  ├─ theme.zod.ts
│  ├─ pageDraft.zod.ts
│  └─ compiled.zod.ts
│
└─ README.md
```

---

## 4) docs (spec bắt buộc để code không trôi)

```
docs/
├─ link_in_bio_requirements_updated.md
├─ theme_system_spec_10_10_final.md
├─ link_in_bio_database_schema_updated.md
├─ frontend_backend_checklist.md
└─ api_contract.md                    # tạo thêm: payload save/publish/compiled_json
```

---

## 5) infra (CDN/cache/image) — tối giản

```
infra/
├─ cloudflare/
│  ├─ cache_rules.md                  # s-maxage, stale-while-revalidate, bypass rules
│  └─ purge_notes.md                  # purge keys on publish
└─ env.example
```

---

## 6) scripts (dev/deploy utilities)

```
scripts/
├─ dev.sh
├─ migrate.sh
└─ seed.sh                            # seed plans + theme presets (optional)
```

---

## 7) REVIEW & CHECKLIST (tự kiểm tra lần cuối)

### 7.1 Tính “đủ” theo requirements
- [x] Public renderer tách khỏi editor bằng route group `(public)`/`(app)`
- [x] Editor có `pages/[pageId]/edit` + `appearance` + `settings`
- [x] Renderer dùng chung (`lib/renderer`) cho preview + public
- [x] Drag-drop blocks + links có `dnd/sortKey.ts`
- [x] Theme deep editing dựa vào `contracts/controls.ts` đọc `meta.contract.controls`
- [x] Backend có module `pages/compiler.go` để compile publish_cache
- [x] Backend có module `public` để serve `/r` từ compiled_json (cacheable)
- [x] Backend có module `access` cho password session 7 ngày (Pro-only)
- [x] Backend có module `billing` để gate Free/Pro
- [x] `migrations/0001_init.sql` đúng theo schema updated
- [x] `packages/shared` giữ contracts + zod schemas để FE/BE sync

### 7.2 Tính “đơn giản, không phức tạp hóa”
- Không tách quá nhiều micro-packages
- Không lạm dụng layers: mỗi module Go chỉ 3 file handler/service/repo
- Web chỉ 1 renderer, reuse tối đa

---

## 8) Gợi ý 2 bước tiếp theo (để AI coding bắt đầu ngay)
1) Tạo `docs/api_contract.md` (chuẩn payload) cho:
   - `GET draft`, `POST save`, `POST publish`, `GET public compiled_json`
2) Tạo skeleton code:
   - SvelteKit routes + renderer stub
   - Fiber router + empty handlers + DB pool
