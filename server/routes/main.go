package routes

import (
	"github.com/dishan1223/mutt/internal/middleware"
	"github.com/dishan1223/mutt/server/handler"
	"github.com/gofiber/fiber/v3"
)

var app *fiber.App

func Init(a *fiber.App) {
	app = a

	v1 := app.Group("/api/v1")

	app.Get("/ping", handler.Ping)

	v1.Get("/ping", handler.Ping)

	auth := v1.Group("/auth")
	auth.Post("/signup", handler.SignUpHandler)
	auth.Post("/login", handler.LoginHandler)
	auth.Post("/refresh", handler.RefreshTokenHandler)
	auth.Post("/logout", middleware.AuthRequired, handler.LogoutHandler)
	auth.Get("/me", middleware.AuthRequired, handler.MeHandler)
}
