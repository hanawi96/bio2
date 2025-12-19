package public

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse follows api_contract.md error format
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

// GET /r - Returns compiled_json from page_publish_cache by Host
// Custom domain: 1 domain = 1 page, always resolve to path "/"
func (h *Handler) GetPage(c *fiber.Ctx) error {
	// Get host from query param (dev helper) or header
	host := c.Query("host")
	if host == "" {
		host = c.Get("Host")
	}

	if host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{
				Code:    "MISSING_HOST",
				Message: "Host header or query param required",
			},
		})
	}

	// Get page result (includes access_type check and compiled data)
	result, err := h.service.GetCompiledPage(c.Context(), host)
	if err != nil {
		if err == ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "Page not found",
				},
			})
		}
		if err == ErrPasswordRequired {
			// Password-protected pages: no caching for security
			c.Set("Cache-Control", "private, no-store")
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Error: ErrorDetail{
					Code:    "PASSWORD_REQUIRED",
					Message: "This page is password protected",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Set cache headers based on access type
	if result.IsPasswordProtected {
		// Even when authorized via cookie, password pages must not be cached publicly
		c.Set("Cache-Control", "private, no-store")
	} else {
		// Public pages: CDN cacheable per api_contract.md section 0.8
		c.Set("Cache-Control", "public, max-age=60, s-maxage=3600, stale-while-revalidate=86400")
	}

	// Use stable ETag from publish cache (computed at publish time)
	c.Set("ETag", `"`+result.ETag+`"`)

	return c.JSON(result.Compiled)
}

// POST /r/password - Verify password (placeholder for slice 1)
func (h *Handler) VerifyPassword(c *fiber.Ctx) error {
	// Password endpoints must never be cached
	c.Set("Cache-Control", "private, no-store")
	return c.Status(fiber.StatusNotImplemented).JSON(ErrorResponse{
		Error: ErrorDetail{
			Code:    "NOT_IMPLEMENTED",
			Message: "Password verification not implemented yet",
		},
	})
}
