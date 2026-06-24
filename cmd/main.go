package main

import (
	"strings"

	"github.com/dishan1223/mutt/consts"
	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/internal/service"
	"github.com/dishan1223/mutt/server/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func init() {
	config.MustLoadEnv()
	config.MustConnectToDB()
	config.MustSyncDatabase()
	config.MustConnectRedis()
	service.MustInitJWT(config.MustGetEnv("JWT_SECRET"))
}

func main() {
	PORT := consts.PORT
	app := fiber.New()

	allowedOrigins := strings.Split(config.MustGetEnv("ALLOWED_ORIGINS"), ",")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Mutt-Key"},
		AllowCredentials: true,
	}))

	routes.Init(app)
	app.Listen(PORT)
}
