package repo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type BlockRepo struct {
	db *pgxpool.Pool
}

func NewBlockRepo(db *pgxpool.Pool) *BlockRepo {
	return &BlockRepo{db: db}
}

// Blocks
func (r *BlockRepo) CreateBlock(ctx context.Context, pageID int64, blockType, sortKey string, refID *int64, content json.RawMessage) (*model.Block, error) {
	var block model.Block
	err := r.db.QueryRow(ctx, `
		INSERT INTO blocks (page_id, type, sort_key, ref_id, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, page_id, type, sort_key, ref_id, content, is_visible, created_at, updated_at
	`, pageID, blockType, sortKey, refID, content).Scan(
		&block.ID, &block.PageID, &block.Type, &block.SortKey, &block.RefID,
		&block.Content, &block.IsVisible, &block.CreatedAt, &block.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *BlockRepo) GetBlocksByPage(ctx context.Context, pageID int64) ([]*model.Block, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, page_id, type, sort_key, ref_id, content, is_visible, created_at, updated_at
		FROM blocks WHERE page_id = $1 ORDER BY sort_key
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*model.Block
	for rows.Next() {
		var b model.Block
		err := rows.Scan(&b.ID, &b.PageID, &b.Type, &b.SortKey, &b.RefID,
			&b.Content, &b.IsVisible, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &b)
	}
	return blocks, nil
}

func (r *BlockRepo) UpdateBlock(ctx context.Context, block *model.Block) error {
	_, err := r.db.Exec(ctx, `
		UPDATE blocks SET sort_key = $2, content = $3, is_visible = $4, updated_at = NOW()
		WHERE id = $1
	`, block.ID, block.SortKey, block.Content, block.IsVisible)
	return err
}

func (r *BlockRepo) DeleteBlock(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM blocks WHERE id = $1`, id)
	return err
}

// Link Groups
func (r *BlockRepo) CreateLinkGroup(ctx context.Context, pageID int64, title *string, layoutType string) (*model.LinkGroup, error) {
	var group model.LinkGroup
	err := r.db.QueryRow(ctx, `
		INSERT INTO link_groups (page_id, title, layout_type)
		VALUES ($1, $2, $3)
		RETURNING id, page_id, title, layout_type, layout_config, style_override, created_at, updated_at
	`, pageID, title, layoutType).Scan(
		&group.ID, &group.PageID, &group.Title, &group.LayoutType,
		&group.LayoutConfig, &group.StyleOverride, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *BlockRepo) GetLinkGroupsByPage(ctx context.Context, pageID int64) ([]*model.LinkGroup, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, page_id, title, layout_type, layout_config, style_override, created_at, updated_at
		FROM link_groups WHERE page_id = $1
	`, pageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.LinkGroup
	for rows.Next() {
		var g model.LinkGroup
		err := rows.Scan(&g.ID, &g.PageID, &g.Title, &g.LayoutType,
			&g.LayoutConfig, &g.StyleOverride, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}

func (r *BlockRepo) UpdateLinkGroup(ctx context.Context, group *model.LinkGroup) error {
	_, err := r.db.Exec(ctx, `
		UPDATE link_groups SET title = $2, layout_type = $3, layout_config = $4, 
		       style_override = $5, updated_at = NOW()
		WHERE id = $1
	`, group.ID, group.Title, group.LayoutType, group.LayoutConfig, group.StyleOverride)
	return err
}

func (r *BlockRepo) DeleteLinkGroup(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM link_groups WHERE id = $1`, id)
	return err
}

// Links
func (r *BlockRepo) CreateLink(ctx context.Context, groupID int64, title, url, sortKey string) (*model.Link, error) {
	var link model.Link
	err := r.db.QueryRow(ctx, `
		INSERT INTO links (group_id, title, url, sort_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id, group_id, title, url, icon_asset_id, sort_key, is_active, created_at, updated_at
	`, groupID, title, url, sortKey).Scan(
		&link.ID, &link.GroupID, &link.Title, &link.URL, &link.IconAssetID,
		&link.SortKey, &link.IsActive, &link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *BlockRepo) GetLinksByGroup(ctx context.Context, groupID int64) ([]*model.Link, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, group_id, title, url, icon_asset_id, sort_key, is_active, created_at, updated_at
		FROM links WHERE group_id = $1 ORDER BY sort_key
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*model.Link
	for rows.Next() {
		var l model.Link
		err := rows.Scan(&l.ID, &l.GroupID, &l.Title, &l.URL, &l.IconAssetID,
			&l.SortKey, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, &l)
	}
	return links, nil
}

func (r *BlockRepo) UpdateLink(ctx context.Context, link *model.Link) error {
	_, err := r.db.Exec(ctx, `
		UPDATE links SET title = $2, url = $3, sort_key = $4, is_active = $5, updated_at = NOW()
		WHERE id = $1
	`, link.ID, link.Title, link.URL, link.SortKey, link.IsActive)
	return err
}

func (r *BlockRepo) DeleteLink(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM links WHERE id = $1`, id)
	return err
}
