package main

import (
	"github.com/dishan1223/mutt/consts"
	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/server/routes"
	"github.com/gofiber/fiber/v3"
)

func init() {
	config.MustLoadEnv()
	config.MustConnectToDB()
}

func main() {
	PORT := consts.PORT
	app := fiber.New()
	routes.Init(app)
	app.Listen(PORT)
}
