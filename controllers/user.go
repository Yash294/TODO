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
	type LoginInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var result LoginInfo

	if err := c.BodyParser(&result); err != nil {
		return err
	}

	success, err := repos.Login(result.Username, result.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to check user login credentials",
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": success,
	})
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

	err := repos.CreateUser(result.Username, result.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create a new user."
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func CheckIfUsernameExists(c *fiber.Ctx) error {
	type Username struct {
		username string
	}

	var result Username

	if err := c.BodyParser(&result); err != nil {
		return err
	}

	exists, err := repos.CheckIfUsernameExists(result.username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to check username existence.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": exists,
	})
}
