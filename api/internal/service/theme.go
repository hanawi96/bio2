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

	// Check if user already has a custom theme
	existing, err := s.themeRepo.GetCustomByUserID(ctx, userID)
	if err == nil && existing != nil {
		// Update existing custom theme
		compiled := s.compileTheme(ctx, presetID, patch)
		err = s.themeRepo.UpdateCustom(ctx, existing.ID, patch, compiled, hash)
		if err != nil {
			return nil, err
		}
		// Return updated theme
		return s.themeRepo.GetCustomByID(ctx, existing.ID)
	}

	// Create new custom theme
	compiled := s.compileTheme(ctx, presetID, patch)
	custom, err := s.themeRepo.CreateCustom(ctx, userID, presetID, patch, hash)
	if err != nil {
		return nil, err
	}
	
	// Update compiled config
	if compiled != nil {
		_ = s.themeRepo.UpdateCustom(ctx, custom.ID, patch, compiled, hash)
	}
	
	return custom, nil
}

func (s *ThemeService) GetUserCustomTheme(ctx context.Context, userID int64) (*model.ThemeCustom, error) {
	return s.themeRepo.GetCustomByUserID(ctx, userID)
}

func (s *ThemeService) DeleteCustomTheme(ctx context.Context, id, userID int64) error {
	return s.themeRepo.DeleteCustom(ctx, id, userID)
}

// compileTheme merges preset with patch (simplified version)
func (s *ThemeService) compileTheme(ctx context.Context, presetID int64, patch json.RawMessage) json.RawMessage {
	preset, err := s.themeRepo.GetPresetByID(ctx, presetID)
	if err != nil {
		return nil
	}
	
	// Deep merge preset.Config with patch
	var presetMap map[string]interface{}
	var patchMap map[string]interface{}
	
	if err := json.Unmarshal(preset.Config, &presetMap); err != nil {
		return nil
	}
	if err := json.Unmarshal(patch, &patchMap); err != nil {
		return nil
	}
	
	merged := deepMerge(presetMap, patchMap)
	compiled, _ := json.Marshal(merged)
	return compiled
}

func deepMerge(dst, src map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range src {
		if dstVal, ok := result[k]; ok {
			if dstMap, dstOk := dstVal.(map[string]interface{}); dstOk {
				if srcMap, srcOk := v.(map[string]interface{}); srcOk {
					result[k] = deepMerge(dstMap, srcMap)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}
