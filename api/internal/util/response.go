package util

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{Success: true, Data: data})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(201).JSON(Response{Success: true, Data: data})
}

func Err(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(Response{Success: false, Error: msg})
}

func BadRequest(c *fiber.Ctx, msg string) error {
	return Err(c, 400, msg)
}

func Unauthorized(c *fiber.Ctx) error {
	return Err(c, 401, "unauthorized")
}

func Forbidden(c *fiber.Ctx) error {
	return Err(c, 403, "forbidden")
}

func NotFound(c *fiber.Ctx) error {
	return Err(c, 404, "not found")
}

func InternalError(c *fiber.Ctx) error {
	return Err(c, 500, "internal server error")
}
