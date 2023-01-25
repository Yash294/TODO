package routes

import (
	"github.com/Yash294/TODO/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(route fiber.Router) {
	route.Get("/login", controllers.RenderLogin)
	route.Get("/signup", controllers.RenderSignup)
	route.Post("/login", controllers.Login)
	route.Post("/signup", controllers.Signup)
	route.Post("/login/resetPassword", controllers.ResetPassword)
}
