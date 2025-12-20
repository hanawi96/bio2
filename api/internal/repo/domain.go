package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type DomainRepo struct {
	db *pgxpool.Pool
}

func NewDomainRepo(db *pgxpool.Pool) *DomainRepo {
	return &DomainRepo{db: db}
}

func (r *DomainRepo) GetByHostname(ctx context.Context, hostname string) (*model.Domain, error) {
	var d model.Domain
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, hostname, status, is_system, created_at, updated_at
		FROM domains WHERE hostname = $1
	`, hostname).Scan(&d.ID, &d.UserID, &d.Hostname, &d.Status, &d.IsSystem, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DomainRepo) GetSystemDomain(ctx context.Context) (*model.Domain, error) {
	var d model.Domain
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, hostname, status, is_system, created_at, updated_at
		FROM domains WHERE is_system = true LIMIT 1
	`).Scan(&d.ID, &d.UserID, &d.Hostname, &d.Status, &d.IsSystem, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DomainRepo) Create(ctx context.Context, userID *int64, hostname string, isSystem bool) (*model.Domain, error) {
	var d model.Domain
	err := r.db.QueryRow(ctx, `
		INSERT INTO domains (user_id, hostname, is_system, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, user_id, hostname, status, is_system, created_at, updated_at
	`, userID, hostname, isSystem).Scan(&d.ID, &d.UserID, &d.Hostname, &d.Status, &d.IsSystem, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DomainRepo) GetRouteByDomainAndPath(ctx context.Context, domainID int64, path string) (*model.PageRoute, error) {
	var route model.PageRoute
	err := r.db.QueryRow(ctx, `
		SELECT id, page_id, domain_id, path, is_current, redirect_to_route_id, created_at
		FROM page_routes WHERE domain_id = $1 AND path = $2 AND is_current = true
	`, domainID, path).Scan(
		&route.ID, &route.PageID, &route.DomainID, &route.Path,
		&route.IsCurrent, &route.RedirectToRouteID, &route.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *DomainRepo) CreateRoute(ctx context.Context, pageID, domainID int64, path string) (*model.PageRoute, error) {
	var route model.PageRoute
	err := r.db.QueryRow(ctx, `
		INSERT INTO page_routes (page_id, domain_id, path)
		VALUES ($1, $2, $3)
		RETURNING id, page_id, domain_id, path, is_current, redirect_to_route_id, created_at
	`, pageID, domainID, path).Scan(
		&route.ID, &route.PageID, &route.DomainID, &route.Path,
		&route.IsCurrent, &route.RedirectToRouteID, &route.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *DomainRepo) GetCurrentRouteByPage(ctx context.Context, pageID int64) (*model.PageRoute, error) {
	var route model.PageRoute
	err := r.db.QueryRow(ctx, `
		SELECT id, page_id, domain_id, path, is_current, redirect_to_route_id, created_at
		FROM page_routes WHERE page_id = $1 AND is_current = true
	`, pageID).Scan(
		&route.ID, &route.PageID, &route.DomainID, &route.Path,
		&route.IsCurrent, &route.RedirectToRouteID, &route.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &route, nil
}
