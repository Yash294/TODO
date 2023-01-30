package main

import (
	//"fmt"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	//"github.com/gofiber/fiber/v2/middleware/session"
	//"github.com/gofiber/fiber/v2/middleware/cors"
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

// var (
// 	store *session.Store
// 	AUTH_KEY string = "authenticated"
// 	USER_ID string = "user_id"
// )

// func AuthMiddleware(c *fiber.Ctx) error {
// 	sess, err := store.Get(c)

// 	fmt.Println(strings.Split(c.Path(), "/")[1])

// 	path := strings.Split(c.Path(), "/")[1] 

// 	if path == "login" || path == "signup" {
// 		return c.Next()
// 	}

// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"message": "Not authorized",
// 		})
// 	}

// 	if sess.Get(AUTH_KEY) == nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"message": "Not authorized",
// 		})
// 	}

// 	return c.Next()
// }

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// store := session.New(session.Config{
	// 	CookieHTTPOnly: true,
	// 	Expiration: time.Hour * 4,
	// })

	// app.Use(AuthMiddleware(), cors.New(cors.Config{
	// 	AllowCredentials: true,
	// 	AllowOrigins: "*",
	// 	AllowHeaders: "Access-Control-Allow-Origin, Content-Type, Origin,"
	// }))

	setupRoutes(app)

	database.Connect()

	app.Listen(":3000")
}
