package service

import (
	"context"
	"encoding/json"
	"errors"

	"linkbio/internal/model"
	"linkbio/internal/repo"
	"linkbio/internal/util"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type BioService struct {
	bioRepo   *repo.BioRepo
	pageRepo  *repo.PageRepo
	blockRepo *repo.BlockRepo
}

func NewBioService(bioRepo *repo.BioRepo, pageRepo *repo.PageRepo, blockRepo *repo.BlockRepo) *BioService {
	return &BioService{
		bioRepo:   bioRepo,
		pageRepo:  pageRepo,
		blockRepo: blockRepo,
	}
}

// BioData represents the full bio structure
type BioData struct {
	Page   *model.BioPage           `json:"page"`
	Blocks []*BlockWithGroup        `json:"blocks"`
}

type BlockWithGroup struct {
	*model.Block
	Group *GroupWithLinks `json:"group,omitempty"`
}

type GroupWithLinks struct {
	*model.LinkGroup
	Links []*model.Link `json:"links"`
}

func (s *BioService) GetBio(ctx context.Context, userID int64) (*BioData, error) {
	// Get or create page
	page, err := s.bioRepo.GetOrCreatePage(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get blocks
	blocks, err := s.blockRepo.GetBlocksByPage(ctx, page.ID)
	if err != nil {
		return nil, err
	}

	// Get groups
	groups, err := s.blockRepo.GetLinkGroupsByPage(ctx, page.ID)
	if err != nil {
		return nil, err
	}

	// Build group map with links
	groupMap := make(map[int64]*GroupWithLinks)
	for _, g := range groups {
		links, err := s.blockRepo.GetLinksByGroup(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		if links == nil {
			links = []*model.Link{}
		}
		groupMap[g.ID] = &GroupWithLinks{
			LinkGroup: g,
			Links:     links,
		}
	}

	// Build blocks with groups
	blocksWithGroups := make([]*BlockWithGroup, 0, len(blocks))
	for _, b := range blocks {
		bwg := &BlockWithGroup{Block: b}
		if b.Type == "link_group" && b.RefID != nil {
			if g, ok := groupMap[*b.RefID]; ok {
				bwg.Group = g
			}
		}
		blocksWithGroups = append(blocksWithGroups, bwg)
	}

	return &BioData{
		Page:   page,
		Blocks: blocksWithGroups,
	}, nil
}

func (s *BioService) AddBlock(ctx context.Context, userID int64, blockType string, content any) (*BlockWithGroup, error) {
	page, err := s.bioRepo.GetOrCreatePage(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get last sort key
	lastKey, _ := s.bioRepo.GetLastBlockSortKey(ctx, page.ID)
	sortKey := util.GenerateSortKey(lastKey, "")

	var refID *int64
	var contentJSON json.RawMessage

	if blockType == "link_group" {
		// Create link group first
		group, err := s.blockRepo.CreateLinkGroup(ctx, page.ID, nil, "list")
		if err != nil {
			return nil, err
		}
		refID = &group.ID
		contentJSON = json.RawMessage(`{}`)
	} else {
		// Serialize content
		if content != nil {
			contentJSON, err = json.Marshal(content)
			if err != nil {
				return nil, err
			}
		} else {
			contentJSON = json.RawMessage(`{}`)
		}
	}

	block, err := s.blockRepo.CreateBlock(ctx, page.ID, blockType, sortKey, refID, contentJSON)
	if err != nil {
		return nil, err
	}

	result := &BlockWithGroup{Block: block}

	if blockType == "link_group" && refID != nil {
		groups, _ := s.blockRepo.GetLinkGroupsByPage(ctx, page.ID)
		for _, g := range groups {
			if g.ID == *refID {
				result.Group = &GroupWithLinks{
					LinkGroup: g,
					Links:     []*model.Link{},
				}
				break
			}
		}
	}

	return result, nil
}

func (s *BioService) UpdateBlock(ctx context.Context, userID, blockID int64, content any, isVisible *bool) (*model.Block, error) {
	// Check ownership
	ownerID, err := s.bioRepo.GetBlockOwnerID(ctx, blockID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ownerID != userID {
		return nil, ErrForbidden
	}

	block, err := s.bioRepo.GetBlockByID(ctx, blockID)
	if err != nil {
		return nil, ErrNotFound
	}

	if content != nil {
		contentJSON, err := json.Marshal(content)
		if err != nil {
			return nil, err
		}
		block.Content = contentJSON
	}

	if isVisible != nil {
		block.IsVisible = *isVisible
	}

	err = s.blockRepo.UpdateBlock(ctx, block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s *BioService) DeleteBlock(ctx context.Context, userID, blockID int64) error {
	ownerID, err := s.bioRepo.GetBlockOwnerID(ctx, blockID)
	if err != nil {
		return ErrNotFound
	}
	if ownerID != userID {
		return ErrForbidden
	}

	return s.blockRepo.DeleteBlock(ctx, blockID)
}

func (s *BioService) AddLink(ctx context.Context, userID, groupID int64, title, url string) (*model.Link, error) {
	// Check ownership
	ownerID, err := s.bioRepo.GetGroupOwnerID(ctx, groupID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ownerID != userID {
		return nil, ErrForbidden
	}

	// Get last sort key
	lastKey, _ := s.bioRepo.GetLastLinkSortKey(ctx, groupID)
	sortKey := util.GenerateSortKey(lastKey, "")

	return s.blockRepo.CreateLink(ctx, groupID, title, url, sortKey)
}

func (s *BioService) UpdateLink(ctx context.Context, userID, linkID int64, title, url string, isActive *bool) (*model.Link, error) {
	ownerID, err := s.bioRepo.GetLinkOwnerID(ctx, linkID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ownerID != userID {
		return nil, ErrForbidden
	}

	link, err := s.bioRepo.GetLinkByID(ctx, linkID)
	if err != nil {
		return nil, ErrNotFound
	}

	if title != "" {
		link.Title = title
	}
	if url != "" {
		link.URL = url
	}
	if isActive != nil {
		link.IsActive = *isActive
	}

	err = s.blockRepo.UpdateLink(ctx, link)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (s *BioService) DeleteLink(ctx context.Context, userID, linkID int64) error {
	ownerID, err := s.bioRepo.GetLinkOwnerID(ctx, linkID)
	if err != nil {
		return ErrNotFound
	}
	if ownerID != userID {
		return ErrForbidden
	}

	return s.blockRepo.DeleteLink(ctx, linkID)
}

func (s *BioService) ReorderBlocks(ctx context.Context, userID int64, blockIDs []int64) error {
	for i, blockID := range blockIDs {
		// Generate sort key based on position
		sortKey := string(rune('A' + i))
		if i >= 26 {
			sortKey = string(rune('A'+i/26-1)) + string(rune('A'+i%26))
		}

		if err := s.bioRepo.UpdateBlockSortKey(ctx, blockID, sortKey); err != nil {
			return err
		}
	}
	return nil
}
