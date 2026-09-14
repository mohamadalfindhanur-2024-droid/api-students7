package middleware

import (
	"github.com/gofiber/fiber/v2"
	"api-students/helper"
)

// RequirePermission menolak request jika role user tidak memiliki permission yang diminta.
// Digunakan pada route level yang keputusannya tidak bergantung pada kepemilikan isi data.
func RequirePermission(perms *helper.PermissionSet, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
		}
		if !perms.Can(user.Role, permission) {
			return helper.Fail(c, fiber.StatusForbidden,
				"role "+user.Role+" tidak memiliki hak "+permission)
		}
		return c.Next()
	}
}

// RequireRole memeriksa nama role secara langsung (sebagai pembanding).
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
		}
		if _, granted := allowed[user.Role]; !granted {
			return helper.Fail(c, fiber.StatusForbidden,
				"role Anda tidak berhak mengakses endpoint ini")
		}
		return c.Next()
	}
}
