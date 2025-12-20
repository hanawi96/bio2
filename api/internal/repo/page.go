package repo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type PageRepo struct {
	db *pgxpool.Pool
}

func NewPageRepo(db *pgxpool.Pool) *PageRepo {
	return &PageRepo{db: db}
}

func (r *PageRepo) Create(ctx context.Context, userID, themePresetID int64, title string) (*model.BioPage, error) {
	var page model.BioPage
	err := r.db.QueryRow(ctx, `
		INSERT INTO bio_pages (user_id, theme_preset_id, title, settings)
		VALUES ($1, $2, $3, '{}'::jsonb)
		RETURNING id, user_id, locale, title, status, access_type, theme_preset_id, 
		          theme_custom_id, theme_mode, settings, created_at, updated_at
	`, userID, themePresetID, title).Scan(
		&page.ID, &page.UserID, &page.Locale, &page.Title, &page.Status,
		&page.AccessType, &page.ThemePresetID, &page.ThemeCustomID,
		&page.ThemeMode, &page.Settings, &page.CreatedAt, &page.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) GetByID(ctx context.Context, id int64) (*model.BioPage, error) {
	var page model.BioPage
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, locale, title, status, access_type, password_hash,
		       theme_preset_id, theme_custom_id, theme_mode, settings, created_at, updated_at
		FROM bio_pages WHERE id = $1
	`, id).Scan(
		&page.ID, &page.UserID, &page.Locale, &page.Title, &page.Status,
		&page.AccessType, &page.PasswordHash, &page.ThemePresetID, &page.ThemeCustomID,
		&page.ThemeMode, &page.Settings, &page.CreatedAt, &page.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) ListByUser(ctx context.Context, userID int64) ([]*model.BioPage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, locale, title, status, access_type,
		       theme_preset_id, theme_custom_id, theme_mode, settings, created_at, updated_at
		FROM bio_pages WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*model.BioPage
	for rows.Next() {
		var page model.BioPage
		err := rows.Scan(
			&page.ID, &page.UserID, &page.Locale, &page.Title, &page.Status,
			&page.AccessType, &page.ThemePresetID, &page.ThemeCustomID,
			&page.ThemeMode, &page.Settings, &page.CreatedAt, &page.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pages = append(pages, &page)
	}
	return pages, nil
}

func (r *PageRepo) Update(ctx context.Context, page *model.BioPage) error {
	_, err := r.db.Exec(ctx, `
		UPDATE bio_pages SET
			locale = $2, title = $3, status = $4, access_type = $5,
			theme_preset_id = $6, theme_custom_id = $7, theme_mode = $8,
			settings = $9, updated_at = NOW()
		WHERE id = $1
	`, page.ID, page.Locale, page.Title, page.Status, page.AccessType,
		page.ThemePresetID, page.ThemeCustomID, page.ThemeMode, page.Settings)
	return err
}

func (r *PageRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM bio_pages WHERE id = $1`, id)
	return err
}

func (r *PageRepo) SavePublishCache(ctx context.Context, pageID int64, compiled json.RawMessage) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO page_publish_cache (page_id, compiled_json)
		VALUES ($1, $2)
		ON CONFLICT (page_id) DO UPDATE SET
			compiled_json = $2, updated_at = NOW()
	`, pageID, compiled)
	return err
}

func (r *PageRepo) GetPublishCache(ctx context.Context, pageID int64) (*model.PagePublishCache, error) {
	var cache model.PagePublishCache
	err := r.db.QueryRow(ctx, `
		SELECT page_id, compiled_json, published_at, updated_at
		FROM page_publish_cache WHERE page_id = $1
	`, pageID).Scan(&cache.PageID, &cache.CompiledJSON, &cache.PublishedAt, &cache.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}
