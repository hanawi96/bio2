package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"linkbio/internal/middleware"
	"linkbio/internal/model"
	"linkbio/internal/service"
	"linkbio/internal/util"
)

type ThemeHandler struct {
	themeService *service.ThemeService
}

func NewThemeHandler(themeService *service.ThemeService) *ThemeHandler {
	return &ThemeHandler{themeService: themeService}
}

func (h *ThemeHandler) ListPresets(c *fiber.Ctx) error {
	tier := c.Query("tier") // free, pro, or empty for all

	presets, err := h.themeService.ListPresets(c.Context(), tier)
	if err != nil {
		return util.InternalError(c)
	}

	if presets == nil {
		presets = []*model.ThemePreset{}
	}

	return util.OK(c, presets)
}

type CreateCustomRequest struct {
	PresetID int64           `json:"preset_id"`
	Patch    json.RawMessage `json:"patch"`
}

func (h *ThemeHandler) CreateCustom(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req CreateCustomRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	if req.PresetID == 0 {
		return util.BadRequest(c, "preset_id required")
	}

	custom, err := h.themeService.CreateOrUpdateCustom(c.Context(), userID, req.PresetID, req.Patch)
	if err != nil {
		return util.InternalError(c)
	}

	return util.Created(c, custom)
}

type ApplyThemeRequest struct {
	PageID        int64  `json:"page_id"`
	PresetID      int64  `json:"preset_id"`
	CustomID      *int64 `json:"custom_id"`
	Mode          string `json:"mode"`
}

func (h *ThemeHandler) Apply(c *fiber.Ctx) error {
	// This would update the page's theme settings
	// For now, just return success - actual implementation would use PageService
	return util.OK(c, fiber.Map{"applied": true})
}
