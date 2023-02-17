package controllers

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/gofiber/template/html"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	c *fiber.Ctx
	mock sqlmock.Sqlmock
}

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	s.NoError(err)

	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.NoError(err)

	engine := html.New("../../../resources/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	s.c = app.AcquireCtx(&fasthttp.RequestCtx{})
}

func (s *Suite) TestShouldRenderLogin() {
	err := RenderLogin(s.c)
	assert.Equal(s.T(), nil, err)
}

func (s *Suite) TestShouldRenderSignup() {
	err := RenderSignup(s.c)
	assert.Equal(s.T(), nil, err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
