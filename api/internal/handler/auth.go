package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"linkbio/internal/middleware"
	"linkbio/internal/service"
	"linkbio/internal/util"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return util.BadRequest(c, "email and password required")
	}

	if len(req.Password) < 6 {
		return util.BadRequest(c, "password must be at least 6 characters")
	}

	user, token, err := h.authService.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrEmailExists {
			return util.BadRequest(c, "email already exists")
		}
		return util.InternalError(c)
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Set true in production
		SameSite: "Lax",
	})

	return util.Created(c, fiber.Map{
		"user":  user,
		"token": token,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	user, token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return util.Err(c, 401, "invalid email or password")
		}
		return util.InternalError(c)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return util.OK(c, fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return util.OK(c, nil)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return util.Unauthorized(c)
	}

	user, err := h.authService.GetUser(c.Context(), userID)
	if err != nil {
		return util.NotFound(c)
	}

	return util.OK(c, user)
}

type SetUsernameRequest struct {
	Username string `json:"username"`
}

func (h *AuthHandler) SetUsername(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return util.Unauthorized(c)
	}

	var req SetUsernameRequest
	if err := c.BodyParser(&req); err != nil {
		return util.BadRequest(c, "invalid request body")
	}

	user, err := h.authService.SetUsername(c.Context(), userID, req.Username)
	if err != nil {
		if err == service.ErrUsernameExists {
			return util.BadRequest(c, "username already taken")
		}
		if err == service.ErrInvalidUsername {
			return util.BadRequest(c, "username must be 3-30 characters, lowercase letters, numbers, underscore only")
		}
		return util.InternalError(c)
	}

	return util.OK(c, user)
}

func (h *AuthHandler) CheckUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	if username == "" {
		return util.BadRequest(c, "username required")
	}

	available, err := h.authService.CheckUsername(c.Context(), username)
	if err != nil {
		return util.InternalError(c)
	}

	return util.OK(c, fiber.Map{"available": available})
}
