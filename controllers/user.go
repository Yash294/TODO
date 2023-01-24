package controllers

import (
	//"fmt"
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

	err := repos.Login(result.Username, result.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User logged in successfully.",
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
		Username    string `json:"username"`
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

	err := repos.ChangePassword(result.Username, result.Password, result.NewPassword)

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

func CheckIfUsernameExists(c *fiber.Ctx) error {
	type UsernameCheck struct {
		Username string `json:"username"`
	}

	var result UsernameCheck

	if err := c.BodyParser(&result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Incorrect data format sent to server.",
		})
	}

	usernameAvailable, err := repos.IsUsernameAvailable(result.Username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Username check successful",
		"data":    usernameAvailable,
	})
}
