package service

import (
	"context"
	"encoding/json"

	"linkbio/internal/repo"
)

type CompilerService struct {
	pageRepo  *repo.PageRepo
	blockRepo *repo.BlockRepo
	themeRepo *repo.ThemeRepo
}

func NewCompilerService(pageRepo *repo.PageRepo, blockRepo *repo.BlockRepo, themeRepo *repo.ThemeRepo) *CompilerService {
	return &CompilerService{
		pageRepo:  pageRepo,
		blockRepo: blockRepo,
		themeRepo: themeRepo,
	}
}

type CompiledPage struct {
	Page   CompiledPageInfo   `json:"page"`
	Theme  json.RawMessage    `json:"theme"`
	Blocks []CompiledBlock    `json:"blocks"`
}

type CompiledPageInfo struct {
	ID       int64           `json:"id"`
	Title    *string         `json:"title"`
	Locale   string          `json:"locale"`
	Mode     string          `json:"mode"`
	Settings json.RawMessage `json:"settings"`
}

type CompiledBlock struct {
	ID        int64           `json:"id"`
	Type      string          `json:"type"`
	Content   json.RawMessage `json:"content"`
	IsVisible bool            `json:"is_visible"`
	// For link_group blocks
	Group *CompiledLinkGroup `json:"group,omitempty"`
}

type CompiledLinkGroup struct {
	ID            int64           `json:"id"`
	Title         *string         `json:"title"`
	LayoutType    string          `json:"layout_type"`
	LayoutConfig  json.RawMessage `json:"layout_config"`
	FinalStyle    json.RawMessage `json:"final_style"` // merged theme defaults + override
	Links         []CompiledLink  `json:"links"`
}

type CompiledLink struct {
	ID       int64   `json:"id"`
	Title    string  `json:"title"`
	URL      string  `json:"url"`
	IsActive bool    `json:"is_active"`
}

func (s *CompilerService) Compile(ctx context.Context, pageID int64) (*CompiledPage, error) {
	// Get page
	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Get theme config
	var themeConfig json.RawMessage
	if page.ThemeCustomID != nil {
		custom, err := s.themeRepo.GetCustomByID(ctx, *page.ThemeCustomID)
		if err != nil {
			return nil, err
		}
		if custom.CompiledConfig != nil {
			themeConfig = custom.CompiledConfig
		} else {
			// Fallback to preset + patch merge (simplified)
			preset, err := s.themeRepo.GetPresetByID(ctx, custom.BasedOnPresetID)
			if err != nil {
				return nil, err
			}
			themeConfig = preset.Config // TODO: merge with patch
		}
	} else {
		preset, err := s.themeRepo.GetPresetByID(ctx, page.ThemePresetID)
		if err != nil {
			return nil, err
		}
		themeConfig = preset.Config
	}

	// Get blocks sorted
	blocks, err := s.blockRepo.GetBlocksByPage(ctx, pageID)
	if err != nil {
		return nil, err
	}

	// Get link groups
	groups, err := s.blockRepo.GetLinkGroupsByPage(ctx, pageID)
	if err != nil {
		return nil, err
	}
	groupMap := make(map[int64]*CompiledLinkGroup)
	for _, g := range groups {
		links, err := s.blockRepo.GetLinksByGroup(ctx, g.ID)
		if err != nil {
			return nil, err
		}

		compiledLinks := make([]CompiledLink, 0, len(links))
		for _, l := range links {
			if l.IsActive {
				compiledLinks = append(compiledLinks, CompiledLink{
					ID:       l.ID,
					Title:    l.Title,
					URL:      l.URL,
					IsActive: l.IsActive,
				})
			}
		}

		// Merge style: theme defaults + group override
		finalStyle := g.StyleOverride
		if finalStyle == nil {
			finalStyle = json.RawMessage(`{}`)
		}

		groupMap[g.ID] = &CompiledLinkGroup{
			ID:           g.ID,
			Title:        g.Title,
			LayoutType:   g.LayoutType,
			LayoutConfig: g.LayoutConfig,
			FinalStyle:   finalStyle,
			Links:        compiledLinks,
		}
	}

	// Build compiled blocks
	compiledBlocks := make([]CompiledBlock, 0, len(blocks))
	for _, b := range blocks {
		if !b.IsVisible {
			continue
		}

		cb := CompiledBlock{
			ID:        b.ID,
			Type:      b.Type,
			Content:   b.Content,
			IsVisible: b.IsVisible,
		}

		if b.Type == "link_group" && b.RefID != nil {
			if group, ok := groupMap[*b.RefID]; ok {
				cb.Group = group
			}
		}

		compiledBlocks = append(compiledBlocks, cb)
	}

	compiled := &CompiledPage{
		Page: CompiledPageInfo{
			ID:       page.ID,
			Title:    page.Title,
			Locale:   page.Locale,
			Mode:     page.ThemeMode,
			Settings: page.Settings,
		},
		Theme:  themeConfig,
		Blocks: compiledBlocks,
	}

	return compiled, nil
}

func (s *CompilerService) Publish(ctx context.Context, pageID int64) error {
	compiled, err := s.Compile(ctx, pageID)
	if err != nil {
		return err
	}

	compiledJSON, err := json.Marshal(compiled)
	if err != nil {
		return err
	}

	// Update page status
	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return err
	}
	page.Status = "published"
	if err := s.pageRepo.Update(ctx, page); err != nil {
		return err
	}

	// Save to cache
	return s.pageRepo.SavePublishCache(ctx, pageID, compiledJSON)
}
