package main

import (
	"github.com/Yash294/TODO/routes"
	"github.com/Yash294/TODO/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(util.NewMiddleware(), cors.New(cors.Config{
	 	AllowCredentials: true,
	 	AllowOrigins: "*",
	 	AllowHeaders: "Access-Control-Allow-Origin, Content-Type, Origin, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	setupRoutes(app)

	util.ConnectRedis()
	util.ConnectDB()

	app.Listen(":3000")
}
