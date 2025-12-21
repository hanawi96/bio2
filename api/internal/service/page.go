package service

import (
	"context"
	"encoding/json"

	"linkbio/internal/model"
	"linkbio/internal/repo"
	"linkbio/internal/util"
)

type PageService struct {
	pageRepo  *repo.PageRepo
	blockRepo *repo.BlockRepo
}

func NewPageService(pageRepo *repo.PageRepo, blockRepo *repo.BlockRepo) *PageService {
	return &PageService{pageRepo: pageRepo, blockRepo: blockRepo}
}

func (s *PageService) Create(ctx context.Context, userID int64, title string, themePresetID int64) (*model.BioPage, error) {
	return s.pageRepo.Create(ctx, userID, themePresetID, title)
}

func (s *PageService) Get(ctx context.Context, pageID int64) (*model.BioPage, error) {
	return s.pageRepo.GetByID(ctx, pageID)
}

func (s *PageService) List(ctx context.Context, userID int64) ([]*model.BioPage, error) {
	return s.pageRepo.ListByUser(ctx, userID)
}

func (s *PageService) Delete(ctx context.Context, pageID int64) error {
	return s.pageRepo.Delete(ctx, pageID)
}

func (s *PageService) Update(ctx context.Context, page *model.BioPage) error {
	return s.pageRepo.Update(ctx, page)
}

// Draft data structure
type DraftData struct {
	Page       *model.BioPage    `json:"page"`
	Blocks     []*model.Block    `json:"blocks"`
	LinkGroups []*model.LinkGroup `json:"link_groups"`
	Links      map[int64][]*model.Link `json:"links"` // groupID -> links
}

func (s *PageService) GetDraft(ctx context.Context, pageID int64) (*DraftData, error) {
	page, err := s.pageRepo.GetByID(ctx, pageID)
	if err != nil {
		return nil, err
	}

	blocks, err := s.blockRepo.GetBlocksByPage(ctx, pageID)
	if err != nil {
		return nil, err
	}

	groups, err := s.blockRepo.GetLinkGroupsByPage(ctx, pageID)
	if err != nil {
		return nil, err
	}

	links := make(map[int64][]*model.Link)
	for _, g := range groups {
		groupLinks, err := s.blockRepo.GetLinksByGroup(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		links[g.ID] = groupLinks
	}

	return &DraftData{
		Page:       page,
		Blocks:     blocks,
		LinkGroups: groups,
		Links:      links,
	}, nil
}

// Save request structures
type SaveBlockReq struct {
	ID        *int64          `json:"id"`
	Type      string          `json:"type"`
	SortKey   string          `json:"sort_key"`
	RefID     *int64          `json:"ref_id"`
	Content   json.RawMessage `json:"content"`
	IsVisible bool            `json:"is_visible"`
	Delete    bool            `json:"delete"`
}

type SaveLinkGroupReq struct {
	ID            *int64          `json:"id"`
	Title         *string         `json:"title"`
	LayoutType    string          `json:"layout_type"`
	LayoutConfig  json.RawMessage `json:"layout_config"`
	StyleOverride json.RawMessage `json:"style_override"`
	Delete        bool            `json:"delete"`
}

type SaveLinkReq struct {
	ID       *int64 `json:"id"`
	GroupID  int64  `json:"group_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	SortKey  string `json:"sort_key"`
	IsActive bool   `json:"is_active"`
	Delete   bool   `json:"delete"`
}

type SaveRequest struct {
	Page       *model.BioPage     `json:"page"`
	Blocks     []SaveBlockReq     `json:"blocks"`
	LinkGroups []SaveLinkGroupReq `json:"link_groups"`
	Links      []SaveLinkReq      `json:"links"`
}

func (s *PageService) Save(ctx context.Context, pageID int64, req *SaveRequest) error {
	// Update page - merge with existing data
	if req.Page != nil {
		existing, err := s.pageRepo.GetByID(ctx, pageID)
		if err != nil {
			return err
		}

		// Merge fields - only update non-zero values
		if req.Page.Locale != "" {
			existing.Locale = req.Page.Locale
		}
		if req.Page.Title != nil {
			existing.Title = req.Page.Title
		}
		if req.Page.Status != "" {
			existing.Status = req.Page.Status
		}
		if req.Page.AccessType != "" {
			existing.AccessType = req.Page.AccessType
		}
		if req.Page.ThemePresetID != 0 {
			existing.ThemePresetID = req.Page.ThemePresetID
		}
		// ALWAYS update ThemeCustomID (even if nil) to allow clearing custom theme
		existing.ThemeCustomID = req.Page.ThemeCustomID
		
		if req.Page.ThemeMode != "" {
			existing.ThemeMode = req.Page.ThemeMode
		}
		if req.Page.Settings != nil {
			existing.Settings = req.Page.Settings
		}

		if err := s.pageRepo.Update(ctx, existing); err != nil {
			return err
		}
	}

	// Process link groups first (blocks may reference them)
	groupIDMap := make(map[int64]int64) // temp ID -> real ID
	for _, g := range req.LinkGroups {
		if g.Delete && g.ID != nil {
			if err := s.blockRepo.DeleteLinkGroup(ctx, *g.ID); err != nil {
				return err
			}
			continue
		}

		if g.ID != nil && *g.ID > 0 {
			// Update existing
			group := &model.LinkGroup{
				ID:            *g.ID,
				Title:         g.Title,
				LayoutType:    g.LayoutType,
				LayoutConfig:  g.LayoutConfig,
				StyleOverride: g.StyleOverride,
			}
			if err := s.blockRepo.UpdateLinkGroup(ctx, group); err != nil {
				return err
			}
		} else {
			// Create new
			group, err := s.blockRepo.CreateLinkGroup(ctx, pageID, g.Title, g.LayoutType)
			if err != nil {
				return err
			}
			if g.ID != nil {
				groupIDMap[*g.ID] = group.ID
			}
		}
	}

	// Process blocks
	for _, b := range req.Blocks {
		if b.Delete && b.ID != nil {
			if err := s.blockRepo.DeleteBlock(ctx, *b.ID); err != nil {
				return err
			}
			continue
		}

		// Map temp group ID to real ID
		refID := b.RefID
		if refID != nil && *refID < 0 {
			if realID, ok := groupIDMap[*refID]; ok {
				refID = &realID
			}
		}

		if b.ID != nil && *b.ID > 0 {
			// Update existing
			block := &model.Block{
				ID:        *b.ID,
				SortKey:   b.SortKey,
				Content:   b.Content,
				IsVisible: b.IsVisible,
			}
			if err := s.blockRepo.UpdateBlock(ctx, block); err != nil {
				return err
			}
		} else {
			// Create new
			_, err := s.blockRepo.CreateBlock(ctx, pageID, b.Type, b.SortKey, refID, b.Content)
			if err != nil {
				return err
			}
		}
	}

	// Process links
	for _, l := range req.Links {
		if l.Delete && l.ID != nil {
			if err := s.blockRepo.DeleteLink(ctx, *l.ID); err != nil {
				return err
			}
			continue
		}

		// Map temp group ID
		groupID := l.GroupID
		if groupID < 0 {
			if realID, ok := groupIDMap[groupID]; ok {
				groupID = realID
			}
		}

		sortKey := l.SortKey
		if sortKey == "" {
			sortKey = util.GenerateSortKey("", "")
		}

		if l.ID != nil && *l.ID > 0 {
			// Update existing
			link := &model.Link{
				ID:       *l.ID,
				Title:    l.Title,
				URL:      l.URL,
				SortKey:  sortKey,
				IsActive: l.IsActive,
			}
			if err := s.blockRepo.UpdateLink(ctx, link); err != nil {
				return err
			}
		} else {
			// Create new
			_, err := s.blockRepo.CreateLink(ctx, groupID, l.Title, l.URL, sortKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
