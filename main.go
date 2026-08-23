package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var bodyMethods = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// requireJSON menolak request dengan body jika Content-Type bukan application/json (mengembalikan 415)
func requireJSON(c *fiber.Ctx) error {
	if bodyMethods[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Tugas Praktikum REST API Students - Pertemuan 2",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			message := "terjadi kesalahan pada server"

			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				message = e.Message
			}

			return fail(c, status, message)
		},
	})

	// Middleware Global
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())

	// Endpoint Root & Health Check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students Server is Running!")
	})

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return ok(c, "server berjalan normal", fiber.Map{"timestamp": time.Now()})
	})

	// Routing Students
	studentsRoute := api.Group("/students", requireJSON)
	studentsRoute.Get("/", listStudents)
	studentsRoute.Get("/:id", getStudent)
	studentsRoute.Post("/", createStudent)
	studentsRoute.Put("/:id", replaceStudent)
	studentsRoute.Patch("/:id", patchStudent)
	studentsRoute.Delete("/:id", deleteStudent)

	// Handler untuk route yang tidak ditemukan (404)
	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	fmt.Println("Server API Students berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}