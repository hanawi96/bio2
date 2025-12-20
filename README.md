# LinkBio

Link in Bio builder - Create beautiful bio pages.

## Tech Stack

- **Frontend**: SvelteKit 2 + Svelte 5
- **Backend**: Go Fiber
- **Database**: PostgreSQL

## Quick Start

### 1. Start Database

```bash
docker-compose up -d
```

### 2. Run Migration

```bash
# Windows PowerShell
Get-Content docs/link_in_bio_init_migration_updated.sql | docker exec -i linkbio_postgres psql -U linkbio -d linkbio
```

### 3. Start Backend

```bash
cd api
go mod tidy
go run cmd/main.go
```

### 4. Start Frontend

```bash
cd web
npm install
npm run dev
```

## Project Structure

```
linkbio/
├── api/                    # Go Fiber backend
│   ├── cmd/main.go
│   └── internal/
│       ├── config/
│       ├── database/
│       ├── handler/
│       ├── middleware/
│       ├── model/
│       ├── repo/
│       ├── service/
│       └── util/
├── web/                    # SvelteKit frontend
│   └── src/
│       ├── lib/
│       │   ├── api/
│       │   ├── stores/
│       │   └── utils/
│       └── routes/
│           ├── (app)/      # Protected routes
│           └── (public)/   # Public pages
├── docs/                   # Specifications
└── docker-compose.yml
```

## API Endpoints

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/me`

### Pages
- `GET /api/pages`
- `POST /api/pages`
- `GET /api/pages/:id/draft`
- `POST /api/pages/:id/save`
- `POST /api/pages/:id/publish`
- `DELETE /api/pages/:id`

### Themes
- `GET /api/themes/presets`
- `POST /api/themes/custom`

### Public
- `GET /r` - Render public page
- `POST /r/password` - Verify password
