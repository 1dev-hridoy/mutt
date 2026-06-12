package handler

import "github.com/gofiber/fiber/v3"

func Test(c fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
