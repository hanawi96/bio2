package repo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type ThemeRepo struct {
	db *pgxpool.Pool
}

func NewThemeRepo(db *pgxpool.Pool) *ThemeRepo {
	return &ThemeRepo{db: db}
}

func (r *ThemeRepo) GetPresets(ctx context.Context, tier string) ([]*model.ThemePreset, error) {
	query := `
		SELECT id, key, name, tier, visibility, is_official, author_user_id, config, created_at, updated_at
		FROM theme_presets WHERE visibility = 'public'
	`
	args := []interface{}{}

	if tier != "" {
		query += " AND tier = $1"
		args = append(args, tier)
	}
	query += " ORDER BY name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []*model.ThemePreset
	for rows.Next() {
		var p model.ThemePreset
		err := rows.Scan(&p.ID, &p.Key, &p.Name, &p.Tier, &p.Visibility,
			&p.IsOfficial, &p.AuthorUserID, &p.Config, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		presets = append(presets, &p)
	}
	return presets, nil
}

func (r *ThemeRepo) GetPresetByID(ctx context.Context, id int64) (*model.ThemePreset, error) {
	var p model.ThemePreset
	err := r.db.QueryRow(ctx, `
		SELECT id, key, name, tier, visibility, is_official, author_user_id, config, created_at, updated_at
		FROM theme_presets WHERE id = $1
	`, id).Scan(&p.ID, &p.Key, &p.Name, &p.Tier, &p.Visibility,
		&p.IsOfficial, &p.AuthorUserID, &p.Config, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ThemeRepo) GetPresetByKey(ctx context.Context, key string) (*model.ThemePreset, error) {
	var p model.ThemePreset
	err := r.db.QueryRow(ctx, `
		SELECT id, key, name, tier, visibility, is_official, author_user_id, config, created_at, updated_at
		FROM theme_presets WHERE key = $1
	`, key).Scan(&p.ID, &p.Key, &p.Name, &p.Tier, &p.Visibility,
		&p.IsOfficial, &p.AuthorUserID, &p.Config, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ThemeRepo) CreateCustom(ctx context.Context, userID, presetID int64, patch json.RawMessage, hash string) (*model.ThemeCustom, error) {
	var t model.ThemeCustom
	err := r.db.QueryRow(ctx, `
		INSERT INTO themes_custom (user_id, based_on_preset_id, patch, hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, based_on_preset_id, name, patch, compiled_config, hash, created_at, updated_at
	`, userID, presetID, patch, hash).Scan(
		&t.ID, &t.UserID, &t.BasedOnPresetID, &t.Name, &t.Patch,
		&t.CompiledConfig, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ThemeRepo) GetCustomByHash(ctx context.Context, userID int64, hash string) (*model.ThemeCustom, error) {
	var t model.ThemeCustom
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, based_on_preset_id, name, patch, compiled_config, hash, created_at, updated_at
		FROM themes_custom WHERE user_id = $1 AND hash = $2
	`, userID, hash).Scan(
		&t.ID, &t.UserID, &t.BasedOnPresetID, &t.Name, &t.Patch,
		&t.CompiledConfig, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ThemeRepo) GetCustomByID(ctx context.Context, id int64) (*model.ThemeCustom, error) {
	var t model.ThemeCustom
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, based_on_preset_id, name, patch, compiled_config, hash, created_at, updated_at
		FROM themes_custom WHERE id = $1
	`, id).Scan(
		&t.ID, &t.UserID, &t.BasedOnPresetID, &t.Name, &t.Patch,
		&t.CompiledConfig, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ThemeRepo) UpdateCustom(ctx context.Context, id int64, patch, compiled json.RawMessage, hash string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE themes_custom SET patch = $2, compiled_config = $3, hash = $4, updated_at = NOW()
		WHERE id = $1
	`, id, patch, compiled, hash)
	return err
}

func (r *ThemeRepo) GetCustomByUserID(ctx context.Context, userID int64) (*model.ThemeCustom, error) {
	var t model.ThemeCustom
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, based_on_preset_id, name, patch, compiled_config, hash, created_at, updated_at
		FROM themes_custom WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT 1
	`, userID).Scan(
		&t.ID, &t.UserID, &t.BasedOnPresetID, &t.Name, &t.Patch,
		&t.CompiledConfig, &t.Hash, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *ThemeRepo) DeleteCustom(ctx context.Context, id, userID int64) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM themes_custom WHERE id = $1 AND user_id = $2
	`, id, userID)
	return err
}
