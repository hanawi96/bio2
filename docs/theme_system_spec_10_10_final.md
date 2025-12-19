# Theme System 10/10 cho Link-in-Bio Builder (Spec + Contract)

> Mục tiêu: Theme phải **đẹp – linh hoạt – an toàn – render nhanh**.  
> Preset **immutable** (không chỉnh trực tiếp). Bất kỳ thay đổi appearance nào khác preset → tạo **Custom Theme** bằng **patch**.

---

## 0) Thuật ngữ
- **Preset Theme**: theme hệ thống (5 theme), chỉ đọc, không bị ghi đè.
- **Custom Theme**: theme do user tạo từ preset + patch (có thể reuse cho nhiều page).
- **Patch**: chỉ lưu các key/value user đã thay đổi so với preset.
- **Compiled Theme**: kết quả merge (preset + patch), dùng để render nhanh.
- **Tokens**: primitive values (màu/size/radius…), không gắn UI.
- **Semantic**: mapping có ý nghĩa UI (primary, text.muted…).
- **Recipes**: preset component (link card, button…), dùng semantic.
- **Page Defaults**: defaults riêng cho Link-in-Bio page (layout, header, link-group defaults…).
- **Override**: cấu hình style riêng của **group** (thắng theme).

---

## 1) Nguyên tắc bắt buộc
1. **Preset immutable**: UI không bao giờ sửa trực tiếp preset.
2. **Custom theme = preset + patch**: patch chứa mọi thay đổi.
3. **UI chỉ dùng semantic/recipes/page**, không đọc trực tiếp tokens gốc.
4. **Style priority** (từ thấp → cao): `Theme (compiled)` → `Group override`.
5. Theme phải **serializable** (JSON), **cache-friendly**, và có **schemaVersion** để migrate.

---


## Marketplace & Tier (Free/Pro)
- Preset themes có thể phân loại:
  - `tier`: `free` | `pro`
  - (tuỳ chọn) `price`: nếu sau này bán theme
- User Free chỉ dùng `free`; user Pro dùng `free` + `pro`.
- Theme preset vẫn **immutable**; user chỉnh khác preset → tạo custom theme bằng patch.

## 2) Theme Object: cấu trúc chuẩn (Top-level)
Theme luôn là 1 JSON object với các section dưới đây:

```json
{
  "meta": {},
  "tokens": {},
  "semantic": {},
  "recipes": {},
  "page": {},
  "background": {},
  "modes": {}
}
```

**Yêu cầu**:
- Không đổi tên section tuỳ tiện.
- Keys ổn định để AI / codegen / validator hoạt động.

---

## 3) `meta` — Contract + quản trị theme (không ảnh hưởng render)
### 3.1 Fields bắt buộc
- `id` (string): định danh theme
- `name` (string)
- `version` (string): semver
- `schemaVersion` (int): version của cấu trúc theme
- `author` (string | object)
- `description` (string, ngắn)
- `supports` (object): capabilities (darkMode, animation, glass…)
- `tier` (string): `free|pro` (preset theme phân loại)
- `visibility` (string): `public|unlisted|private` (tuỳ marketplace)
- `preview` (object): thumbnailUrl, demoUrl (optional)

### 3.2 `meta.contract` (BẮT BUỘC cho builder)
Mục tiêu: định nghĩa **appearance controls** và **guardrails** để editor tạo UI chỉnh theme đúng cách.

- `controls[]`: danh sách control hiển thị trong Appearance.
  - `keyPath`: đường dẫn trong theme (vd `page.layout.textAlign`)
  - `type`: `select|slider|color|toggle|number`
  - `label`: tên hiển thị
  - `options`: (select) list option
  - `min/max/step`: (slider/number)
  - `affects`: gợi ý nhóm thành phần bị ảnh hưởng (page, header, linkGroup)
- `constraints`: validate nâng cao (range, enum)
- `defaultsForMissing`: fallback khi key thiếu

**Ví dụ**:
```json
{
  "meta": {
    "schemaVersion": 1,
    "contract": {
      "controls": [
        { "keyPath": "page.layout.textAlign", "type": "select",
          "options": ["left","center","right"], "label": "Text align (page)" },

        { "keyPath": "page.layout.baseFontSize", "type": "select",
          "options": ["S","M","L","XL"], "label": "Font size (page)" },

        { "keyPath": "page.layout.pagePadding", "type": "slider",
          "min": 0, "max": 32, "step": 1, "label": "Page padding" },

        { "keyPath": "page.layout.blockGap", "type": "slider",
          "min": 0, "max": 32, "step": 1, "label": "Block spacing" },

        { "keyPath": "page.mode", "type": "select",
          "options": ["light","dark","compact"], "label": "Mode (page)" },

        { "keyPath": "page.defaults.linkGroup.textAlign", "type": "select",
          "options": ["left","center","right"], "label": "Link group text align (default)" },

        { "keyPath": "page.defaults.linkGroup.fontSize", "type": "select",
          "options": ["S","M","L","XL"], "label": "Link group font size (default)" },

        { "keyPath": "background.effects.blur", "type": "slider",
          "min": 0, "max": 12, "step": 1, "label": "Background blur" },

        { "keyPath": "background.effects.dim", "type": "slider",
          "min": 0, "max": 0.8, "step": 0.05, "label": "Background dim" }
      ],
      "constraints": {
        "background.effects.dim": { "min": 0, "max": 0.8 },
        "background.effects.blur": { "min": 0, "max": 12 }
      },
      "notes": [
        "Theme controls chỉ chỉnh ở mức theme/page defaults.",
        "Muốn link style riêng: tạo Link Group riêng và set group style_override."
      ]
    }
  }
}
```

---

## 4) `tokens` — Design primitives (atomic, không biết UI)
Chỉ chứa giá trị thô. Không dùng tên kiểu “buttonPrimary”.

### 4.1 Nhóm tokens chuẩn
- `color`: palette (gray.50..900, blue.100..700…)
- `typography`: fontFamily, fontSizeScale, fontWeight, lineHeight
- `space`: spacing scale (2,4,8,12,16…)
- `size`: icon/avatar/button heights (generic)
- `radius`: radius scale
- `elevation`: shadow levels
- `motion`: duration/easing
- `zIndex`: layers
- `breakpoint`: responsive breakpoints

**Rule**: UI không đọc tokens trực tiếp, chỉ thông qua semantic/recipes/page.

---

## 5) `semantic` — Mapping có ý nghĩa UI
UI dùng semantic để ổn định khi đổi theme.

### 5.1 `semantic.color`
- `primary`, `secondary`, `accent`
- `text.default`, `text.muted`, `text.invert`
- `border.default`, `divider`
- `surface.page`, `surface.card`, `surface.overlay`
- `danger`, `success`, `warning`

### 5.2 `semantic.typography`
- `heading`, `body`, `caption`, `button`
Mỗi item nên map về tokens typography (size/weight/lineHeight).

### 5.3 `semantic.surface`
Kết hợp color + elevation + blur/glass (nếu có).

---

## 6) `recipes` — Preset component (Agnostic)
### 6.1 Bắt buộc phải có (Link-in-bio)
- `linkItem`: style 1 item link (card/button)
- `linkGroup`: container cho nhóm link
- `header`: avatar/name/bio area
- `button` (optional)
- `badge` (optional)

### 6.2 Variants (BẮT BUỘC cho linkGroup layout)
`linkGroup` phải hỗ trợ layout:
- `list`
- `cards`
- `grid`

Ví dụ:
```json
{
  "recipes": {
    "linkGroup": {
      "base": { "gap": "space.3" },
      "variants": {
        "layout": {
          "list":  { "columns": 1 },
          "cards": { "columns": 1, "cardStyle": "elevated" },
          "grid":  { "columns": 2 }
        }
      }
    }
  }
}
```

---

## 7) `page` — Defaults riêng cho Link-in-Bio Builder (BẮT BUỘC)
Đây là phần giúp builder merge style đúng và ổn định.

### 7.1 `page.layout`
- `maxWidth`
- `pagePadding`
- `blockGap`
- `textAlign` (page-level default align)
- `baseFontSize` (S/M/L/XL)
- `contentAlign` (left/center/right nếu cần)
- (optional) `safeArea` (mobile)

### 7.2 `page.defaults` (BẮT BUỘC)
Defaults cho các nhóm nội dung để inherit:
- `linkGroup` (textAlign, fontSize, radius, padding, shadow…)
- `textBlock` (heading/body/caption styles)
- `imageBlock` (radius, shadow)
- `productBlock` (card styles)

**Ví dụ**:
```json
{
  "page": {
    "layout": { "maxWidth": 520, "pagePadding": 16, "blockGap": 12,
                "textAlign": "center", "baseFontSize": "M" },
    "defaults": {
      "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "radius.4" }
    }
  }
}
```

---

## 8) `background` — Wallpaper + effects (preset/upload)
### 8.1 Cấu trúc chuẩn
```json
{
  "background": {
    "wallpaper": { "kind": "preset|upload", "assetId": 123 },
    "effects": {
      "blur": 0,
      "dim": 0.0,
      "overlayColor": "rgba(0,0,0,0)"
    }
  }
}
```

### 8.2 Guardrails khuyến nghị (để đẹp + nhanh)
- `blur`: 0..12
- `dim`: 0..0.8
- overlayColor default nên tồn tại (để đảm bảo readability)

---

## 9) `modes` — Light/Dark/Compact (optional nhưng nên có)
`modes` là object map name → override.
- override có thể đè `tokens` hoặc `semantic` hoặc `page`.

Ví dụ:
```json
{
  "modes": {
    "dark": { "semantic": { "color": { "surface": { "page": "color.gray.900" }}}}
  }
}
```

---

## 10) Custom Theme Workflow (Preset → Patch → Compiled)
### 10.1 Quy tắc tạo custom
- User chọn preset A.
- Khi user đổi bất kỳ control nào trong Appearance:
  - Nếu page đang dùng preset trực tiếp → tạo `themes_custom` mới:
    - `basedOnPresetId = presetA`
    - `patch = {<keyPath>: <value>, ...}`
  - Gán `theme_custom_id` cho page.
- Các lần chỉnh tiếp theo: cập nhật patch của custom theme.

### 10.2 Patch format (chuẩn cho AI coding)
- Patch là object chứa các keyPath đầy đủ (khuyến nghị dùng dạng nested JSON giống theme).
- Chỉ chứa phần thay đổi.

Ví dụ patch:
```json
{
  "page": { "layout": { "textAlign": "left", "baseFontSize": "L" } },
  "background": { "effects": { "dim": 0.35 } }
}
```

### 10.3 Compiled theme
- `compiled = deepMerge(preset, patch, modeOverride?)`
- Lưu `compiled_config` để render nhanh.
- Khi publish page: dùng compiled theme để tạo `compiled_json` public.

---

## 11) Group override (đúng logic builder)
### 11.1 Group override structure
Group chỉ lưu override keys cần thiết (links trong group sẽ **thừa kế** các style này):

```json
{ "textAlign": "center", "fontSize": "XL" }
```

### 11.2 Merge rule khi render
- `groupFinalStyle = deepMerge(theme.page.defaults.linkGroup, group.override)`
- Nếu group có key → key đó thắng, không bị ảnh hưởng bởi appearance/theme.

---


## 11.3 Quy tắc “muốn link style riêng → tạo group riêng” (BẮT BUỘC)
- Hệ thống **không hỗ trợ style theo từng link** và **không có block-level override**.
- Nếu user muốn một link (hoặc một cụm link) có style khác các link còn lại:
  1) **Tạo một Link Group mới** (có thể chỉ chứa 1 link)
  2) Set `style_override` cho group đó
  3) Các link trong group sẽ thừa kế style của group

**Ví dụ** (CTA link riêng):
- Group `CTA` có 1 link “Contact me”
- `style_override`: `{ "textAlign": "left", "fontSize": "XL" }`
→ chỉ nhóm CTA khác, các group khác vẫn theo theme.



## 12) Performance & Cache (bắt buộc để “siêu mượt”)
### 12.1 Public cache
- Khi Publish: compile page thành `page_publish_cache.compiled_json`.
- Public request: chỉ resolve route + trả compiled_json (rất nhanh).

### 12.2 Payload size
- Theme object nên giới hạn “hợp lý”:
  - Preset đầy đủ (ok)
  - Patch nhỏ
  - Public compiled_json chỉ chứa cái cần render

---

## 13) Minimal “Theme Preset” mẫu (skeleton)
> AI có thể dùng skeleton này để sinh theme preset hợp lệ.

```json
{
  "meta": {
    "id": "preset.themeA",
    "name": "Theme A",
    "version": "1.0.0",
    "schemaVersion": 1,
    "author": "system",
    "description": "Clean center layout",
    "supports": { "modes": ["light","dark"], "animation": true },
    "contract": {
      "controls": [
        { "keyPath": "page.layout.textAlign", "type": "select",
          "options": ["left","center","right"], "label": "Text align (page)" },
        { "keyPath": "page.layout.baseFontSize", "type": "select",
          "options": ["S","M","L","XL"], "label": "Font size (page)" },
        { "keyPath": "page.layout.pagePadding", "type": "slider",
          "min": 0, "max": 32, "step": 1, "label": "Page padding" },
        { "keyPath": "page.layout.blockGap", "type": "slider",
          "min": 0, "max": 32, "step": 1, "label": "Block spacing" },

        { "keyPath": "page.mode", "type": "select",
          "options": ["light","dark","compact"], "label": "Mode (page)" },

        { "keyPath": "page.defaults.linkGroup.textAlign", "type": "select",
          "options": ["left","center","right"], "label": "Link group text align (default)" },
        { "keyPath": "page.defaults.linkGroup.fontSize", "type": "select",
          "options": ["S","M","L","XL"], "label": "Link group font size (default)" },

        { "keyPath": "background.effects.blur", "type": "slider",
          "min": 0, "max": 12, "step": 1, "label": "Background blur" },
        { "keyPath": "background.effects.dim", "type": "slider",
          "min": 0, "max": 0.8, "step": 0.05, "label": "Background dim" }
      ]
    }
  },
  "tokens": {
    "color": { "gray": { "50": "#f9fafb", "900": "#111827" } },
    "typography": { "fontFamily": { "sans": "Inter" } },
    "space": { "1": 4, "2": 8, "3": 12, "4": 16 },
    "radius": { "4": 16 },
    "elevation": { "sm": "0 1px 2px rgba(0,0,0,0.08)" }
  },
  "semantic": {
    "color": {
      "primary": "tokens.color.gray.900",
      "text": { "default": "tokens.color.gray.900", "muted": "tokens.color.gray.50" },
      "surface": { "page": "tokens.color.gray.50", "card": "tokens.color.gray.50" }
    },
    "typography": { "body": { "fontFamily": "tokens.typography.fontFamily.sans" } }
  },
  "recipes": {
    "linkItem": { "base": { "radius": "tokens.radius.4", "shadow": "tokens.elevation.sm" } },
    "linkGroup": {
      "base": { "gap": "tokens.space.3" },
      "variants": { "layout": { "list": { "columns": 1 }, "cards": { "columns": 1 }, "grid": { "columns": 2 } } }
    }
  },
  "page": {
    "layout": { "maxWidth": 520, "pagePadding": 16, "blockGap": 12, "textAlign": "center", "baseFontSize": "M" },
    "defaults": { "linkGroup": { "textAlign": "center", "fontSize": "M", "radius": "tokens.radius.4" } }
  },
  "background": {
    "wallpaper": { "kind": "preset", "assetId": 0 },
    "effects": { "blur": 0, "dim": 0.0, "overlayColor": "rgba(0,0,0,0)" }
  },
  "modes": {}
}
```

---

## 14) Checklist “Theme 10/10”
- [x] Preset immutable + custom patch workflow
- [x] `meta.contract` để generate Appearance UI + validate
- [x] tokens/semantic/recipes tách bạch
- [x] `page.layout` + `page.defaults` cho builder merge đúng (không block override, không per-link style)
- [x] LinkGroup recipes có variants list/cards/grid
- [x] Background wallpaper (preset/upload) + blur/dim/overlayColor
- [x] Modes override (optional)
- [x] Cache-friendly, serialize JSON dễ dàng, schemaVersion rõ

