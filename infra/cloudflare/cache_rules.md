# Cloudflare Cache Rules

## Public Pages (/r)

Cache-Control header từ origin:
```
Cache-Control: public, max-age=60, s-maxage=3600, stale-while-revalidate=86400
```

- `max-age=60`: Browser cache 1 phút
- `s-maxage=3600`: CDN cache 1 giờ
- `stale-while-revalidate=86400`: Serve stale trong 24h khi revalidate

## Purge on Publish

Khi user publish page:
1. API gọi Cloudflare Purge API
2. Purge theo hostname: `https://{hostname}/*`

## Bypass Rules

Bypass cache khi:
- Cookie `pats` present (password session)
- Query param `?preview=1`
