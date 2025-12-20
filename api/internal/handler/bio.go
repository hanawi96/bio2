package handler

import (
	"github.com/gofiber/fiber/v2"
	"linkbio/internal/middleware"
	"linkbio/internal/service"
	"linkbio/internal/util"
)

type BioHandler struct {
	bioService *service.BioService
}

func NewBioHandler(bioService *service.BioService) *BioHandler {
	return &BioHandler{bioService: bioService}
}

// Get full bio data (blocks + groups + links)
func (h *BioHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	data, err := h.bioService.GetBio(c.Context(), userID)
	if err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, data)
}

// Add a new block
type AddBlockRequest struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

func (h *BioHandler) AddBlock(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req AddBlockRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	block, err := h.bioService.AddBlock(c.Context(), userID, req.Type, req.Content)
	if err != nil {
		return util.InternalError(c)
	}

	return util.Created(c, block)
}

// Update profile (display name and bio)
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func (h *BioHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	err := h.bioService.UpdateProfile(c.Context(), userID, req.DisplayName, req.Bio)
	if err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"message": "Profile updated successfully"})
}

// Update social links
type UpdateSocialLinksRequest struct {
	Instagram string `json:"instagram"`
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	TikTok    string `json:"tiktok"`
	YouTube   string `json:"youtube"`
	LinkedIn  string `json:"linkedin"`
	GitHub    string `json:"github"`
	Website   string `json:"website"`
}

func (h *BioHandler) UpdateSocialLinks(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req UpdateSocialLinksRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	err := h.bioService.UpdateSocialLinks(c.Context(), userID, req)
	if err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"message": "Social links updated successfully"})
}

// Update block
type UpdateBlockRequest struct {
	Content   any   `json:"content"`
	IsVisible *bool `json:"is_visible"`
}

func (h *BioHandler) UpdateBlock(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	blockID, err := parseID(c, "id")
	if err != nil {
		return util.BadRequest(c, "invalid block id")
	}

	var req UpdateBlockRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	block, err := h.bioService.UpdateBlock(c.Context(), userID, blockID, req.Content, req.IsVisible)
	if err != nil {
		if err == service.ErrNotFound {
			return util.NotFound(c)
		}
		if err == service.ErrForbidden {
			return util.Forbidden(c)
		}
		return util.InternalError(c)
	}

	return util.OK(c, block)
}

// Delete block
func (h *BioHandler) DeleteBlock(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	blockID, err := parseID(c, "id")
	if err != nil {
		return util.BadRequest(c, "invalid block id")
	}

	err = h.bioService.DeleteBlock(c.Context(), userID, blockID)
	if err != nil {
		if err == service.ErrNotFound {
			return util.NotFound(c)
		}
		if err == service.ErrForbidden {
			return util.Forbidden(c)
		}
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"deleted": true})
}

// Add link to a group
type AddLinkRequest struct {
	GroupID int64  `json:"group_id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}

func (h *BioHandler) AddLink(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req AddLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	if req.Title == "" || req.URL == "" {
		return util.BadRequest(c, "title and url required")
	}

	link, err := h.bioService.AddLink(c.Context(), userID, req.GroupID, req.Title, req.URL)
	if err != nil {
		if err == service.ErrForbidden {
			return util.Forbidden(c)
		}
		return util.InternalError(c)
	}

	return util.Created(c, link)
}

// Update link
type UpdateLinkRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	IsActive *bool  `json:"is_active"`
}

func (h *BioHandler) UpdateLink(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	linkID, err := parseID(c, "id")
	if err != nil {
		return util.BadRequest(c, "invalid link id")
	}

	var req UpdateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	link, err := h.bioService.UpdateLink(c.Context(), userID, linkID, req.Title, req.URL, req.IsActive)
	if err != nil {
		if err == service.ErrNotFound {
			return util.NotFound(c)
		}
		if err == service.ErrForbidden {
			return util.Forbidden(c)
		}
		return util.InternalError(c)
	}

	return util.OK(c, link)
}

// Delete link
func (h *BioHandler) DeleteLink(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	linkID, err := parseID(c, "id")
	if err != nil {
		return util.BadRequest(c, "invalid link id")
	}

	err = h.bioService.DeleteLink(c.Context(), userID, linkID)
	if err != nil {
		if err == service.ErrNotFound {
			return util.NotFound(c)
		}
		if err == service.ErrForbidden {
			return util.Forbidden(c)
		}
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"deleted": true})
}

// Reorder blocks
type ReorderBlocksRequest struct {
	BlockIDs []int64 `json:"block_ids"`
}

func (h *BioHandler) ReorderBlocks(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req ReorderBlocksRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	err := h.bioService.ReorderBlocks(c.Context(), userID, req.BlockIDs)
	if err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"reordered": true})
}

// Helper
func parseID(c *fiber.Ctx, param string) (int64, error) {
	id, err := c.ParamsInt(param)
	return int64(id), err
}
