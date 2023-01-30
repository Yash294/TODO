package routes

import (
	"github.com/Yash294/TODO/controllers"
	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(route fiber.Router) {
	route.Get("/dashboard", controllers.RenderTasks)
	route.Post("/add", controllers.AddTask)
	route.Delete("/delete", controllers.DeleteTask)
	route.Patch("/edit", controllers.EditTask)
}