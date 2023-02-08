package routes

import (
	"github.com/Yash294/TODO/controllers"
	"github.com/Yash294/TODO/util"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(route fiber.Router) {
	route.Get("/dashboard", util.RequireSession, controllers.RenderTasks)
	route.Post("/add", util.RequireSession, controllers.AddTask)
	route.Post("/edit", util.RequireSession, controllers.EditTask)
	route.Post("/delete", util.RequireSession, controllers.DeleteTask)
}