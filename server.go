package main

import (
	"github.com/Yash294/TODO/routes"
	"github.com/Yash294/TODO/app/database"
	"github.com/Yash294/TODO/app/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func setupRoutes(app *fiber.App) {
	userGroup := app.Group("/user")
	taskGroup := app.Group("/task")

	routes.UserRoutes(userGroup)
	routes.TaskRoutes(taskGroup)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/user/login")
	})
}

func main() {

	engine := html.New("./resources/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./resources")

	setupRoutes(app)

	redis.ConnectRedis()
	database.ConnectDB()

	app.Listen(":3000")
}
