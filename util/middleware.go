package util

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func NewMiddleware() fiber.Handler {
	return AuthMiddleware
}

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := Store.Get(c)

	path := c.Path()
	
	fmt.Println(c.Path())

	if path == "/" || path == "/user/login" || path == "/user/signup" {
		return c.Next()
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authorized",
		})
	}

	if sess.Get(AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authorized",
		})
	}

	return c.Next()
}

func GetSessionUserID(c *fiber.Ctx) (uint, error) {
	sess, err := Store.Get(c)
	if err != nil {
		return 0, errors.New("session does not exist")
	}

	userId := sess.Get(USER_ID)
	if userId == nil {
		return 0, errors.New("not authorized in this session")
	}
	return userId.(uint), nil
}