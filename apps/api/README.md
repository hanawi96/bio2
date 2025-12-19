# LinkBio API (Go Fiber)

## Structure

```
apps/api/
├─ cmd/api/main.go          # Entry point
├─ internal/
│  ├─ config/               # Configuration
│  ├─ db/                   # Database pool
│  ├─ middleware/           # Auth, recover, etc.
│  ├─ router/               # Route setup
│  ├─ modules/
│  │  └─ public/            # GET /r, POST /r/password
│  └─ shared/               # Common types
└─ migrations/              # SQL migrations
```

## Run

```bash
# Set environment
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/linkbio?sslmode=disable"

# Run
go run cmd/api/main.go
```

## Endpoints (Slice 1)

- `GET /health` - Health check
- `GET /r` - Public page render (by Host header)
- `POST /r/password` - Password verification (placeholder)
