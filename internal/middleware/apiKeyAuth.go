package middleware

import (
	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/internal/service"
	"github.com/dishan1223/mutt/models"
	"github.com/gofiber/fiber/v3"
)

func APIKeyAuth(c fiber.Ctx) error {
	apiKey := c.Get("X-Mutt-Key")
	if apiKey == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "API key required. Pass it in X-Mutt-Key header.",
		})
	}

	hashedKey := service.HashAPIKey(apiKey)

	var project models.Project
	if err := config.DB.Where("api_key = ?", hashedKey).First(&project).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid API key",
		})
	}

	c.Locals("projectID", project.ID)
	return c.Next()
}
