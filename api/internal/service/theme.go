package service

import (
	"context"
	"encoding/json"
	"fmt"

	"linkbio/internal/model"
	"linkbio/internal/repo"
	"linkbio/internal/util"
)

type ThemeService struct {
	themeRepo *repo.ThemeRepo
}

func NewThemeService(themeRepo *repo.ThemeRepo) *ThemeService {
	return &ThemeService{themeRepo: themeRepo}
}

func (s *ThemeService) ListPresets(ctx context.Context, tier string) ([]*model.ThemePreset, error) {
	return s.themeRepo.GetPresets(ctx, tier)
}

func (s *ThemeService) GetPreset(ctx context.Context, id int64) (*model.ThemePreset, error) {
	return s.themeRepo.GetPresetByID(ctx, id)
}

func (s *ThemeService) CreateOrUpdateCustom(ctx context.Context, userID, presetID int64, patch json.RawMessage) (*model.ThemeCustom, error) {
	// Generate hash for deduplication
	hashInput := fmt.Sprintf("%d:%s", presetID, string(patch))
	hash := util.SHA256(hashInput)

	// Check if exists
	existing, err := s.themeRepo.GetCustomByHash(ctx, userID, hash)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new
	return s.themeRepo.CreateCustom(ctx, userID, presetID, patch, hash)
}
