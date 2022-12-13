package main

import (
	"collaborators-tracking-platforms/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	routes.CollaboratorsRoute(app)

	app.Listen(":8080")
}
