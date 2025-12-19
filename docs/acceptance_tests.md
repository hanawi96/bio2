# Acceptance Tests (Manual) — Link-in-Bio (Final)

> Mục tiêu: bạn không biết code vẫn có thể kiểm tra AI làm đúng hay không.
> Mỗi mục có **Steps** và **Expected**.

---

## 0) Setup
- Chạy migration + seed (FREE/PRO + theme presets).
- Start API + Web.
- Có ít nhất 1 theme free và 1 theme pro trong marketplace.

---

## 1) Auth & Session
### 1.1 Register/Login
**Steps**
1) Register bằng email/password
2) Login
3) Reload trang dashboard

**Expected**
- Đăng nhập thành công, không bị logout khi reload.
- Dashboard load được danh sách pages.

---

## 2) Create Page + Draft
### 2.1 Create Page
**Steps**
1) Tạo page mới (locale vi, theme preset free)
2) Vào editor

**Expected**
- Page có status `draft`
- Preview hiển thị đúng header + empty blocks

### 2.2 Save button (không auto-save)
**Steps**
1) Thêm text block, thêm link group block
2) Không bấm Save → reload trang
3) Làm lại, lần này bấm Save → reload

**Expected**
- Nếu không Save: thay đổi không được lưu
- Nếu Save: reload vẫn còn

---

## 3) Drag & Drop
### 3.1 Reorder blocks
**Steps**
1) Tạo 3 blocks (text/link_group/image)
2) Drag đổi thứ tự
3) Save → reload

**Expected**
- Thứ tự đúng như drag
- sort_key được update hợp lệ (không trùng)

### 3.2 Reorder links inside group
**Steps**
1) Trong link group tạo 3 links
2) Drag đổi thứ tự links
3) Save → reload

**Expected**
- Thứ tự links đúng như drag

---

## 4) Theme & Appearance (preset immutable + custom patch)
### 4.1 Apply preset + edit appearance
**Steps**
1) Apply theme preset A
2) Vào Appearance chỉnh 1 giá trị (vd textAlign từ center → left)
3) Save

**Expected**
- Hệ thống tạo custom theme (patch) hoặc update custom theme
- Theme preset không bị thay đổi
- Page dùng custom theme id (không ghi đè preset)

### 4.2 Group override wins theme
**Steps**
1) Theme page textAlign = center
2) Set group A style_override textAlign = left
3) Save → Publish
4) Public view

**Expected**
- Group A căn trái
- Các group khác (không override) vẫn theo theme (center)

---

## 5) Draft vs Published
### 5.1 Public chỉ đổi khi Publish
**Steps**
1) Publish lần 1
2) Mở public page → ghi lại nội dung
3) Sửa draft (thêm block/text) → Save (không publish)
4) Refresh public page

**Expected**
- Public page KHÔNG đổi khi chỉ Save
- Khi Publish lần 2 → public page đổi

---

## 6) Public performance (Publish cache + CDN)
### 6.1 Public endpoint returns compiled_json
**Steps**
1) Publish
2) Mở public
3) Inspect network `/r`

**Expected**
- Response nhanh
- Có `Cache-Control` + `ETag`
- Payload là compiled_json (blocks sorted, group final_style)

---

## 7) Password Page (Pro-only) + remember 7 days
### 7.1 Free cannot enable password
**Steps**
1) User plan FREE
2) Vào settings -> bật password page

**Expected**
- Bị chặn: trả lỗi `FORBIDDEN_PRO_REQUIRED` hoặc UI thông báo upgrade

### 7.2 Pro enable password
**Steps**
1) Upgrade user → Pro (seed hoặc giả lập)
2) Bật password page, set password
3) Publish
4) Mở public page

**Expected**
- Public trả `PASSWORD_REQUIRED`
- Nhập đúng password → vào page
- Refresh/đóng mở lại trong 7 ngày vẫn vào được (cookie/session)

---

## 8) Custom Domain (Pro-only, 1 domain = 1 page)
### 8.1 Free cannot add custom domain
**Steps**
1) User plan FREE
2) Add domain `yen.dev`

**Expected**
- Bị chặn (Pro-only)

### 8.2 Pro attach custom domain
**Steps**
1) User plan PRO
2) Add domain + attach page
3) Request với Host=yen.dev

**Expected**
- Domain trỏ đúng page tại `/`
- Không có nested path cho custom domain (luôn root)

---

## 9) Theme marketplace free/pro
### 9.1 Free cannot apply pro theme
**Steps**
1) User FREE
2) Mở marketplace, chọn theme tier=pro và apply

**Expected**
- Bị chặn, UI yêu cầu upgrade

### 9.2 Pro can apply pro theme
**Steps**
1) User PRO
2) Apply pro theme
3) Publish

**Expected**
- Apply thành công, public hiển thị đúng

---

## 10) GDPR minimum delete
### 10.1 Delete account deletes pages & assets
**Steps**
1) Tạo page, upload 1 image asset, publish
2) Call DELETE /api/account
3) Login lại hoặc inspect DB

**Expected**
- User bị xoá
- Pages, blocks, groups, links, publish cache, sessions bị xoá
- Assets thuộc user bị xoá theo policy (DB cascade)
