package controllers

import (
	"github.com/Yash294/TODO/app/models"
	"github.com/Yash294/TODO/app/repos"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/redis"
	"github.com/gofiber/fiber/v2"
)

func RenderLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func RenderSignup(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	var data models.UserDTO

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	result, err := repos.Login(&data, database.DB, repos.EncryptionInstance)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	sess, err := redis.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	sess.Set(redis.AUTH_KEY, true)
	sess.Set(redis.USER_ID, result)

	if err := sess.Save(); err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "log in successful",
	})
}

func Signup(c *fiber.Ctx) error {
	var data models.UserDTO

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := repos.CreateUser(&data, database.DB, repos.EncryptionInstance, repos.CopierInstance); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "sign up successful",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	var data models.UserDTO

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := repos.ChangePassword(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "password changed successfully",
	})
}

func Logout(c *fiber.Ctx) error {

	sess, err := redis.Store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Already logged out.",
		})
	}

	if err = sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout Successful.",
	})
}
