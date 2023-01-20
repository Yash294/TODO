package main

import (
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func setupRoutes(app *fiber.App) {
	userGroup := app.Group("/user")
	taskGroup := app.Group("/task")

	routes.UserRoutes(userGroup)
	routes.TaskRoutes(taskGroup)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/user/signup")
	})

	// app.Get("/tasks/:userId", func(c *fiber.Ctx) error {
	// 	return c.Render("dashboard", fiber.Map{})
	// })
}

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	database.Connect()

	setupRoutes(app)

	app.Listen(":3000")
}
