package main

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})
}
