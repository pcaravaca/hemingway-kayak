package middleware

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// RequireResourceOwnership verifica en la base que el recurso (:id) pertenezca
// al usuario autenticado. table y ownerColumn son constantes internas, nunca
// entrada del usuario, por lo que no hay inyección.
//
// Ej.: RequireResourceOwnership(db, "trips", "user_id", "id")
func RequireResourceOwnership(db *sql.DB, table, ownerColumn, param string) fiber.Handler {
	query := "SELECT CAST(" + ownerColumn + " AS VARCHAR(36)) FROM " + table + " WHERE id = $1"

	return func(c *fiber.Ctx) error {
		id := c.Params(param)
		if id == "" {
			return c.Status(400).JSON(fiber.Map{"success": false, "error": "missing resource id"})
		}

		var owner string
		err := db.QueryRowContext(c.Context(), query, id).Scan(&owner)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"success": false, "error": "not found"})
		}
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": "database error"})
		}
		if owner != UserID(c) {
			return c.Status(403).JSON(fiber.Map{"success": false, "error": "forbidden"})
		}
		return c.Next()
	}
}
