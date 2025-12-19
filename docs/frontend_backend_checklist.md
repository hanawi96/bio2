# Frontend + Backend Checklist (Svelte 5 + Go Fiber) — Siêu nhanh / Siêu mượt / Siêu nhẹ / Siêu ổn định

> Bản tổng hợp trực quan “cần dùng những gì” cho hệ Link-in-Bio đã chốt requirements:
- Save theo nút (không undo/redo)
- Live preview realtime **local** trong editor
- Publish cache + CDN cache
- Custom domain: **1 domain = 1 page**, dùng ở `/`
- Theme marketplace: preset **free/pro**, preset immutable, custom theme bằng patch
- Modes: light/dark/compact
- Password page Pro-only, nhớ 7 ngày
- Link style theo **group**, muốn link style riêng → tạo group riêng

---

## A) Frontend (Svelte 5)

### A1) 2 bề mặt app
1) **Editor App** (`/app`)
- Dashboard: pages, themes, domains, billing (Pro)
- Page editor: blocks + groups + links + drag-drop
- Appearance editor: chỉnh theme “rất sâu” theo `meta.contract.controls`
- Live preview mockup realtime

2) **Public Renderer**
- Render từ `compiled_json` (đã compile khi publish)
- JS tối giản, không dính editor deps

---

### A2) State & data model (để mượt)
- Dùng Svelte 5 `$state` cho draft:
  - `page.settings` (header + cover)
  - `blocks[]` (sort_key)
  - `linkGroups{}` + `links{}` (sort_key)
  - `themeSelection` (preset/custom + mode)
  - `themePatchDraft` (appearance patch)
- Normalize dữ liệu (maps byId) để:
  - update O(1)
  - drag-drop không re-render toàn trang

---

### A3) Renderer dùng chung (đúng “What you see is what you publish”)
- Tạo 1 bộ renderer component:
  - **DraftRenderer** (editor preview): render từ draft state local
  - **PublicRenderer** (public): render từ compiled_json
- Input chuẩn (khuyến nghị):
  - `{ settings, theme, blocks, groups, links }`

---

### A4) Drag & Drop (blocks + links)
- Thư viện: `@dnd-kit` (khuyên dùng)
- Reorder bằng `sort_key` (LexoRank/fractional)
  - mỗi lần kéo thả chỉ cần đổi `sort_key` của item bị move
- Persist khi **Save**

---

### A5) Appearance / Theme deep editing (vũ khí cạnh tranh)
- FE tự generate UI chỉnh theme từ:
  - `theme.meta.contract.controls[]`
  - `constraints`
- Output:
  - **patch JSON** (diff so với preset)
- Rule:
  - thay đổi khác preset → server tạo/update custom theme patch (preset immutable)

---

### A6) Public rendering tối ưu
- Public chỉ:
  - fetch compiled_json (hoặc server inject)
  - render
  - lazy-load ảnh
- Không tải editor code

---

### A7) FE performance checklist
- Code-splitting:
  - editor routes chunk riêng
  - appearance/theme editor chunk riêng
- Ảnh:
  - responsive images, webp/avif
- Render:
  - key theo id
  - tránh props “to” (truyền object lớn) nếu không cần
- Networking:
  - Save theo nút → ít write, ít request

---

### A8) Libs tối thiểu (gọn)
- fetch wrapper (nhẹ)
- `zod` (validate payload + contract)
- `@dnd-kit` (drag-drop)
- toast/snackbar nhẹ

---

## B) Backend (Go Fiber)

### B1) Modules/services nên có
1) **Auth**
- Email + Password
- OAuth (Google, Facebook)

2) **Entitlements (Free/Pro)**
- Check plan/subscription
- Gate: custom domain, password page, pro themes

3) **Pages (Draft)**
- load draft / save draft (Save button)

4) **Publisher**
- compile & publish cache
- purge CDN

5) **Public**
- resolve host → page
- trả compiled_json
- password verify + session 7d (Pro-only)

6) **Themes**
- presets marketplace (free/pro)
- custom themes patch + compiled_config

7) **Assets**
- upload signed URL
- persist metadata

8) **Domains (Pro)**
- add domain
- attach domain → page (1 domain = 1 page)

---

### B2) Endpoints tối thiểu (gợi ý)
**Public**
- `GET /r` (đọc Host, resolve page, trả compiled_json)
- `POST /r/password` (verify password, set cookie 7d)

**Editor**
- `GET /api/pages/:id/draft`
- `POST /api/pages/:id/save`
- `POST /api/pages/:id/publish`

**Themes**
- `GET /api/themes/presets?tier=free|pro`
- `POST /api/themes/custom` (create/update patch)
- `POST /api/themes/apply`

**Domains (Pro)**
- `POST /api/domains`
- `POST /api/domains/:id/attach-page`

**Assets**
- `POST /api/assets/upload-url`
- `POST /api/assets/complete`

---

### B3) Publish compile (trung tâm performance)
Khi Publish:
- Load page + settings + theme (preset + patch + mode)
- Load blocks (sort_key)
- Với link_group:
  - load group + links (sort_key)
  - merge style: `theme.page.defaults.linkGroup` + `group.style_override`
- Build `compiled_json` “ready to render”
- Save `page_publish_cache`
- Purge CDN theo hostname

Public request:
- resolve domain → page_id (custom domain 1 page)
- read `page_publish_cache`
- trả response + cache headers

---

### B4) DB access nhanh
- Postgres + `pgxpool`
- Ít round-trips (1–2 query cho public)
- Prepared statements + pool tuning

---

### B5) Caching (DB + CDN)
- DB: `page_publish_cache` là source of truth
- CDN cache theo hostname
- Trả headers:
  - `Cache-Control` (s-maxage, stale-while-revalidate)
  - `ETag` (hash compiled_json)
- Publish → purge CDN key

---

### B6) Password page “remember 7 days” (Pro-only)
- `page_access_sessions`:
  - `token_hash`, `expires_at`
- Cookie httpOnly/secure/sameSite
- Verify session on public request

---

### B7) GDPR tối thiểu (bắt buộc)
- Endpoint: `DELETE /api/account`
- DB cascade delete:
  - pages, blocks, groups, links, sessions, publish_cache
  - assets (user_upload) theo policy
- Storage deletion:
  - xoá object trong storage theo batch/job (khuyến nghị)

---

### B8) Stability/observability (để “siêu ổn định”)
- Structured logs (JSON)
- Metrics: latency p50/p95, db time, publish time
- Alert basic (5xx rate, db connection issues)

---

## C) Infra / CDN / Images

### C1) CDN (bắt buộc)
- Cache public responses
- Purge khi publish
- HTTPS + redirect http→https do CDN/infra xử lý

### C2) Image pipeline (để nhẹ)
- Khuyên dùng: Cloudflare Images (resize/format/edge cache)
- Alternative: S3/R2 + presigned + worker resize (phức tạp hơn)

### C3) Deploy (để low latency Bangkok/SEA)
- Origin ở Singapore
- Static FE assets trên CDN

---

## D) Chuẩn dữ liệu cần “đóng”
1) **Draft payload**: normalized maps + sort_key
2) **Theme patch**: nested JSON patch
3) **Compiled JSON**: ready-to-render, embed groupFinalStyle
4) **Theme contract**: controls/constraints/keyPaths để auto-generate Appearance UI

---

## E) Tóm tắt “siêu nhanh, siêu mượt, siêu nhẹ”
- Preview: 100% local render trong editor (instant)
- Public: publish cache + CDN cache (fastest)
- Complexity thấp: không undo/redo, không per-link style, save theo nút
- Theme cạnh tranh: contract-driven + marketplace + modes

