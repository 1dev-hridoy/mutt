package routes

import (
	"github.com/dishan1223/mutt/server/handler"
	"github.com/gofiber/fiber/v3"
)

var app *fiber.App

func Init(a *fiber.App) {
	app = a

	v1 := app.Group("/api/v1")

	app.Get("/test", handler.Test)

	v1.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
