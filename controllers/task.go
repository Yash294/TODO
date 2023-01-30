package controllers

import (
	//"fmt"
	"github.com/Yash294/TODO/models"
	"github.com/Yash294/TODO/repos"
	"github.com/gofiber/fiber/v2"
)

func RenderTasks(c *fiber.Ctx) error {
	var data models.Task

	data.Assignee = "yashp@gmail.com"

	result, err := repos.GetTasks(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Render("task", fiber.Map{
		"data": result,
	})
}

func AddTask(c *fiber.Ctx) error {
	var data models.Task

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	data.Assignee = "yashp@gmail.com"

	if err := repos.AddTask(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task created successfully.",
	})
}

func DeleteTask(c *fiber.Ctx) error {
	var data models.Task

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	data.Assignee = "yashp@gmail.com"

	if err := repos.DeleteTask(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task deleted successfully.",
	})
}

func EditTask(c *fiber.Ctx) error {
	var data models.Task

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if err := repos.EditTask(&data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Task edited successfully.",
	})
}
