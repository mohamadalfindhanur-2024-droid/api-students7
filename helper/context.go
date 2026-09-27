package helper

import (
	"github.com/gofiber/fiber/v2"
	"api-students/app/model"
)

// LocalsAuthUser adalah kunci penyimpanan identitas pemakai di dalam context request Locals.
const LocalsAuthUser = "authUser"

// CurrentUser membaca identitas pemakai yang terotentikasi dari Locals context.
func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	user, ok := c.Locals(LocalsAuthUser).(model.AuthUser)
	return user, ok
}

// RequestID mengambil request ID unik dari context Locals request.
func RequestID(c *fiber.Ctx) string {
	if reqID, ok := c.Locals("requestid").(string); ok && reqID != "" {
		return reqID
	}
	return ""
}
