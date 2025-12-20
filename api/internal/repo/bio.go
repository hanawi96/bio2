package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type BioRepo struct {
	db *pgxpool.Pool
}

func NewBioRepo(db *pgxpool.Pool) *BioRepo {
	return &BioRepo{db: db}
}

func (r *BioRepo) GetOrCreatePage(ctx context.Context, userID int64) (*model.BioPage, error) {
	// Try to get existing page
	var page model.BioPage
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, locale, title, status, access_type,
		       theme_preset_id, theme_custom_id, theme_mode, settings, created_at, updated_at
		FROM bio_pages WHERE user_id = $1 LIMIT 1
	`, userID).Scan(
		&page.ID, &page.UserID, &page.Locale, &page.Title, &page.Status,
		&page.AccessType, &page.ThemePresetID, &page.ThemeCustomID,
		&page.ThemeMode, &page.Settings, &page.CreatedAt, &page.UpdatedAt,
	)

	if err == nil {
		return &page, nil
	}

	// Create new page
	var presetID int64
	err = r.db.QueryRow(ctx, `SELECT id FROM theme_presets LIMIT 1`).Scan(&presetID)
	if err != nil {
		presetID = 1
	}

	err = r.db.QueryRow(ctx, `
		INSERT INTO bio_pages (user_id, theme_preset_id, title, settings)
		VALUES ($1, $2, 'My Bio', '{}'::jsonb)
		RETURNING id, user_id, locale, title, status, access_type,
		          theme_preset_id, theme_custom_id, theme_mode, settings, created_at, updated_at
	`, userID, presetID).Scan(
		&page.ID, &page.UserID, &page.Locale, &page.Title, &page.Status,
		&page.AccessType, &page.ThemePresetID, &page.ThemeCustomID,
		&page.ThemeMode, &page.Settings, &page.CreatedAt, &page.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (r *BioRepo) GetBlockByID(ctx context.Context, id int64) (*model.Block, error) {
	var b model.Block
	err := r.db.QueryRow(ctx, `
		SELECT id, page_id, type, sort_key, ref_id, content, is_visible, created_at, updated_at
		FROM blocks WHERE id = $1
	`, id).Scan(&b.ID, &b.PageID, &b.Type, &b.SortKey, &b.RefID,
		&b.Content, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BioRepo) GetBlockOwnerID(ctx context.Context, blockID int64) (int64, error) {
	var ownerID int64
	err := r.db.QueryRow(ctx, `
		SELECT bp.user_id
		FROM blocks b
		JOIN bio_pages bp ON b.page_id = bp.id
		WHERE b.id = $1
	`, blockID).Scan(&ownerID)
	return ownerID, err
}

func (r *BioRepo) GetGroupOwnerID(ctx context.Context, groupID int64) (int64, error) {
	var ownerID int64
	err := r.db.QueryRow(ctx, `
		SELECT bp.user_id
		FROM link_groups lg
		JOIN bio_pages bp ON lg.page_id = bp.id
		WHERE lg.id = $1
	`, groupID).Scan(&ownerID)
	return ownerID, err
}

func (r *BioRepo) GetLinkOwnerID(ctx context.Context, linkID int64) (int64, error) {
	var ownerID int64
	err := r.db.QueryRow(ctx, `
		SELECT bp.user_id
		FROM links l
		JOIN link_groups lg ON l.group_id = lg.id
		JOIN bio_pages bp ON lg.page_id = bp.id
		WHERE l.id = $1
	`, linkID).Scan(&ownerID)
	return ownerID, err
}

func (r *BioRepo) GetLinkByID(ctx context.Context, id int64) (*model.Link, error) {
	var l model.Link
	err := r.db.QueryRow(ctx, `
		SELECT id, group_id, title, url, icon_asset_id, sort_key, is_active, created_at, updated_at
		FROM links WHERE id = $1
	`, id).Scan(&l.ID, &l.GroupID, &l.Title, &l.URL, &l.IconAssetID,
		&l.SortKey, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *BioRepo) GetLastBlockSortKey(ctx context.Context, pageID int64) (string, error) {
	var sortKey string
	err := r.db.QueryRow(ctx, `
		SELECT sort_key FROM blocks WHERE page_id = $1 ORDER BY sort_key DESC LIMIT 1
	`, pageID).Scan(&sortKey)
	if err != nil {
		return "", err
	}
	return sortKey, nil
}

func (r *BioRepo) GetLastLinkSortKey(ctx context.Context, groupID int64) (string, error) {
	var sortKey string
	err := r.db.QueryRow(ctx, `
		SELECT sort_key FROM links WHERE group_id = $1 ORDER BY sort_key DESC LIMIT 1
	`, groupID).Scan(&sortKey)
	if err != nil {
		return "", err
	}
	return sortKey, nil
}

func (r *BioRepo) UpdateBlockSortKey(ctx context.Context, blockID int64, sortKey string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE blocks SET sort_key = $2, updated_at = NOW() WHERE id = $1
	`, blockID, sortKey)
	return err
}
