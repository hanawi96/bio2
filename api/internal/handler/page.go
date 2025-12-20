package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"linkbio/internal/middleware"
	"linkbio/internal/model"
	"linkbio/internal/service"
	"linkbio/internal/util"
)

type PageHandler struct {
	pageService     *service.PageService
	compilerService *service.CompilerService
}

func NewPageHandler(pageService *service.PageService, compilerService *service.CompilerService) *PageHandler {
	return &PageHandler{
		pageService:     pageService,
		compilerService: compilerService,
	}
}

type CreatePageRequest struct {
	Title         string `json:"title"`
	ThemePresetID int64  `json:"theme_preset_id"`
}

func (h *PageHandler) Create(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	if req.ThemePresetID == 0 {
		req.ThemePresetID = 1 // Default theme
	}

	page, err := h.pageService.Create(c.Context(), userID, req.Title, req.ThemePresetID)
	if err != nil {
		return util.InternalError(c)
	}

	return util.Created(c, page)
}

func (h *PageHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	pages, err := h.pageService.List(c.Context(), userID)
	if err != nil {
		return util.InternalError(c)
	}

	if pages == nil {
		pages = []*model.BioPage{}
	}

	return util.OK(c, pages)
}

func (h *PageHandler) Get(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return util.BadRequest(c, "invalid page id")
	}

	page, err := h.pageService.Get(c.Context(), pageID)
	if err != nil {
		return util.NotFound(c)
	}

	if page.UserID != userID {
		return util.Forbidden(c)
	}

	return util.OK(c, page)
}

func (h *PageHandler) GetDraft(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return util.BadRequest(c, "invalid page id")
	}

	draft, err := h.pageService.GetDraft(c.Context(), pageID)
	if err != nil {
		return util.NotFound(c)
	}

	if draft.Page.UserID != userID {
		return util.Forbidden(c)
	}

	return util.OK(c, draft)
}

func (h *PageHandler) Save(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return util.BadRequest(c, "invalid page id")
	}

	// Check ownership
	page, err := h.pageService.Get(c.Context(), pageID)
	if err != nil {
		return util.NotFound(c)
	}
	if page.UserID != userID {
		return util.Forbidden(c)
	}

	var req service.SaveRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	if err := h.pageService.Save(c.Context(), pageID, &req); err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"saved": true})
}

func (h *PageHandler) Publish(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return util.BadRequest(c, "invalid page id")
	}

	// Check ownership
	page, err := h.pageService.Get(c.Context(), pageID)
	if err != nil {
		return util.NotFound(c)
	}
	if page.UserID != userID {
		return util.Forbidden(c)
	}

	if err := h.compilerService.Publish(c.Context(), pageID); err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"published": true})
}

func (h *PageHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return util.BadRequest(c, "invalid page id")
	}

	// Check ownership
	page, err := h.pageService.Get(c.Context(), pageID)
	if err != nil {
		return util.NotFound(c)
	}
	if page.UserID != userID {
		return util.Forbidden(c)
	}

	if err := h.pageService.Delete(c.Context(), pageID); err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"deleted": true})
}
