package routes

import (
	"github.com/Yash294/TODO/app/controllers/api"
	"github.com/Yash294/TODO/util"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(route fiber.Router) {
	route.Get("/login", util.RequireSession, controllers.RenderLogin)
	route.Get("/signup", util.RequireSession, controllers.RenderSignup)
	route.Get("/logout", util.RequireSession, controllers.Logout)
	route.Post("/login", util.RequireSession, controllers.Login)
	route.Post("/signup", util.RequireSession, controllers.Signup)
	route.Post("/resetPassword", util.RequireSession, controllers.ResetPassword)
}
