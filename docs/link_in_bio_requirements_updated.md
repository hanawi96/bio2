# Tài liệu yêu cầu (PRD + Data Model) — Link-in-Bio Builder

> Mục tiêu: xây dựng hệ thống tạo “Link in Bio” cho user với editor kéo-thả mượt, theme preset + custom theme, hỗ trợ custom domain, public load cực nhanh (cache publish), và sẵn sàng mở rộng.

---

## 1) Tổng quan sản phẩm

### 1.1 Luồng chính
1. User đăng ký tài khoản thành công.
2. User tạo một Bio Page và chọn đường dẫn (path) public, ví dụ:
   - Domain hệ thống: `domain.com/yendev96` (tương đương `domain.com/{path}`)
   - Custom domain: `yen.dev/` hoặc `yen.dev/about`
3. Chọn **theme preset** (hệ thống có sẵn **5 theme preset**).
4. Vào trang **thiết kế bio (editor)** để thêm/sắp xếp nội dung.
5. Vào trang **appearance** để tuỳ chỉnh giao diện (theme).
6. Khi mọi thứ ổn: user bấm **Publish** → public page mới cập nhật.

### 1.2 Các loại nội dung (Block)
- Hệ thống hỗ trợ nhiều loại block: `link_group`, `text`, `image`, `product`, `spacer`, … (có thể mở rộng).
- **Link** được tổ chức theo **Group** (link group). Link không style theo từng link, mà **thừa kế style từ group**.

### 1.3 Kéo-thả (Drag & Drop)
- User có thể kéo-thả:
  - Các block trên page để đổi thứ tự hiển thị (chuyên nghiệp).
  - Các link trong từng link group để đổi thứ tự.

---

## 2) Quy tắc kế thừa style (Theme → Group)

### 2.1 Quy tắc ưu tiên
1. **Theme** quyết định style mặc định cho cả page (các token như text align, font size, border, spacing, …).
2. **Group override**: nếu group có cấu hình riêng cho một thuộc tính nào đó, thì **thuộc tính đó dùng giá trị của group**, không bị ảnh hưởng bởi theme.
3. (Tuỳ chọn tương lai) **Block override**: nếu sau này cho phép block text/image có style riêng.

### 2.2 Ví dụ đúng theo yêu cầu
- Theme A: `textAlign=left`, `fontSize=M`
- Group link A override: `textAlign=center`, `fontSize=XL`
- Khi user vào appearance đổi theme thành `textAlign=right`, `fontSize=L`:
  - Các group không override → theo theme mới (right + L)
  - Group A vẫn giữ override (center + XL)

---

## 3) Đa ngôn ngữ (Locale)
- Mỗi Bio Page hiển thị **chỉ 1 ngôn ngữ tại 1 thời điểm**.
- Page có trường `locale`: `vi` hoặc `en`.
- Không yêu cầu switch ngôn ngữ trên public page.

---

## 4) Domain / Path / Redirect

### 4.1 Custom domain
- User có thể dùng **custom domain** của riêng mình trong tương lai, ví dụ `yen.dev`.

### 4.2 Path (thay cho slug)
- Hỗ trợ:
  - Root: `yen.dev/`
  - Nested path: `yen.dev/about`, `yen.dev/shop/item1`
- Path được chuẩn hoá:
  - Luôn bắt đầu bằng `/`
  - Root là `/`
  - Không có dấu `/` ở cuối (trừ root)

### 4.3 Đổi path và redirect
- User có thể đổi path.
- Khi đổi path, hệ thống phải **redirect path cũ → path mới**.

---

## 5) Draft / Publish / Live preview trong editor

### 5.1 Draft vs Published
- Có 2 trạng thái:
  - **Draft**: user chỉnh sửa trong editor, lưu nháp.
  - **Published**: public page chỉ thay đổi khi user bấm **Publish**.

### 5.2 Mockup/live preview trong editor & appearance
- Trong màn thiết kế bio và appearance, user cần **xem preview ngay** (mockup realtime).
- Không bắt buộc phải có “public preview link” riêng; preview realtime có thể render trực tiếp từ state editor.

### 5.3 Cache cho public load siêu nhanh
- Khi user bấm Publish, backend sẽ **compile** toàn bộ page thành 1 JSON “đã sắp xếp + đã merge style” và lưu cache.
- Public request chỉ cần đọc cache để render nhanh.

---

## 6) Wallpaper background (preset + upload) + effects
- Wallpaper background hỗ trợ:
  - **Preset image** (hệ thống)
  - **Upload image** (user)
- Có hiệu ứng (effects):
  - `blur`
  - `dim`
  - `overlayColor` (màu phủ)

---

## 7) Private / Password page
- Page có chế độ **private/password**:
  - `public`: ai cũng xem được
  - `password`: cần nhập mật khẩu để xem
- Lưu password dưới dạng **hash** (bcrypt/argon2id).

---

## 8) Analytics
- Analytics **để sau** (không implement ngay).
- Thiết kế nên “ready” nhưng không bắt buộc ghi event hiện tại.

---

## 9) Product block — phương án tối ưu khuyến nghị
### Giai đoạn 1 (khuyến nghị)
- `product` là “product card link” nhập tay, lưu trong `blocks.content`:
  - `title`, `price`, `currency`, `url`, `image_asset_id`
- Ưu điểm: đơn giản, nhẹ, render nhanh, phù hợp MVP.

### Giai đoạn 2 (mở rộng sau)
- Thêm `source` (`manual|catalog`) + `product_id` (nếu có catalog) mà không phá cấu trúc hiện tại.

---

## 10) Không giới hạn blocks/links
- Yêu cầu nghiệp vụ: “không giới hạn”.
- Khuyến nghị kỹ thuật: nên có **soft limit** (theo plan hoặc theo kích thước compiled_json) để chống abuse và giữ performance.

---

# 11) Data Model / Database Schema (đề xuất)

> DB khuyến nghị: **Postgres** (RDBMS) + **JSONB** cho style/config để linh hoạt.

## 11.1 Bảng `domains`
Quản lý domain hệ thống và custom domain.

- `id`
- `user_id` (nullable; domain hệ thống có thể null)
- `hostname` (unique)
- `status` (`pending|active|disabled`)
- `is_system` (bool)
- timestamps

**Index/Constraint**
- `UNIQUE(hostname)`
- `INDEX(user_id)`

---

## 11.2 Bảng `bio_pages`
Thông tin page.

- `id`
- `user_id`
- `locale` (`vi|en`)
- `title`
- `theme_preset_id`
- `theme_custom_id` (nullable)
- `status` (`draft|published`)
- `access_type` (`public|password`)
- `password_hash` (nullable)
- timestamps

**Index**
- `INDEX(user_id, created_at)`

---

## 11.3 Bảng `page_routes` (path + redirect history)
Lưu path hiện tại và lịch sử để redirect.

- `id`
- `page_id` (FK bio_pages)
- `domain_id` (FK domains)
- `path` (text) — `/`, `/about`, `/shop/item1`
- `is_current` (bool)
- `redirect_to_route_id` (self-FK, nullable)
- `created_at`

**Index/Constraint quan trọng**
- Partial unique: `UNIQUE(domain_id, path) WHERE is_current=true`
- `INDEX(domain_id, path)` (lookup public)
- `INDEX(page_id, is_current)`

---

## 11.4 Bảng `theme_presets` (5 theme hệ thống)
- `id`
- `key` (unique)
- `name`
- `config` (JSONB full)
- timestamps

---

## 11.5 Bảng `themes_custom` (theme custom tái sử dụng)
- `id`
- `user_id`
- `based_on_preset_id`
- `name`
- `patch` (JSONB — chỉ phần override)
- `compiled_config` (JSONB — merge sẵn để nhanh)
- `hash` (dedupe)
- timestamps

**Index/Constraint**
- `INDEX(user_id)`
- `UNIQUE(user_id, hash)`

---

## 11.6 Bảng `assets` (wallpaper preset + upload)
- `id`
- `user_id` (nullable)
- `scope` (`system_preset|user_upload`)
- `type` (`image`)
- `provider` (`s3|r2|cf_images`)
- `storage_key`
- `url` (optional)
- `mime_type`, `size_bytes`, `width`, `height`
- timestamps

**Index**
- `INDEX(user_id, created_at)`

---

## 11.7 Bảng `blocks` (layout 1 cột + drag-drop)
- `id`
- `page_id`
- `type` (`link_group|text|image|product|spacer|...`)
- `sort_key` (string — phục vụ reorder mượt, update 1 row)
- `ref_id` (nullable; nếu `type=link_group` thì trỏ `link_groups.id`)
- `content` (JSONB)
- `style_override` (JSONB, optional)
- `is_visible` (bool)
- timestamps

**Index**
- `INDEX(page_id, sort_key)`
- `INDEX(page_id, type)`

---

## 11.8 Bảng `link_groups`
- `id`
- `page_id`
- `title`
- `layout_type` (`list|cards|grid`)
- `layout_config` (JSONB: columns, gap, …)
- `style_override` (JSONB: override theme keys cho group)
- timestamps

**Index**
- `INDEX(page_id)`

---

## 11.9 Bảng `links`
- `id`
- `group_id`
- `title`
- `url`
- `icon_asset_id` (nullable)
- `sort_key` (string)
- `is_active` (bool)
- timestamps

**Index**
- `INDEX(group_id, sort_key)`
- `INDEX(group_id, is_active)`

---

## 11.10 Bảng `page_publish_cache` (public load cực nhanh)
- `page_id` (PK/FK)
- `compiled_json` (JSONB)
- `published_at`, `updated_at`

**Mục tiêu**
- Public request chỉ cần resolve route + đọc compiled_json.

---

# 12) Cấu trúc JSON style (khuyến nghị)

## 12.1 Theme config (token-based)
Ví dụ:

```json
{
  "typography": { "fontFamily": "Inter", "baseSize": "M", "textAlign": "left" },
  "spacing": { "pagePadding": 16, "blockGap": 12 },
  "shape": { "borderRadius": 16, "borderWidth": 1 },
  "background": {
    "wallpaper": { "kind": "preset", "asset_id": 123 },
    "effects": { "blur": 8, "dim": 0.3, "overlayColor": "rgba(0,0,0,0.25)" }
  },
  "defaults": {
    "linkGroup": { "textAlign": "center", "fontSize": "M" }
  }
}
```

## 12.2 Group override (chỉ lưu key cần override)
```json
{ "textAlign": "center", "fontSize": "XL" }
```

---

# 13) Compile/PUBLISH output (để FE render nhẹ)

## 13.1 Yêu cầu compiled_json
- Đã chọn theme cuối (preset + patch hoặc compiled_config)
- Blocks đã sắp xếp theo sort_key
- Mỗi link_group block embed:
  - group data + `finalStyle` = merge(theme defaults → group override)
  - links sorted theo sort_key
- FE render chỉ cần map 1 vòng, không cần gọi API con.

---

# 14) Quy tắc kỹ thuật để “siêu mượt”
- Drag & drop dùng `sort_key` (LexoRank/fractional indexing) để reorder không phải update hàng loạt.
- Public đọc `page_publish_cache.compiled_json`.
- Password page dùng `password_hash` (argon2id/bcrypt); xác thực xong mới trả compiled_json.
- Analytics deferred (để sau).

---

## 15) Checklist tính năng (đủ - không sót)

- [x] User đăng ký / tạo nhiều bio pages
- [x] Mỗi page có 1 domain + path (root + nested)
- [x] Đổi path có redirect path cũ
- [x] 5 theme preset
- [x] Appearance tuỳ chỉnh tạo custom theme (không ghi đè preset), custom theme reuse nhiều page
- [x] Theme quyết định page style; group override thắng theme
- [x] Block system (text/image/product/link_group/…)
- [x] Link theo group; style theo group, không style từng link
- [x] Group layout: list/cards/grid + config (columns, gap,…)
- [x] Drag & drop blocks và drag & drop links
- [x] Draft / Publish (public chỉ đổi khi Publish)
- [x] Mockup realtime trong editor/appearance
- [x] Wallpaper: preset + upload + blur/dim/overlay
- [x] Locale per page (vi/en)
- [x] Private/password page
- [x] Analytics để sau
- [x] Product block tối ưu (manual card, mở đường catalog sau)
- [x] Không giới hạn nghiệp vụ, khuyến nghị soft limit kỹ thuật



---

## 16) Plans (Free/Pro) & Feature gating

- Có 2 gói: **Free / Pro**
- Quyền theo plan:
  - **Custom domain**: chỉ **Pro**
  - **Password page**: chỉ **Pro**
- Các tính năng khác trong MVP: dùng chung cho cả Free/Pro (trừ khi về sau quy định thêm).


---

## 4) Domain / Path / Redirect (cập nhật)

- Hỗ trợ domain hệ thống (ví dụ `domain.com`) để host nhiều page theo path.
- Custom domain (ví dụ `yen.dev`) **chỉ trỏ tới 1 page duy nhất** và dùng tại **root**: `yen.dev/`.
- Khi user đổi path trên domain hệ thống hoặc đổi root route của custom domain, hệ thống phải hỗ trợ **redirect** route cũ → route mới.


---

## 6) Theme Marketplace & Modes

- Có **Theme Marketplace**.
- Theme preset hệ thống gồm:
  - **Free themes**
  - **Pro themes** (chỉ user Pro sử dụng)
- Theme preset là **immutable**.
- Nếu user chỉnh appearance khác preset → hệ thống tạo **custom theme** (patch) và gán cho page (custom theme có thể reuse nhiều page).
- Modes: có ngay (`light/dark/compact`).


---

## 7) Private / Password page (cập nhật)

- Page có chế độ `public` hoặc `password`.
- **Password page chỉ dành cho Pro**.
- Password là **1 mật khẩu cho cả page**.
- Người xem nhập đúng password sẽ được “nhớ” trong **7 ngày** (session/token).


---

## 17) GDPR/PDPA tối thiểu

- Có tối thiểu các chức năng:
  - **Delete account**
  - **Delete pages**
  - **Delete assets** (theo policy: asset upload của user bị xoá khi user xoá tài khoản)
- Không yêu cầu analytics/compliance nâng cao trong giai đoạn đầu.


---

## 18) Performance & Caching (cập nhật)

- Mục tiêu: public load **siêu nhanh**.
- Caching: **DB + CDN cache** (không cần Redis giai đoạn đầu)
  - Publish → compile JSON và lưu vào DB (publish cache)
  - Public request → CDN cache response theo `hostname` (và path nếu domain hệ thống)
  - Infra xử lý http→https và www/non-www (không xử lý ở app).
