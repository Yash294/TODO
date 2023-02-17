package repos

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type MockEncryption struct {
	mock.Mock
}

func (m *MockEncryption) createHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockEncryption) comparePasswordAndHash(password string, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}

type MockCopier struct {
	mock.Mock
}

func (n *MockCopier) copy(toValue interface{}, fromValue interface{}) error {
	args := n.Called(toValue, fromValue)
	return args.Error(0)
}

var (
	userId      = uint(1)
	email       = "hello@gmail.com"
	password    = "password"
	newPassword = "new password"
	hash        = "hash"
)

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

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}