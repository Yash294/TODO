package controllers

import (
	"github.com/Yash294/TODO/app/models"
	"github.com/Yash294/TODO/app/repos"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/app/middleware"
	"github.com/gofiber/fiber/v2"
)

func RenderTasks(c *fiber.Ctx) error {
	userId, err := middleware.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	result, err := repos.GetTasks(userId, database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	email, err := repos.GetUser(userId, database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Render("task", fiber.Map{
		"username": email,
		"data":     result,
	})
}

func AddTask(c *fiber.Ctx) error {
	userId, err := middleware.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var data models.TaskDTO

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	taskID, err := repos.AddTask(&data, userId, database.DB, repos.CopierInstance)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "task created successfully",
		"taskID": taskID,
	})
}

func EditTask(c *fiber.Ctx) error {
	_, err := middleware.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var data models.TaskDTO

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err = repos.EditTask(&data, database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "task edited successfully",
	})
}

func DeleteTask(c *fiber.Ctx) error {
	_, err := middleware.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var data models.TaskDTO

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err = repos.DeleteTask(&data, database.DB); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "task deleted successfully",
	})
}
