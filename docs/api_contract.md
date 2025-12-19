# API Contract (Final) — Link-in-Bio (Svelte 5 + Go Fiber)

> Mục tiêu: để AI implement **đúng logic**, FE/BE **khớp payload**, public **siêu nhanh** nhờ publish cache + CDN.

## 0) Quy ước chung

### 0.1 Base URL
- Dev: `http://localhost:8080`
- API prefix (private/editor): `/api`
- Public renderer endpoint: `/r` (không prefix `/api`)

### 0.2 Auth & session
- Dashboard/editor dùng **cookie session** (httpOnly).
- Public password page dùng **cookie token** (httpOnly) để “remember 7 days”.

**Cookie names (khuyến nghị)**
- `sid` — editor session (login)
- `pats` — public access token (password page session), scoped per page

### 0.3 Plans (Free/Pro) — Feature gating (enforced in application layer)
- Custom domain: **Pro-only**
- Password page: **Pro-only**
- Pro theme presets: **Pro-only** (Free chỉ xem/không apply)

### 0.4 Draft vs Published
- Editor thao tác trên **draft**.
- Public chỉ đọc **published compiled_json**.
- Public chỉ thay đổi khi user bấm **Publish**.

### 0.5 Locale
- Mỗi page chỉ có 1 locale tại 1 thời điểm: `vi|en`.

### 0.6 Sort / Drag & Drop
- Reorder dùng `sort_key` (string, fractional/lexorank).
- `blocks.sort_key`: thứ tự block trên page.
- `links.sort_key`: thứ tự link trong group.

### 0.7 Error format (bắt buộc thống nhất)
Mọi lỗi trả về JSON:

```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "Human readable message",
    "details": { }
  }
}
```

### 0.8 Standard response headers (public)
- `ETag: "<hash>"`
- `Cache-Control: public, max-age=60, s-maxage=3600, stale-while-revalidate=86400`

> Publish sẽ purge CDN theo hostname.

---

## 1) Public Rendering

### 1.1 GET `/r`
**Mục đích**: trả nội dung public của page theo `Host` (custom domain 1 page) hoặc theo system domain + path (nếu bạn dùng system domain).

**Request**
- Dựa vào HTTP `Host` header
- Path: với custom domain luôn `/` (app/infra enforce), nên `/r` có thể chạy ở mọi path.
- Query (dev helper, optional):
  - `?host=yen.dev&path=/` (chỉ dùng dev nếu không set Host)

**Response 200**
- `compiled_json` (schema ở mục 6.2)

**Response 401 (password required)**
```json
{
  "error": {
    "code": "PASSWORD_REQUIRED",
    "message": "This page is password protected"
  }
}
```

**Response 404**
```json
{
  "error": { "code": "NOT_FOUND", "message": "Page not found" }
}
```

---

### 1.2 POST `/r/password`
**Mục đích**: verify password của page, tạo session 7 ngày.

**Request**
```json
{
  "password": "string"
}
```

**Response 200**
- Set cookie `pats` (httpOnly, secure, 7 days)
```json
{ "ok": true }
```

**Response 401**
```json
{
  "error": { "code": "INVALID_PASSWORD", "message": "Invalid password" }
}
```

---

## 2) Auth (Editor)

### 2.1 POST `/api/auth/register`
```json
{
  "email": "user@example.com",
  "password": "string",
  "display_name": "Yen"
}
```
**Response 200**
```json
{ "user": { "id": 1, "email": "user@example.com", "plan": "FREE" } }
```

### 2.2 POST `/api/auth/login`
```json
{
  "email": "user@example.com",
  "password": "string"
}
```
**Response 200**
- Set cookie `sid`
```json
{ "user": { "id": 1, "email": "user@example.com", "plan": "FREE" } }
```

### 2.3 POST `/api/auth/logout`
**Response 200**
```json
{ "ok": true }
```

### 2.4 OAuth (Google/Facebook)
- `/api/auth/oauth/:provider/start`
- `/api/auth/oauth/:provider/callback`
> Payload phụ thuộc provider; tối thiểu yêu cầu tạo user + set cookie `sid`.

---

## 3) Billing / Entitlements

### 3.1 GET `/api/billing/me`
**Response 200**
```json
{
  "plan": "FREE",
  "status": "active",
  "features": {
    "custom_domain": false,
    "password_page": false,
    "pro_themes": false
  }
}
```

---

## 4) Pages (Draft workflow)

### 4.1 POST `/api/pages`
**Mục đích**: tạo page mới.

**Request**
```json
{
  "title": "My Bio",
  "locale": "vi",
  "theme_preset_key": "theme_a"
}
```

**Response 200**
```json
{ "page_id": 101 }
```

---

### 4.2 GET `/api/pages`
**Response 200**
```json
{
  "pages": [
    {
      "id": 101,
      "title": "My Bio",
      "locale": "vi",
      "status": "draft",
      "updated_at": "2025-12-19T00:00:00Z"
    }
  ]
}
```

---

### 4.3 GET `/api/pages/:id/draft`
**Mục đích**: load draft state cho editor.

**Response 200**
- Draft schema ở mục 6.1.

---

### 4.4 POST `/api/pages/:id/save`
**Mục đích**: Save theo nút, ghi draft state.

**Request**
- Draft payload (schema 6.1)

**Response 200**
```json
{ "ok": true, "updated_at": "2025-12-19T00:00:00Z" }
```

---

### 4.5 POST `/api/pages/:id/publish`
**Mục đích**: compile + ghi `page_publish_cache` + purge CDN.

**Request**
```json
{ "note": "optional publish note" }
```

**Response 200**
```json
{ "ok": true, "published_at": "2025-12-19T00:00:00Z" }
```

---

### 4.6 PATCH `/api/pages/:id/access`
**Mục đích**: set public/password (Pro-only).

**Request**
```json
{
  "access_type": "public",
  "password": "optional when access_type=password"
}
```

**Rules**
- Nếu `access_type=password`:
  - require Pro
  - hash password và lưu vào page
  - reset all existing `page_access_sessions`

**Response 200**
```json
{ "ok": true }
```

---

### 4.7 DELETE `/api/account`
**Mục đích**: GDPR tối thiểu (delete account, pages, assets).
**Response 200**
```json
{ "ok": true }
```

---

## 5) Themes (Marketplace + Custom patch)

### 5.1 GET `/api/themes/presets?tier=free|pro`
**Response 200**
```json
{
  "themes": [
    {
      "id": 1,
      "key": "theme_a",
      "name": "Theme A",
      "tier": "free",
      "preview": { "thumbnail_url": "..." }
    }
  ]
}
```

### 5.2 POST `/api/themes/custom`
**Mục đích**
- tạo hoặc update custom theme patch cho user
- preset immutable, patch khác preset → custom theme

**Request**
```json
{
  "based_on_preset_key": "theme_a",
  "name": "My Custom A",
  "patch": { }
}
```

**Response 200**
```json
{ "theme_custom_id": 9001 }
```

### 5.3 POST `/api/themes/apply`
**Mục đích**: apply preset/custom + mode vào page.

**Request**
```json
{
  "page_id": 101,
  "theme_preset_key": "theme_a",
  "theme_custom_id": 9001,
  "theme_mode": "dark"
}
```

**Rules**
- Nếu preset `tier=pro` thì require Pro.
- `theme_mode`: `light|dark|compact`.

**Response 200**
```json
{ "ok": true }
```

---

## 6) Draft / Compiled Payload Schemas (bắt buộc)

### 6.1 Draft payload (editor)
> Normalize theo id để update nhanh, nhưng vẫn dễ serialize.

```json
{
  "page": {
    "id": 101,
    "title": "My Bio",
    "locale": "vi",
    "status": "draft",
    "theme": {
      "preset_key": "theme_a",
      "custom_id": 9001,
      "mode": "light"
    },
    "access_type": "public",
    "settings": {
      "header": {
        "avatar_asset_id": 100,
        "name": "Yen",
        "bio": "Hello",
        "verified": true,
        "social": [
          { "type": "facebook", "url": "https://facebook.com/..." }
        ]
      },
      "cover": { "kind": "color", "color": "#111111", "asset_id": null }
    }
  },

  "blocks": [
    { "id": 1, "type": "link_group", "sort_key": "a0", "ref_id": 10, "content": {}, "is_visible": true },
    { "id": 2, "type": "text", "sort_key": "a1", "ref_id": null, "content": { "text": "About me", "variant": "heading" }, "is_visible": true }
  ],

  "link_groups": {
    "10": {
      "id": 10,
      "title": "Links",
      "layout_type": "list",
      "layout_config": {},
      "style_override": { "textAlign": "left", "fontSize": "XL" }
    }
  },

  "links": {
    "10": [
      { "id": 1001, "title": "Instagram", "url": "https://instagram.com/...", "icon_asset_id": null, "sort_key": "a0", "is_active": true },
      { "id": 1002, "title": "YouTube", "url": "https://youtube.com/...", "icon_asset_id": null, "sort_key": "a1", "is_active": true }
    ]
  }
}
```

**Rules**
- Không có style theo từng link.
- Không có block-level override.
- Muốn link style riêng → tạo group riêng + style_override group.

---

### 6.2 Compiled payload (public)
> Đây là payload public trả về từ `/r` (đã merge theme + group override + sort sẵn).

```json
{
  "version": 1,
  "page": {
    "id": 101,
    "locale": "vi",
    "settings": { }
  },

  "theme": {
    "preset_key": "theme_a",
    "mode": "light",
    "compiled": { }
  },

  "blocks": [
    {
      "type": "link_group",
      "group": {
        "id": 10,
        "title": "Links",
        "layout_type": "list",
        "layout_config": {},
        "final_style": { },
        "links": [
          { "title": "Instagram", "url": "https://instagram.com/...", "icon_url": null }
        ]
      }
    },
    {
      "type": "text",
      "content": { "text": "About me", "variant": "heading" }
    }
  ]
}
```

**Compiled rules**
- `blocks` đã sort theo `sort_key`.
- `links` trong group đã sort theo `sort_key`.
- `final_style` đã merge `theme.page.defaults.linkGroup` + `group.style_override`.

---

## 7) Assets

### 7.1 POST `/api/assets/upload-url`
**Request**
```json
{ "type": "image", "mime_type": "image/png", "size_bytes": 123456 }
```
**Response 200**
```json
{
  "upload_url": "https://...",
  "asset_temp_key": "temp_..."
}
```

### 7.2 POST `/api/assets/complete`
**Request**
```json
{
  "asset_temp_key": "temp_...",
  "width": 800,
  "height": 600
}
```
**Response 200**
```json
{ "asset_id": 555, "url": "https://cdn/..." }
```

---

## 8) Domains (Pro-only)

### 8.1 POST `/api/domains`
**Request**
```json
{ "hostname": "yen.dev" }
```
**Rules**
- require Pro

**Response 200**
```json
{ "domain_id": 77, "status": "pending" }
```

### 8.2 POST `/api/domains/:id/attach-page`
**Request**
```json
{ "page_id": 101 }
```
**Rules**
- require Pro
- custom domain: 1 domain = 1 page
- attach sẽ tạo/đảm bảo route current path = `/`

**Response 200**
```json
{ "ok": true }
```
