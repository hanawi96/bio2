package pages

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// GetThemePresetID returns preset ID by key
func (r *Repo) GetThemePresetID(ctx context.Context, key string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `SELECT id FROM theme_presets WHERE key = $1`, key).Scan(&id)
	return id, err
}

// GetThemePresetConfig returns preset config JSON
func (r *Repo) GetThemePresetConfig(ctx context.Context, presetID int64) (map[string]interface{}, error) {
	var configJSON []byte
	err := r.pool.QueryRow(ctx, `SELECT config FROM theme_presets WHERE id = $1`, presetID).Scan(&configJSON)
	if err != nil {
		return nil, err
	}
	var config map[string]interface{}
	if err := json.Unmarshal(configJSON, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// CreatePage creates a new page
func (r *Repo) CreatePage(ctx context.Context, userID, themePresetID int64, title, locale string) (int64, error) {
	var pageID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bio_pages (user_id, title, locale, theme_preset_id, theme_mode, settings)
		VALUES ($1, $2, $3, $4, 'light', '{}')
		RETURNING id
	`, userID, title, locale, themePresetID).Scan(&pageID)
	return pageID, err
}

// ListPages returns all pages for a user
func (r *Repo) ListPages(ctx context.Context, userID int64) ([]PageListItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(title, ''), locale, status, updated_at
		FROM bio_pages
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []PageListItem
	for rows.Next() {
		var p PageListItem
		if err := rows.Scan(&p.ID, &p.Title, &p.Locale, &p.Status, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	if pages == nil {
		pages = []PageListItem{}
	}
	return pages, nil
}

// GetPageBasic returns basic page info
func (r *Repo) GetPageBasic(ctx context.Context, pageID int64) (*DraftPage, int64, error) {
	var page DraftPage
	var themePresetID int64
	var themeCustomID *int64
	var settingsJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, COALESCE(title, ''), locale, status, access_type,
		       theme_preset_id, theme_custom_id, theme_mode, settings
		FROM bio_pages WHERE id = $1
	`, pageID).Scan(
		&page.ID, &page.Title, &page.Locale, &page.Status, &page.AccessType,
		&themePresetID, &themeCustomID, &page.Theme.Mode, &settingsJSON,
	)
	if err != nil {
		return nil, 0, err
	}

	// Get preset key
	r.pool.QueryRow(ctx, `SELECT key FROM theme_presets WHERE id = $1`, themePresetID).Scan(&page.Theme.PresetKey)
	page.Theme.CustomID = themeCustomID

	// Parse settings
	if len(settingsJSON) > 0 {
		json.Unmarshal(settingsJSON, &page.Settings)
	}

	return &page, themePresetID, nil
}

// GetBlocks returns all blocks for a page
func (r *Repo) GetBlocks(ctx context.Context, pageID int64) ([]DraftBlock, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, type, sort_key, ref_id, content, is_visible
		FROM blocks
		WHERE page_id = $1
		ORDER BY sort_key
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []DraftBlock
	for rows.Next() {
		var b DraftBlock
		var contentJSON []byte
		if err := rows.Scan(&b.ID, &b.Type, &b.SortKey, &b.RefID, &contentJSON, &b.IsVisible); err != nil {
			return nil, err
		}
		b.Content = make(map[string]interface{})
		json.Unmarshal(contentJSON, &b.Content)
		blocks = append(blocks, b)
	}
	if blocks == nil {
		blocks = []DraftBlock{}
	}
	return blocks, nil
}

// GetLinkGroups returns all link groups for a page
func (r *Repo) GetLinkGroups(ctx context.Context, pageID int64) (map[string]DraftGroup, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(title, ''), layout_type, layout_config, style_override
		FROM link_groups
		WHERE page_id = $1
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make(map[string]DraftGroup)
	for rows.Next() {
		var g DraftGroup
		var layoutConfigJSON, styleOverrideJSON []byte
		if err := rows.Scan(&g.ID, &g.Title, &g.LayoutType, &layoutConfigJSON, &styleOverrideJSON); err != nil {
			return nil, err
		}
		g.LayoutConfig = make(map[string]interface{})
		g.StyleOverride = make(map[string]interface{})
		json.Unmarshal(layoutConfigJSON, &g.LayoutConfig)
		if styleOverrideJSON != nil {
			json.Unmarshal(styleOverrideJSON, &g.StyleOverride)
		}
		groups[string(rune(g.ID))] = g
	}
	return groups, nil
}

// GetLinks returns all links grouped by group_id
func (r *Repo) GetLinks(ctx context.Context, pageID int64) (map[string][]DraftLink, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, l.group_id, l.title, l.url, l.icon_asset_id, l.sort_key, l.is_active
		FROM links l
		JOIN link_groups lg ON lg.id = l.group_id
		WHERE lg.page_id = $1
		ORDER BY l.sort_key
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make(map[string][]DraftLink)
	for rows.Next() {
		var l DraftLink
		var groupID int64
		if err := rows.Scan(&l.ID, &groupID, &l.Title, &l.URL, &l.IconAssetID, &l.SortKey, &l.IsActive); err != nil {
			return nil, err
		}
		key := string(rune(groupID))
		links[key] = append(links[key], l)
	}
	return links, nil
}

// SavePageSettings updates page settings and theme
func (r *Repo) SavePageSettings(ctx context.Context, pageID int64, page *DraftPage, themePresetID int64) error {
	settingsJSON, _ := json.Marshal(page.Settings)
	_, err := r.pool.Exec(ctx, `
		UPDATE bio_pages
		SET title = $2, locale = $3, theme_preset_id = $4, theme_custom_id = $5,
		    theme_mode = $6, settings = $7, updated_at = NOW()
		WHERE id = $1
	`, pageID, page.Title, page.Locale, themePresetID, page.Theme.CustomID, page.Theme.Mode, settingsJSON)
	return err
}

// DeleteBlocks deletes all blocks for a page
func (r *Repo) DeleteBlocks(ctx context.Context, pageID int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM blocks WHERE page_id = $1`, pageID)
	return err
}

// DeleteLinkGroups deletes all link groups (and cascades to links)
func (r *Repo) DeleteLinkGroups(ctx context.Context, pageID int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM link_groups WHERE page_id = $1`, pageID)
	return err
}

// InsertLinkGroup inserts a link group
func (r *Repo) InsertLinkGroup(ctx context.Context, pageID int64, g *DraftGroup) (int64, error) {
	layoutConfigJSON, _ := json.Marshal(g.LayoutConfig)
	styleOverrideJSON, _ := json.Marshal(g.StyleOverride)

	var newID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO link_groups (page_id, title, layout_type, layout_config, style_override)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, pageID, g.Title, g.LayoutType, layoutConfigJSON, styleOverrideJSON).Scan(&newID)
	return newID, err
}

// InsertBlock inserts a block
func (r *Repo) InsertBlock(ctx context.Context, pageID int64, b *DraftBlock) (int64, error) {
	contentJSON, _ := json.Marshal(b.Content)

	var newID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO blocks (page_id, type, sort_key, ref_id, content, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, pageID, b.Type, b.SortKey, b.RefID, contentJSON, b.IsVisible).Scan(&newID)
	return newID, err
}

// InsertLink inserts a link
func (r *Repo) InsertLink(ctx context.Context, groupID int64, l *DraftLink) (int64, error) {
	var newID int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO links (group_id, title, url, icon_asset_id, sort_key, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, groupID, l.Title, l.URL, l.IconAssetID, l.SortKey, l.IsActive).Scan(&newID)
	return newID, err
}

// UpdatePageStatus updates page status
func (r *Repo) UpdatePageStatus(ctx context.Context, pageID int64, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE bio_pages SET status = $2, updated_at = NOW() WHERE id = $1`, pageID, status)
	return err
}

// UpsertPublishCache inserts or updates publish cache
func (r *Repo) UpsertPublishCache(ctx context.Context, pageID int64, compiledJSON []byte, etag string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO page_publish_cache (page_id, compiled_json, etag, published_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (page_id) DO UPDATE SET
			compiled_json = EXCLUDED.compiled_json,
			etag = EXCLUDED.etag,
			published_at = NOW(),
			updated_at = NOW()
	`, pageID, compiledJSON, etag)
	return err
}

// GetPageUpdatedAt returns updated_at for a page
func (r *Repo) GetPageUpdatedAt(ctx context.Context, pageID int64) (time.Time, error) {
	var updatedAt time.Time
	err := r.pool.QueryRow(ctx, `SELECT updated_at FROM bio_pages WHERE id = $1`, pageID).Scan(&updatedAt)
	return updatedAt, err
}


// Pool returns the database pool for direct access
func (r *Repo) Pool() *pgxpool.Pool {
	return r.pool
}

// GetPublishCacheByPageID returns compiled_json for a page by ID
func (r *Repo) GetPublishCacheByPageID(ctx context.Context, pageID int64) (map[string]interface{}, error) {
	query := `SELECT compiled_json FROM page_publish_cache WHERE page_id = $1`

	var compiledJSON []byte
	err := r.pool.QueryRow(ctx, query, pageID).Scan(&compiledJSON)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(compiledJSON, &result); err != nil {
		return nil, err
	}

	return result, nil
}
