package routes

import (
	"github.com/Yash294/TODO/app/controllers/api"
	"github.com/Yash294/TODO/app/middleware"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(route fiber.Router) {
	route.Get("/dashboard", middleware.RequireSession, controllers.RenderTasks)
	route.Post("/add", middleware.RequireSession, controllers.AddTask)
	route.Post("/edit", middleware.RequireSession, controllers.EditTask)
	route.Post("/delete", middleware.RequireSession, controllers.DeleteTask)
}