package util

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var (
	Store *session.Store
	AUTH_KEY string = "authenticated"
	USER_ID string = "user_id"
)

func ConnectRedis() error {

	var storage fiber.Storage

	config, err := LoadRedisConfig(".")
	if err != nil {
		panic("Cannot load config")
	}

	rds := redis.Config{
		Host:     config.RedisHost,
		Port:     config.RedisPort,
		Username: config.RedisUser,
		Password: config.RedisPass,
	}

	storage = redis.New(rds)

	Store = session.New(session.Config{
		Storage:        storage,
		Expiration:     time.Hour * 4,
		CookieHTTPOnly: true,
	})

	return nil
}
