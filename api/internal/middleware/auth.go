package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try cookie first
		token := c.Cookies("token")

		// Then try Authorization header
		if token == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims := jwt.MapClaims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !parsed.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("userID", int64(userID))
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) int64 {
	if id, ok := c.Locals("userID").(int64); ok {
		return id
	}
	return 0
}
