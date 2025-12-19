# Link-in-Bio Builder

Monorepo: Svelte 5 (SvelteKit) + Go Fiber + PostgreSQL

## Structure

```
linkbio/
├─ apps/
│  ├─ web/          # SvelteKit (Editor + Public renderer)
│  └─ api/          # Go Fiber API
├─ packages/
│  └─ shared/       # Types/contracts dùng chung
├─ docs/            # Specs & documentation
├─ infra/           # CDN/cache notes
└─ scripts/         # Dev/deploy utilities
```

## Quick Start

### Prerequisites
- Node.js 20+
- Go 1.21+
- PostgreSQL 15+
- psql CLI

### 1. Setup Database
```bash
# Create database
createdb linkbio

# Run migrations (Windows: dùng psql trực tiếp)
psql -d linkbio -f apps/api/migrations/0001_init.sql

# Hoặc Linux/Mac:
./scripts/migrate.sh
```

### 2. Seed Data
```bash
# Windows:
psql -d linkbio -f scripts/seed.sql

# Linux/Mac:
./scripts/seed.sh
```

### 3. Start Development

**Terminal 1 - API (port 8080):**
```bash
cd apps/api
set DATABASE_URL=postgres://postgres:postgres@localhost:5432/linkbio?sslmode=disable
go run cmd/api/main.go
```

**Terminal 2 - Web (port 5173):**
```bash
cd apps/web
npm install
npm run dev
```

### 4. Verify
- API Health: http://localhost:8080/health
- Web: http://localhost:5173

---

## Vertical Slice 1: Public Render + Publish Cache

### Endpoints
- `GET /r` - Returns compiled_json from page_publish_cache by Host
- `POST /r/password` - Password verification (placeholder - returns NOT_IMPLEMENTED)

### Renderer
- Public route: `(public)/[...path]` - rest parameter matches `/` (empty) and all paths
- Supported blocks: `text`, `link_group`
- Note: SvelteKit rest parameter `[...path]` is inherently optional, matching both `/` and `/any/path`

### Domain Logic
- Custom domain: 1 domain = 1 page, always resolves to path `/`
- System domain: supports multiple pages (future)

---

## Verification Checklist (theo acceptance_tests.md)

### Setup
```bash
# 1. Migrate
psql -d linkbio -f apps/api/migrations/0001_init.sql

# 2. Seed
psql -d linkbio -f scripts/seed.sql

# 3. Start API
cd apps/api && go run cmd/api/main.go

# 4. Start Web (terminal khác)
cd apps/web && npm run dev
```

### Test Case 1: 200 OK - Returns compiled_json
```bash
curl -i "http://localhost:8080/r?host=localhost:5173"
```
**Expected:**
- Status: `200 OK`
- Header `ETag`: có giá trị (vd: `"a1b2c3..."`)
- Header `Cache-Control`: `public, max-age=60, s-maxage=3600, stale-while-revalidate=86400`
- Body: JSON với `version`, `page`, `theme`, `blocks`

### Test Case 2: 404 NOT_FOUND - Page không tồn tại
```bash
curl -i "http://localhost:8080/r?host=unknown.domain.com"
```
**Expected:**
- Status: `404 Not Found`
- Body:
```json
{"error":{"code":"NOT_FOUND","message":"Page not found"}}
```

### Test Case 3: 401 PASSWORD_REQUIRED - Page có password
```bash
# Cần update seed để có page với access_type='password'
# Hoặc update trực tiếp trong DB:
# UPDATE bio_pages SET access_type='password' WHERE id=1;
curl -i "http://localhost:8080/r?host=localhost:5173"
```
**Expected:**
- Status: `401 Unauthorized`
- Body:
```json
{"error":{"code":"PASSWORD_REQUIRED","message":"This page is password protected"}}
```

### Test Case 4: POST /r/password - Placeholder
```bash
curl -i -X POST "http://localhost:8080/r/password" \
  -H "Content-Type: application/json" \
  -d '{"password":"test"}'
```
**Expected:**
- Status: `501 Not Implemented`
- Body:
```json
{"error":{"code":"NOT_IMPLEMENTED","message":"Password verification not implemented yet"}}
```

### Test Case 5: Web renders page
1. Mở http://localhost:5173
2. **Expected:** Hiển thị page với header (name, bio) và blocks (link_group, text)

---

## Vertical Slice 2: Draft Save + Publish Compile

### Endpoints
- `POST /api/pages` - Create new page
- `GET /api/pages` - List pages
- `GET /api/pages/:id/draft` - Get draft state
- `POST /api/pages/:id/save` - Save draft state
- `POST /api/pages/:id/publish` - Compile and publish

### Editor UI
- `/app/pages` - List pages + create
- `/app/pages/[pageId]/edit` - Editor with preview

### Verification (acceptance_tests.md mục 5.1)

#### Test: Save không làm public đổi
```bash
# 1. Tạo page mới
curl -X POST "http://localhost:8080/api/pages" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Page","locale":"vi","theme_preset_key":"theme_a"}'
# Response: {"page_id": 2}

# 2. Publish lần 1
curl -X POST "http://localhost:8080/api/pages/2/publish" \
  -H "Content-Type: application/json" -d '{}'

# 3. Ghi lại public content
curl "http://localhost:8080/r?host=localhost:5173" > public_v1.json

# 4. Save draft (thêm block mới) - không publish
curl -X POST "http://localhost:8080/api/pages/2/save" \
  -H "Content-Type: application/json" \
  -d '{"page":{"id":2,"title":"Test Page","locale":"vi","status":"published","theme":{"preset_key":"theme_a","custom_id":null,"mode":"light"},"access_type":"public","settings":{"header":{"name":"Test","bio":"Hello","verified":false,"social":[]}}},"blocks":[{"id":1,"type":"text","sort_key":"a0","ref_id":null,"content":{"text":"NEW TEXT","variant":"heading"},"is_visible":true}],"link_groups":{},"links":{}}'

# 5. Check public - KHÔNG đổi
curl "http://localhost:8080/r?host=localhost:5173" > public_v2.json
# public_v1.json và public_v2.json phải giống nhau

# 6. Publish lần 2
curl -X POST "http://localhost:8080/api/pages/2/publish" \
  -H "Content-Type: application/json" -d '{}'

# 7. Check public - ĐÃ đổi
curl "http://localhost:8080/r?host=localhost:5173"
# Phải có "NEW TEXT" trong blocks
```

#### Test via UI
1. Mở http://localhost:5173/app/pages
2. Tạo page mới
3. Thêm Text block, thêm Link Group block
4. Bấm Save → reload trang → thay đổi vẫn còn
5. Bấm Publish → mở public page → thấy nội dung mới

---

## Error Response Format (theo api_contract.md)
```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "Human readable message",
    "details": {}
  }
}
```
