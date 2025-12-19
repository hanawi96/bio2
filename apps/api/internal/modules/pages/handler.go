package pages

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// POST /api/pages - Create new page
func (h *Handler) CreatePage(c *fiber.Ctx) error {
	var req CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid request body"},
		})
	}

	// MVP: hardcode user_id = 1 (no auth in slice 2)
	userID := int64(1)

	pageID, err := h.service.CreatePage(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "CREATE_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(fiber.Map{"page_id": pageID})
}

// GET /api/pages - List pages
func (h *Handler) ListPages(c *fiber.Ctx) error {
	// MVP: hardcode user_id = 1
	userID := int64(1)

	pages, err := h.service.ListPages(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "LIST_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(fiber.Map{"pages": pages})
}

// GET /api/pages/:id/draft - Get draft state
func (h *Handler) GetDraft(c *fiber.Ctx) error {
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_ID", Message: "Invalid page ID"},
		})
	}

	draft, err := h.service.GetDraft(c.Context(), pageID)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: ErrorDetail{Code: "NOT_FOUND", Message: "Page not found"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "GET_DRAFT_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(draft)
}

// POST /api/pages/:id/save - Save draft state
func (h *Handler) SaveDraft(c *fiber.Ctx) error {
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_ID", Message: "Invalid page ID"},
		})
	}

	var draft DraftPayload
	if err := c.BodyParser(&draft); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid draft payload"},
		})
	}

	updatedAt, err := h.service.SaveDraft(c.Context(), pageID, &draft)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: ErrorDetail{Code: "NOT_FOUND", Message: "Page not found"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "SAVE_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(fiber.Map{
		"ok":         true,
		"updated_at": updatedAt.Format(time.RFC3339),
	})
}

// POST /api/pages/:id/publish - Compile and publish
func (h *Handler) Publish(c *fiber.Ctx) error {
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_ID", Message: "Invalid page ID"},
		})
	}

	publishedAt, err := h.service.Publish(c.Context(), pageID)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: ErrorDetail{Code: "NOT_FOUND", Message: "Page not found"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "PUBLISH_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(fiber.Map{
		"ok":           true,
		"published_at": publishedAt.Format(time.RFC3339),
	})
}

// GET /api/pages/:id/published - Get published compiled_json
func (h *Handler) GetPublished(c *fiber.Ctx) error {
	pageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_ID", Message: "Invalid page ID"},
		})
	}

	compiled, err := h.service.GetPublished(c.Context(), pageID)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: ErrorDetail{Code: "NOT_FOUND", Message: "Page not published yet"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "GET_PUBLISHED_FAILED", Message: err.Error()},
		})
	}

	return c.JSON(compiled)
}
