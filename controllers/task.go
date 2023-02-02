package controllers

import (
	"fmt"
	"github.com/Yash294/TODO/models"
	"github.com/Yash294/TODO/repos"
	"github.com/Yash294/TODO/util"
	"github.com/gofiber/fiber/v2"
)

func RenderTasks(c *fiber.Ctx) error {

	userId, err := util.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	result, err := repos.GetTasks(userId)
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

	userId, err := util.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	var data models.Task
	data.Assignee = userId

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	if err = repos.AddTask(&data); err != nil {
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

func EditTask(c *fiber.Ctx) error {
	userId, err := util.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	var data models.Task
	data.Assignee = userId

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	fmt.Println(data)

	if err = repos.EditTask(&data); err != nil {
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

func DeleteTask(c *fiber.Ctx) error {
	userId, err := util.GetSessionUserID(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	var data models.Task
	data.Assignee = userId

	if err = c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error,
		})
	}

	if err = repos.DeleteTask(&data); err != nil {
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
