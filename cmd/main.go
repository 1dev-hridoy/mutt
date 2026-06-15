package main

import (
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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false, // Must be false when AllowOrigins is "*"
	}))

	routes.Init(app)
	app.Listen(PORT)
}
