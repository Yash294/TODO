package repos

import (
	"errors"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Yash294/TODO/app/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	db *gorm.DB
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
	userId = uint(1)
	email = "hello@gmail.com"
	password = "password"
	newPassword = "new password"
	hash = "hash"	
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

// Tests for Login
func (s *Suite) TestShouldLogin() {

	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(userId, email, hash))

	userId, err := Login(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertCalled(s.T(), "comparePasswordAndHash", password, hash)

	s.NoError(err)
	s.Equal(uint(1), userId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotLoginEmailDoesNotExist() {

	m := new(MockEncryption)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	userId, err := Login(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertNotCalled(s.T(), "comparePasswordAndHash")

	s.Error(err)
	s.Equal(uint(0), userId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotLoginRetrieveInfoFail() {

	m := new(MockEncryption)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(sqlmock.ErrCancelled)

	userId, err := Login(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertNotCalled(s.T(), "comparePasswordAndHash")

	s.Error(err)
	s.NotEqual(err, gorm.ErrRecordNotFound)
	s.Equal(uint(0), userId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotLoginIncorrectPassword() {

	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(false, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(userId, email, hash))

	userId, err := Login(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertCalled(s.T(), "comparePasswordAndHash", password, hash)

	s.Error(err)
	s.Equal(uint(0), userId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldCreateUser() {

	n := new(MockCopier)
	m := new(MockEncryption)
	
	n.On("copy", &models.User{Email: email, Password: password}, &models.UserDTO{Email: email, Password: password, NewPassword: newPassword}).Return(nil)
	m.On("createHash", password).Return(hash, nil)
	
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, email, hash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	s.mock.ExpectCommit()

	err := CreateUser(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m, n)

	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())

}

func (s *Suite) TestShouldNotCreateUserCopyFail() {

	n := new(MockCopier)
	m := new(MockEncryption)
	
	n.On("copy", &models.User{Email: email, Password: password}, &models.UserDTO{Email: email, Password: password, NewPassword: newPassword}).Return(errors.New("cannot map data"))
	m.On("createHash", password).Return(hash, nil)

	err := CreateUser(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m, n)

	m.AssertNotCalled(s.T(), "createHash", password)

	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotCreateUserHashFail() {

	n := new(MockCopier)
	m := new(MockEncryption)
	
	n.On("copy", &models.User{Email: email, Password: password}, &models.UserDTO{Email: email, Password: password, NewPassword: newPassword}).Return(nil)
	m.On("createHash", password).Return(password, errors.New("cannot hash password"))
	
	err := CreateUser(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m, n)

	n.AssertNumberOfCalls(s.T(), "copy", 1)
	m.AssertCalled(s.T(), "createHash", password)

	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotCreateUserDBFail() {

	n := new(MockCopier)
	m := new(MockEncryption)
	
	n.On("copy", &models.User{Email: email, Password: password}, &models.UserDTO{Email: email, Password: password, NewPassword: newPassword}).Return(nil)
	m.On("createHash", password).Return(hash, nil)
	
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, email, hash).
		WillReturnError(errors.New("username exists"))

	err := CreateUser(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m, n)

	n.AssertNumberOfCalls(s.T(), "copy", 1)
	m.AssertCalled(s.T(), "createHash", password)

	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldChangePassword() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password).Return(true, nil)
	m.On("createHash", password).Return(hash, nil)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}