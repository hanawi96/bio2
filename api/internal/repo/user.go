package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"linkbio/internal/model"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, username, display_name, avatar_asset_id, is_active, created_at, updated_at
	`, email, passwordHash).Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName, &user.AvatarAssetID,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, username, display_name, avatar_asset_id, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.DisplayName,
		&user.AvatarAssetID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, username, display_name, avatar_asset_id, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.DisplayName,
		&user.AvatarAssetID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, username, display_name, avatar_asset_id, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Username, &user.DisplayName,
		&user.AvatarAssetID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) SetUsername(ctx context.Context, userID int64, username string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET username = $2, updated_at = NOW() WHERE id = $1
	`, userID, username)
	return err
}

func (r *UserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
	`, username).Scan(&exists)
	return exists, err
}
