package public

import (
	"context"
	"errors"
)

var (
	ErrNotFound         = errors.New("page not found")
	ErrPasswordRequired = errors.New("password required")
)

// PageResult contains compiled page data and metadata
type PageResult struct {
	Compiled            map[string]interface{}
	ETag                string
	IsPasswordProtected bool
}

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// GetCompiledPage resolves host -> page -> compiled_json
// Custom domain: 1 domain = 1 page, always at path "/"
func (s *Service) GetCompiledPage(ctx context.Context, host string) (*PageResult, error) {
	// 1. Resolve domain to page_id (custom domain always uses "/")
	pageID, accessType, err := s.repo.ResolveHostToPage(ctx, host)
	if err != nil {
		return nil, ErrNotFound
	}

	// 2. Check access type - password protected pages require auth
	if accessType == "password" {
		// TODO: Check cookie session for password access (slice 2)
		// For now, always return PASSWORD_REQUIRED
		return nil, ErrPasswordRequired
	}

	// 3. Get compiled JSON and stable ETag from cache
	cacheResult, err := s.repo.GetPublishCache(ctx, pageID)
	if err != nil {
		return nil, ErrNotFound
	}

	return &PageResult{
		Compiled:            cacheResult.Compiled,
		ETag:                cacheResult.ETag,
		IsPasswordProtected: accessType == "password",
	}, nil
}
