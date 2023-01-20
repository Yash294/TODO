package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/Yash294/TODO/repos"
	"github.com/gofiber/fiber/v2"
)

func RenderLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func RenderSignup(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	return nil
}

func Signup(c *fiber.Ctx) error {
	type NewSignup struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var result NewSignup

	if err := c.BodyParser(&result); err != nil {
		return err
	}

	dbErr := repos.CreateUser(result.Username, result.Password)

	if dbErr != nil {
		return c.JSON(fiber.Map{"status": "fail", "message": "Signup FAIL", "data": result})
	}

	fmt.Println("DB SUCCESS")
	return c.JSON(fiber.Map{"status": "success", "message": "Signup SUCCESS", "data": result})
}

func GetAllUsernames(c *fiber.Ctx) error {

	usernames, dbErr := repos.GetAllUsernames()

	if dbErr != nil {
		return c.JSON(fiber.Map{"status": "fail", "message": "Get User FAIL"})
	}

	bytes, _ := json.Marshal(usernames)

	return c.SendString(string(bytes))
}
