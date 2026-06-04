package main

import (
	"os"

	"github.com/dishan1223/mutt/internal/config"
	"github.com/gofiber/fiber/v3"
)

func init() {
	config.LoadEnv()
	config.ConnectToDB()
}

func main() {
	// app.Listen(":PORT")
	PORT := ":" + os.Getenv("PORT")

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("hello from mutt")
	})

	app.Listen(PORT)
}
