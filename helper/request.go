package helper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"api-students/app/model"

	"github.com/gofiber/fiber/v2"
)

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func ParamID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func ParseCursorQuery(c *fiber.Ctx) (model.CursorQuery, error) {
	q := model.CursorQuery{
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
	}

	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}

	if rawIsActive := strings.TrimSpace(c.Query("is_active")); rawIsActive != "" {
		v, err := strconv.ParseBool(rawIsActive)
		if err != nil {
			return q, BadRequest("parameter is_active harus berupa boolean")
		}
		q.IsActive = &v
	}

	if rawCursor := strings.TrimSpace(c.Query("cursor")); rawCursor != "" {
		cursor, err := DecodeCursor(rawCursor)
		if err != nil {
			return q, BadRequest("cursor tidak valid")
		}
		q.After = &cursor
	}

	return q, nil
}
