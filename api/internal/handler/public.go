package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"linkbio/internal/repo"
	"linkbio/internal/util"
)

type PublicHandler struct {
	pageRepo   *repo.PageRepo
	domainRepo *repo.DomainRepo
}

func NewPublicHandler(pageRepo *repo.PageRepo, domainRepo *repo.DomainRepo) *PublicHandler {
	return &PublicHandler{
		pageRepo:   pageRepo,
		domainRepo: domainRepo,
	}
}

func (h *PublicHandler) Render(c *fiber.Ctx) error {
	// Get hostname from Host header
	host := c.Get("Host")
	if host == "" {
		return util.BadRequest(c, "missing host header")
	}

	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Get path from query or default to /
	path := c.Query("path", "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Find domain
	domain, err := h.domainRepo.GetByHostname(c.Context(), host)
	if err != nil {
		// Try system domain
		domain, err = h.domainRepo.GetSystemDomain(c.Context())
		if err != nil {
			return util.NotFound(c)
		}
		// For system domain, path comes from query
	}

	// Find route
	route, err := h.domainRepo.GetRouteByDomainAndPath(c.Context(), domain.ID, path)
	if err != nil {
		return util.NotFound(c)
	}

	// Check if redirect
	if route.RedirectToRouteID != nil {
		// TODO: implement redirect
		return util.NotFound(c)
	}

	// Get page
	page, err := h.pageRepo.GetByID(c.Context(), route.PageID)
	if err != nil {
		return util.NotFound(c)
	}

	// Check password protection
	if page.AccessType == "password" {
		// TODO: check session cookie
		return util.Err(c, 401, "password required")
	}

	// Get cached compiled JSON
	cache, err := h.pageRepo.GetPublishCache(c.Context(), page.ID)
	if err != nil {
		return util.NotFound(c)
	}

	// Set cache headers
	c.Set("Cache-Control", "public, max-age=60, s-maxage=300, stale-while-revalidate=86400")
	c.Set("Content-Type", "application/json")

	return c.Send(cache.CompiledJSON)
}

type VerifyPasswordRequest struct {
	PageID   int64  `json:"page_id"`
	Password string `json:"password"`
}

func (h *PublicHandler) VerifyPassword(c *fiber.Ctx) error {
	var req VerifyPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	page, err := h.pageRepo.GetByID(c.Context(), req.PageID)
	if err != nil {
		return util.NotFound(c)
	}

	if page.AccessType != "password" || page.PasswordHash == nil {
		return util.BadRequest(c, "page is not password protected")
	}

	if !util.CheckPassword(req.Password, *page.PasswordHash) {
		return util.Err(c, 401, "invalid password")
	}

	// TODO: create session and set cookie for 7 days

	// Return compiled JSON
	cache, err := h.pageRepo.GetPublishCache(c.Context(), page.ID)
	if err != nil {
		return util.NotFound(c)
	}

	return c.Send(cache.CompiledJSON)
}
