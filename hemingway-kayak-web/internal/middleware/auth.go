// Package middleware contiene los middlewares de autenticación y autorización.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"hemingway-kayak-web/pkg/auth"
)

const localsUserID = "userId"

// JWTAuth exige un Bearer token válido y no revocado, y deja la identidad en
// c.Locals("userId") / c.Locals("email").
func JWTAuth(secret []byte, revoked *auth.Revocation) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{"success": false, "error": "missing authorization header"})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"success": false, "error": "invalid authorization header"})
		}

		claims, err := auth.Parse(secret, tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"success": false, "error": "invalid or expired token"})
		}

		if revoked != nil && revoked.IsRevoked(tokenString) {
			return c.Status(401).JSON(fiber.Map{"success": false, "error": "token revoked"})
		}

		c.Locals(localsUserID, claims.UserID)
		c.Locals("email", claims.Email)
		return c.Next()
	}
}

// RequireOwnership garantiza que el parámetro de ruta pertenezca al usuario autenticado.
func RequireOwnership(param string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Params(param) != UserID(c) {
			return c.Status(403).JSON(fiber.Map{"success": false, "error": "forbidden"})
		}
		return c.Next()
	}
}

// UserID devuelve el id del usuario autenticado ("" si no hubo middleware).
func UserID(c *fiber.Ctx) string {
	id, _ := c.Locals(localsUserID).(string)
	return id
}
