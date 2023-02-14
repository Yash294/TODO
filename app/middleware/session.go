package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/Yash294/TODO/redis"
)

func RequireSession(c *fiber.Ctx) error {
	sess, err := redis.Store.Get(c)

	path := c.Path()

	if path == "/" || path == "/user/login" || path == "/user/signup" || path == "/user/resetPassword" {
		return c.Next()
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authorized",
		})
	}

	if sess.Get(redis.AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Not authorized",
		})
	}

	return c.Next()
}

func GetSessionUserID(c *fiber.Ctx) (uint, error) {
	sess, err := redis.Store.Get(c)
	if err != nil {
		return 0, errors.New("session does not exist")
	}

	userId := sess.Get(redis.USER_ID)
	if userId == nil {
		return 0, errors.New("not authorized in this session")
	}
	return userId.(uint), nil
}