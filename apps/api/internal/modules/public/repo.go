package public

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PublishCacheResult contains compiled_json and stable etag
type PublishCacheResult struct {
	Compiled map[string]interface{}
	ETag     string
}

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// ResolveHostToPage finds page_id from domain hostname
// Custom domain: 1 domain = 1 page, always at path "/"
// System domain: can have multiple pages at different paths
// Returns page_id, access_type, error
func (r *Repo) ResolveHostToPage(ctx context.Context, hostname string) (int64, string, error) {
	// Custom domain logic: always resolve to "/" (1 domain = 1 page)
	// This query handles both custom domain (path="/") and system domain
	query := `
		SELECT bp.id, bp.access_type
		FROM domains d
		JOIN page_routes pr ON pr.domain_id = d.id
		JOIN bio_pages bp ON bp.id = pr.page_id
		WHERE d.hostname = $1
		  AND pr.path = '/'
		  AND pr.is_current = true
		  AND d.status = 'active'
		  AND bp.status = 'published'
		LIMIT 1
	`

	var pageID int64
	var accessType string
	err := r.pool.QueryRow(ctx, query, hostname).Scan(&pageID, &accessType)
	if err != nil {
		return 0, "", err
	}

	return pageID, accessType, nil
}

// GetPublishCache returns compiled_json and stable etag for a page
func (r *Repo) GetPublishCache(ctx context.Context, pageID int64) (*PublishCacheResult, error) {
	query := `
		SELECT compiled_json, etag
		FROM page_publish_cache
		WHERE page_id = $1
	`

	var compiledJSON []byte
	var etag string
	err := r.pool.QueryRow(ctx, query, pageID).Scan(&compiledJSON, &etag)
	if err != nil {
		return nil, err
	}

	var compiled map[string]interface{}
	if err := json.Unmarshal(compiledJSON, &compiled); err != nil {
		return nil, err
	}

	return &PublishCacheResult{
		Compiled: compiled,
		ETag:     etag,
	}, nil
}
