package controllers

import (
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
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var result LoginInfo

	if err := c.BodyParser(&result); err != nil {
		return err
	}

	err := repos.Login(result.Email, result.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	// sess, sessErr := store.Get(c)
	// if sessErr != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": sessErr.Error,
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User logged in successfully.",
	})
}

func Signup(c *fiber.Ctx) error {
	type NewSignup struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var result NewSignup

	if err := c.BodyParser(&result); err != nil {
		fmt.Println(result.Email, result.Password)
		return err
	}

	err := repos.CreateUser(result.Email, result.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User signed up successfully.",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	type UserInfo struct {
		Email    string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}

	var result UserInfo

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Incorrect data format sent to server.",
		})
	}

	err := repos.ChangePassword(result.Email, result.Password, result.NewPassword)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Password changed successfully.",
	})
}
