package pages

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
)

var ErrNotFound = errors.New("not found")

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// GetRepo returns the repo for direct access when needed
func (s *Service) GetRepo() *Repo {
	return s.repo
}

// CreatePage creates a new page
func (s *Service) CreatePage(ctx context.Context, userID int64, req *CreatePageRequest) (int64, error) {
	// Default values
	locale := req.Locale
	if locale == "" {
		locale = "vi"
	}
	themeKey := req.ThemePresetKey
	if themeKey == "" {
		themeKey = "theme_a"
	}

	// Get theme preset ID
	themePresetID, err := s.repo.GetThemePresetID(ctx, themeKey)
	if err != nil {
		return 0, fmt.Errorf("theme preset not found: %s", themeKey)
	}

	return s.repo.CreatePage(ctx, userID, themePresetID, req.Title, locale)
}

// ListPages returns all pages for a user
func (s *Service) ListPages(ctx context.Context, userID int64) ([]PageListItem, error) {
	return s.repo.ListPages(ctx, userID)
}

// GetDraft returns the draft state for a page
func (s *Service) GetDraft(ctx context.Context, pageID int64) (*DraftPayload, error) {
	// Get page basic info
	page, _, err := s.repo.GetPageBasic(ctx, pageID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Get blocks
	blocks, err := s.repo.GetBlocks(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Get link groups
	groupsRaw, err := s.repo.GetLinkGroups(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Get links
	linksRaw, err := s.repo.GetLinks(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Convert to string keys (group ID as string)
	linkGroups := make(map[string]DraftGroup)
	for _, g := range groupsRaw {
		linkGroups[strconv.FormatInt(g.ID, 10)] = g
	}

	links := make(map[string][]DraftLink)
	for k, v := range linksRaw {
		// k is already string representation of group ID
		links[k] = v
	}

	// Fix: re-fetch groups with proper string keys
	linkGroups = make(map[string]DraftGroup)
	rows, _ := s.repo.Pool().Query(ctx, `
		SELECT id, COALESCE(title, ''), layout_type, layout_config, style_override
		FROM link_groups WHERE page_id = $1
	`, pageID)
	defer rows.Close()
	for rows.Next() {
		var g DraftGroup
		var layoutConfigJSON, styleOverrideJSON []byte
		rows.Scan(&g.ID, &g.Title, &g.LayoutType, &layoutConfigJSON, &styleOverrideJSON)
		g.LayoutConfig = make(map[string]interface{})
		g.StyleOverride = make(map[string]interface{})
		json.Unmarshal(layoutConfigJSON, &g.LayoutConfig)
		if styleOverrideJSON != nil {
			json.Unmarshal(styleOverrideJSON, &g.StyleOverride)
		}
		linkGroups[strconv.FormatInt(g.ID, 10)] = g
	}

	// Re-fetch links with proper string keys
	links = make(map[string][]DraftLink)
	linkRows, _ := s.repo.Pool().Query(ctx, `
		SELECT l.id, l.group_id, l.title, l.url, l.icon_asset_id, l.sort_key, l.is_active
		FROM links l
		JOIN link_groups lg ON lg.id = l.group_id
		WHERE lg.page_id = $1
		ORDER BY l.sort_key
	`, pageID)
	defer linkRows.Close()
	for linkRows.Next() {
		var l DraftLink
		var groupID int64
		linkRows.Scan(&l.ID, &groupID, &l.Title, &l.URL, &l.IconAssetID, &l.SortKey, &l.IsActive)
		key := strconv.FormatInt(groupID, 10)
		links[key] = append(links[key], l)
	}

	return &DraftPayload{
		Page:       *page,
		Blocks:     blocks,
		LinkGroups: linkGroups,
		Links:      links,
	}, nil
}

// SaveDraft saves the draft state
func (s *Service) SaveDraft(ctx context.Context, pageID int64, draft *DraftPayload) (time.Time, error) {
	// Verify page exists
	_, themePresetID, err := s.repo.GetPageBasic(ctx, pageID)
	if err != nil {
		return time.Time{}, ErrNotFound
	}

	// Get theme preset ID from key if changed
	if draft.Page.Theme.PresetKey != "" {
		newPresetID, err := s.repo.GetThemePresetID(ctx, draft.Page.Theme.PresetKey)
		if err == nil {
			themePresetID = newPresetID
		}
	}

	// Save page settings
	if err := s.repo.SavePageSettings(ctx, pageID, &draft.Page, themePresetID); err != nil {
		return time.Time{}, err
	}

	// Delete existing blocks and groups (cascade deletes links)
	if err := s.repo.DeleteBlocks(ctx, pageID); err != nil {
		return time.Time{}, err
	}
	if err := s.repo.DeleteLinkGroups(ctx, pageID); err != nil {
		return time.Time{}, err
	}

	// Map old group IDs to new IDs
	groupIDMap := make(map[string]int64)

	// Insert link groups first
	for oldIDStr, g := range draft.LinkGroups {
		newID, err := s.repo.InsertLinkGroup(ctx, pageID, &g)
		if err != nil {
			return time.Time{}, err
		}
		groupIDMap[oldIDStr] = newID
	}

	// Insert blocks with updated ref_id
	for _, b := range draft.Blocks {
		block := b
		if block.Type == "link_group" && block.RefID != nil {
			oldRefIDStr := strconv.FormatInt(*block.RefID, 10)
			if newID, ok := groupIDMap[oldRefIDStr]; ok {
				block.RefID = &newID
			}
		}
		if _, err := s.repo.InsertBlock(ctx, pageID, &block); err != nil {
			return time.Time{}, err
		}
	}

	// Insert links
	for oldGroupIDStr, linkList := range draft.Links {
		newGroupID, ok := groupIDMap[oldGroupIDStr]
		if !ok {
			continue
		}
		for _, l := range linkList {
			if _, err := s.repo.InsertLink(ctx, newGroupID, &l); err != nil {
				return time.Time{}, err
			}
		}
	}

	return s.repo.GetPageUpdatedAt(ctx, pageID)
}

// Publish compiles and publishes the page
func (s *Service) Publish(ctx context.Context, pageID int64) (time.Time, error) {
	// Get draft
	draft, err := s.GetDraft(ctx, pageID)
	if err != nil {
		return time.Time{}, err
	}

	// Get theme preset config
	themePresetID, _ := s.repo.GetThemePresetID(ctx, draft.Page.Theme.PresetKey)
	themeConfig, _ := s.repo.GetThemePresetConfig(ctx, themePresetID)

	// Compile
	compiled := s.compile(draft, themeConfig)

	// Serialize
	compiledJSON, err := json.Marshal(compiled)
	if err != nil {
		return time.Time{}, err
	}

	// Generate ETag
	hash := sha256.Sum256(compiledJSON)
	etag := hex.EncodeToString(hash[:16])

	// Save to cache
	if err := s.repo.UpsertPublishCache(ctx, pageID, compiledJSON, etag); err != nil {
		return time.Time{}, err
	}

	// Update page status
	if err := s.repo.UpdatePageStatus(ctx, pageID, "published"); err != nil {
		return time.Time{}, err
	}

	return time.Now(), nil
}

// compile transforms draft to compiled payload
func (s *Service) compile(draft *DraftPayload, themeConfig map[string]interface{}) map[string]interface{} {
	// Sort blocks by sort_key
	sortedBlocks := make([]DraftBlock, len(draft.Blocks))
	copy(sortedBlocks, draft.Blocks)
	sort.Slice(sortedBlocks, func(i, j int) bool {
		return sortedBlocks[i].SortKey < sortedBlocks[j].SortKey
	})

	// Build compiled blocks
	compiledBlocks := make([]map[string]interface{}, 0)
	for _, block := range sortedBlocks {
		if !block.IsVisible {
			continue
		}

		if block.Type == "text" {
			compiledBlocks = append(compiledBlocks, map[string]interface{}{
				"type":    "text",
				"content": block.Content,
			})
		} else if block.Type == "link_group" && block.RefID != nil {
			groupIDStr := strconv.FormatInt(*block.RefID, 10)
			group, ok := draft.LinkGroups[groupIDStr]
			if !ok {
				continue
			}

			// Get links for this group and sort
			groupLinks := draft.Links[groupIDStr]
			sort.Slice(groupLinks, func(i, j int) bool {
				return groupLinks[i].SortKey < groupLinks[j].SortKey
			})

			// Build compiled links
			compiledLinks := make([]map[string]interface{}, 0)
			for _, link := range groupLinks {
				if !link.IsActive {
					continue
				}
				compiledLinks = append(compiledLinks, map[string]interface{}{
					"title":    link.Title,
					"url":      link.URL,
					"icon_url": nil,
				})
			}

			// Merge final_style: theme defaults + group override
			finalStyle := s.mergeFinalStyle(themeConfig, group.StyleOverride)

			compiledBlocks = append(compiledBlocks, map[string]interface{}{
				"type": "link_group",
				"group": map[string]interface{}{
					"id":            group.ID,
					"title":         group.Title,
					"layout_type":   group.LayoutType,
					"layout_config": group.LayoutConfig,
					"final_style":   finalStyle,
					"links":         compiledLinks,
				},
			})
		}
	}

	// Build compiled theme
	compiledTheme := s.buildCompiledTheme(themeConfig, draft.Page.Theme.Mode)

	// Build settings
	settings := make(map[string]interface{})
	if draft.Page.Settings.Header != nil {
		settings["header"] = map[string]interface{}{
			"avatar_asset_id": draft.Page.Settings.Header.AvatarAssetID,
			"name":            draft.Page.Settings.Header.Name,
			"bio":             draft.Page.Settings.Header.Bio,
			"verified":        draft.Page.Settings.Header.Verified,
			"social":          draft.Page.Settings.Header.Social,
		}
	}
	if draft.Page.Settings.Cover != nil {
		settings["cover"] = map[string]interface{}{
			"kind":     draft.Page.Settings.Cover.Kind,
			"color":    draft.Page.Settings.Cover.Color,
			"asset_id": draft.Page.Settings.Cover.AssetID,
		}
	}

	return map[string]interface{}{
		"version": 1,
		"page": map[string]interface{}{
			"id":       draft.Page.ID,
			"locale":   draft.Page.Locale,
			"settings": settings,
		},
		"theme": map[string]interface{}{
			"preset_key": draft.Page.Theme.PresetKey,
			"mode":       draft.Page.Theme.Mode,
			"compiled":   compiledTheme,
		},
		"blocks": compiledBlocks,
	}
}

// mergeFinalStyle merges theme defaults with group override
func (s *Service) mergeFinalStyle(themeConfig map[string]interface{}, styleOverride map[string]interface{}) map[string]interface{} {
	finalStyle := make(map[string]interface{})

	// Get theme defaults for linkGroup
	if page, ok := themeConfig["page"].(map[string]interface{}); ok {
		if defaults, ok := page["defaults"].(map[string]interface{}); ok {
			if linkGroup, ok := defaults["linkGroup"].(map[string]interface{}); ok {
				for k, v := range linkGroup {
					finalStyle[k] = v
				}
			}
		}
	}

	// Override with group style
	for k, v := range styleOverride {
		finalStyle[k] = v
	}

	return finalStyle
}

// buildCompiledTheme extracts compiled theme from config
func (s *Service) buildCompiledTheme(themeConfig map[string]interface{}, mode string) map[string]interface{} {
	compiled := make(map[string]interface{})

	// Extract page layout and defaults
	if page, ok := themeConfig["page"].(map[string]interface{}); ok {
		compiled["page"] = page
	} else {
		// Minimal defaults
		compiled["page"] = map[string]interface{}{
			"layout": map[string]interface{}{
				"textAlign":    "center",
				"baseFontSize": "M",
				"pagePadding":  16,
				"blockGap":     12,
			},
			"defaults": map[string]interface{}{
				"linkGroup": map[string]interface{}{
					"textAlign": "center",
					"fontSize":  "M",
					"radius":    16,
				},
			},
		}
	}

	// Extract background
	if bg, ok := themeConfig["background"].(map[string]interface{}); ok {
		compiled["background"] = bg
	} else {
		compiled["background"] = map[string]interface{}{
			"kind":  "color",
			"color": "#0B0F19",
		}
	}

	return compiled
}

// GetPublished returns the published compiled_json for a page
func (s *Service) GetPublished(ctx context.Context, pageID int64) (map[string]interface{}, error) {
	result, err := s.repo.GetPublishCacheByPageID(ctx, pageID)
	if err != nil {
		return nil, ErrNotFound
	}
	return result, nil
}
