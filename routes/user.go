package routes

import (
	"github.com/Yash294/TODO/app/controllers/api"
	"github.com/Yash294/TODO/app/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(route fiber.Router) {
	route.Get("/login", middleware.RequireSession, controllers.RenderLogin)
	route.Get("/signup", middleware.RequireSession, controllers.RenderSignup)
	route.Get("/logout", middleware.RequireSession, controllers.Logout)
	route.Post("/login", middleware.RequireSession, controllers.Login)
	route.Post("/signup", middleware.RequireSession, controllers.Signup)
	route.Post("/resetPassword", middleware.RequireSession, controllers.ResetPassword)
}
