# Database Schema (Postgres) — Link-in-Bio Builder (Final)

> Mục tiêu schema: **đúng requirements**, dễ mở rộng, **public load siêu nhanh** nhờ publish cache, hỗ trợ **custom domain + nested path + redirect**, **draft/publish**, **theme preset immutable + custom patch**, **link style theo group**, **drag & drop** mượt bằng `sort_key`.

---

## 0) Quy ước & nguyên tắc thiết kế

### 0.1 Core rules (khớp requirements)
- **Một user có thể tạo nhiều Bio Page**.
- Mỗi page có **1 ngôn ngữ tại 1 thời điểm**: `locale = vi|en`.
- **Draft vs Published**: public page chỉ thay đổi khi user bấm **Publish**.
- **Custom domain**: user có thể gắn domain riêng (vd `yen.dev`).
- **Path** (thay slug): hỗ trợ `/` (root) và nested `/a/b`.
- Đổi path phải **redirect** path cũ → path mới.
- Theme preset **immutable**. Hễ appearance khác preset → tạo **custom theme** bằng **patch**.
- **Link style không theo từng link**. Muốn link style riêng → tạo **Link Group riêng** và style group đó.
- **Không có block-level override** (không style per-block). Style chỉ: **Theme → Group override**.
- Wallpaper background: **preset + upload** + effects blur/dim/overlayColor.
- Page có chế độ **password protected**.

### 0.2 Drag & drop
- Dùng `sort_key` (string) cho reorder:
  - `blocks.sort_key` (reorder block toàn page)
  - `links.sort_key` (reorder link trong group)
- Lợi: kéo thả chỉ update **1 row** mỗi lần.

### 0.3 Kiểu dữ liệu
- Khuyến nghị: `bigserial` cho id (hoặc UUID nếu bạn thích).
- Dùng `jsonb` cho:
  - theme config/patch
  - block content
  - group layout_config
  - group style_override (chỉ các key override)

---

## 1) User & Auth / Permission

## 1.5) Plans & Subscriptions (Free/Pro)
- `plans`: FREE/PRO
- `subscriptions`: user đang ở plan nào (dùng để gate custom domain + password page)


### 1.1 `users`
- Login: **Email + Password**
- OAuth: **Google, Facebook**

```sql
create table users (
  id bigserial primary key,

  email text unique,
  password_hash text, -- bcrypt/argon2id (nullable nếu chỉ OAuth)

  display_name text,
  avatar_asset_id bigint null,

  is_active boolean not null default true,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  -- basic sanity
  constraint chk_auth_method check (
    (email is not null) or (password_hash is not null)
  )
);
```

> Ghi chú: nếu bạn muốn bắt buộc email cho cả OAuth thì giữ `email not null`. Nếu muốn cho phép OAuth-only không email, bỏ unique email hoặc dùng `citext` + nullable carefully.

### 1.2 `oauth_accounts` (liên kết OAuth)
```sql
create table oauth_accounts (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,

  provider text not null,       -- 'google' | 'facebook'
  provider_user_id text not null,

  access_token text null,
  refresh_token text null,
  token_expires_at timestamptz null,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  unique(provider, provider_user_id),
  unique(user_id, provider)
);

create index idx_oauth_user on oauth_accounts(user_id);
```

### 1.3 Permission
- Giai đoạn này: **owner-only**.
- Mọi page thuộc về đúng `users.id` (không collaborator/team).

---

## 2) Domains & Routing (custom domain + path + redirect)

### 2.1 `domains`
- Quản lý domain hệ thống và custom domain.
- Mỗi `hostname` là unique.

```sql
create table domains (
  id bigserial primary key,
  user_id bigint null references users(id) on delete set null,

  hostname text not null unique,     -- 'domain.com', 'yen.dev'
  status text not null default 'pending', -- pending|active|disabled
  is_system boolean not null default false,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_domains_user on domains(user_id);
```

### 2.2 `page_routes` (path hiện tại + lịch sử redirect)
- Hỗ trợ root `/` và nested `/a/b`.
- **Chỉ path current** trong cùng domain được unique.

```sql
create table page_routes (
  id bigserial primary key,

  page_id bigint not null references bio_pages(id) on delete cascade,
  domain_id bigint not null references domains(id) on delete cascade,

  path text not null, -- '/', '/about', '/shop/item1'
  is_current boolean not null default true,

  redirect_to_route_id bigint null references page_routes(id),
  created_at timestamptz not null default now()
);

-- unique only for current routes
create unique index uq_domain_current_path
  on page_routes(domain_id, path)
  where is_current = true;

create index idx_route_lookup on page_routes(domain_id, path);
create index idx_route_current_by_page on page_routes(page_id, is_current);
```

**Chuẩn hoá path (app-level requirement)**
- luôn bắt đầu bằng `/`
- root là `/`
- không có trailing `/` (trừ `/`)
- lower-case nếu bạn muốn (khuyến nghị)

---

## 3) Pages (draft/publish + locale + password)

### 3.1 `bio_pages`
```sql
create table bio_pages (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,

  locale text not null default 'vi', -- vi|en
  title text null,

  status text not null default 'draft', -- draft|published

  -- password protected
  access_type text not null default 'public', -- public|password
  password_hash text null, -- bcrypt/argon2id
  password_updated_at timestamptz null,

  -- theme
  theme_preset_id bigint not null references theme_presets(id),
  theme_custom_id bigint null references themes_custom(id),
  theme_mode text not null default 'light', -- light|dark|compact

  -- page settings (header, cover, misc)
  settings jsonb not null default '{}'::jsonb,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint chk_locale check (locale in ('vi','en')),
  constraint chk_access_type check (access_type in ('public','password'))
);

create index idx_pages_user on bio_pages(user_id, created_at);
create index idx_pages_status on bio_pages(status, updated_at);
```

---

## 4) Theme system (preset immutable + custom patch + compiled)

### 4.1 `theme_presets` (5 theme hệ thống)
```sql
create table theme_presets (
  id bigserial primary key,
  key text not null unique,  -- 'theme_a', 'theme_b', ...
  name text not null,

  tier text not null default 'free',       -- free|pro
  visibility text not null default 'public', -- public|unlisted|private
  is_official boolean not null default true,
  author_user_id bigint null,

  -- full theme JSON (tokens/semantic/recipes/page/background/modes/meta)
  config jsonb not null,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
```

### 4.2 `themes_custom` (tạo khi appearance khác preset)
- Preset immutable; custom theme lưu **patch** (khác biệt).
- Có thể reuse cho nhiều page.

```sql
create table themes_custom (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,

  based_on_preset_id bigint not null references theme_presets(id),
  name text null,

  patch jsonb not null,           -- chỉ phần override
  compiled_config jsonb null,     -- merge sẵn (preset + patch) để render nhanh

  hash text not null,             -- dedupe: hash(presetId + patch normalized)

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  unique(user_id, hash)
);

create index idx_custom_themes_user on themes_custom(user_id);
```

**Workflow bắt buộc**
- Page đang dùng preset → user chỉnh appearance (khác preset) → tạo `themes_custom`, gán `bio_pages.theme_custom_id`.
- Các chỉnh sửa tiếp theo: update `patch` (và `compiled_config` nếu bạn lưu).

---

## 5) Assets (wallpaper preset + upload)

### 5.1 `assets`
- Dùng cho wallpaper/background, icon link, ảnh block, avatar…
- `scope` phân biệt preset hệ thống vs upload user.

```sql
create table assets (
  id bigserial primary key,

  user_id bigint null references users(id) on delete set null,
  scope text not null, -- system_preset|user_upload
  type text not null,  -- image (future: video)

  provider text not null, -- s3|r2|cf_images|...
  storage_key text not null,
  url text null,

  mime_type text null,
  size_bytes bigint null,
  width int null,
  height int null,

  created_at timestamptz not null default now(),

  constraint chk_asset_scope check (scope in ('system_preset','user_upload'))
);

create index idx_assets_user on assets(user_id, created_at);
```

---

## 6) Blocks (layout 1 cột + drag-drop)

### 6.1 `blocks`
- Một page là danh sách block 1 cột, sort theo `sort_key`.
- Không có style_override per-block.
- `type='link_group'` thì `ref_id` trỏ tới `link_groups.id`.
- Các block khác lưu data trong `content` (JSONB).

```sql
create table blocks (
  id bigserial primary key,
  page_id bigint not null references bio_pages(id) on delete cascade,

  type text not null,        -- link_group|text|image|product|spacer|...
  sort_key text not null,    -- LexoRank/fractional key

  ref_id bigint null,        -- when type=link_group => link_groups.id
  content jsonb not null default '{}'::jsonb,

  is_visible boolean not null default true,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_blocks_page_sort on blocks(page_id, sort_key);
create index idx_blocks_page_type on blocks(page_id, type);
```

**Block content (gợi ý)**
- text:
  - `{ "text": "...", "variant": "body|heading|caption" }`
- image:
  - `{ "assetId": 123, "alt": "..." }`
- product (GĐ1):
  - `{ "title": "...", "price": 199000, "currency": "VND", "url": "...", "image_asset_id": 123 }`

---

## 7) Link Groups & Links (style theo group)

### 7.1 `link_groups`
- Group có layout: list/cards/grid.
- Group có `style_override` để override theme (chỉ trong group).
- Không style theo từng link.

```sql
create table link_groups (
  id bigserial primary key,
  page_id bigint not null references bio_pages(id) on delete cascade,

  title text null,

  layout_type text not null default 'list', -- list|cards|grid
  layout_config jsonb not null default '{}'::jsonb, -- {columns,gap,...}

  style_override jsonb null, -- override theme defaults for this group

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint chk_layout_type check (layout_type in ('list','cards','grid'))
);

create index idx_groups_page on link_groups(page_id);
```

### 7.2 `links`
- Links nằm trong group, reorder bằng `sort_key`.

```sql
create table links (
  id bigserial primary key,
  group_id bigint not null references link_groups(id) on delete cascade,

  title text not null,
  url text not null,

  icon_asset_id bigint null references assets(id),

  sort_key text not null,
  is_active boolean not null default true,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_links_group_sort on links(group_id, sort_key);
create index idx_links_group_active on links(group_id, is_active);
```

**Rule quan trọng**
- Muốn một link style khác → tạo **group riêng** và set `link_groups.style_override`.

---

## 8) Publish cache (public load siêu nhanh)

### 8.1 `page_publish_cache`
- Khi Publish, compile page thành JSON đã merge theme + group override + blocks order.

```sql
create table page_publish_cache (
  page_id bigint primary key references bio_pages(id) on delete cascade,

  compiled_json jsonb not null,

  published_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
```

**Public request tối ưu**
1) resolve domain+path → tìm `page_routes` current
2) check `bio_pages.access_type`
3) trả `page_publish_cache.compiled_json`

> compiled_json nên embed sẵn blocks + group + links sorted và style đã merge để FE render 1 vòng.

---

## 9) Password-protected access (khuyến nghị triển khai)
Schema ở trên đủ cho password hash; bạn có thể chọn thêm session table (tuỳ kiến trúc).

### Option A (nhẹ DB): JWT
- Không cần table.
- Verify password → phát JWT scoped `page_id`.

### Option B (stateful): `page_access_sessions`
```sql
create table page_access_sessions (
  id bigserial primary key,
  page_id bigint not null references bio_pages(id) on delete cascade,

  token_hash text not null,
  expires_at timestamptz not null,

  created_at timestamptz not null default now(),

  unique(page_id, token_hash)
);

create index idx_page_access_exp on page_access_sessions(page_id, expires_at);
```

---

## 10) Analytics (để sau)
Không bắt buộc hiện tại. Khi cần, khuyến nghị dùng bảng aggregate theo ngày để nhẹ:
- `page_daily_stats(page_id, date, views, ...)`
- `link_daily_stats(link_id, date, clicks, ...)`

---

## 11) Ràng buộc & tính nhất quán (khuyến nghị)
- App-level validation:
  - path canonicalization (như mục 2.2)
  - url validation cho links
  - JSON schema validation cho theme preset/custom (theo `schemaVersion`)
- Khi xoá page:
  - cascade xoá routes, blocks, groups, links, publish_cache
- Khi publish:
  - compile theme (preset + patch + mode nếu có) + merge group overrides

---

## 12) Tóm tắt “đủ khớp yêu cầu”
- [x] Auth: Email+Password + OAuth (Google/Facebook)
- [x] User tạo nhiều page
- [x] Locale per page (vi/en)
- [x] Draft/Publish (public đổi khi Publish)
- [x] Custom domain + root `/` + nested path
- [x] Đổi path có redirect history
- [x] 5 theme preset immutable
- [x] Appearance chỉnh khác preset → tạo custom theme patch, reuse nhiều page
- [x] Blocks + drag-drop bằng sort_key
- [x] Link theo group, style theo group (không per-link style)
- [x] Wallpaper preset + upload + effects
- [x] Password protected page
- [x] Public load nhanh bằng page_publish_cache

---

## Plans & Subscriptions (gợi ý DDL)
```sql
create table plans (
  id bigserial primary key,
  code text not null unique, -- FREE|PRO
  name text not null,
  created_at timestamptz not null default now()
);

create table subscriptions (
  id bigserial primary key,
  user_id bigint not null references users(id) on delete cascade,
  plan_id bigint not null references plans(id),
  status text not null, -- active|canceled|past_due
  current_period_end timestamptz null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_subscriptions_user on subscriptions(user_id);
```
