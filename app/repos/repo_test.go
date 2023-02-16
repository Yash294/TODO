package repos

import (
	"database/sql"
	"database/sql/driver"
	"errors"
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

var (
	taskId = uint(1)
	taskName = "todo"
	description = "finish the todo app"
	assignee = userId
	isDone = false
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

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)
	m.On("createHash", newPassword).Return(hash, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(userId, email, hash))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password"=$1,"updated_at"=$2 WHERE email = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs(hash, anyTime{}, email).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotChangePasswordEmailDoesNotExist() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)
	m.On("createHash", newPassword).Return(hash, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertNotCalled(s.T(), "comparePasswordAndHash", password, hash)
	m.AssertNotCalled(s.T(), "createHash", newPassword)

	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotChangePasswordRetrieveInfoFail() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)
	m.On("createHash", newPassword).Return(hash, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(sqlmock.ErrCancelled)

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertNotCalled(s.T(), "comparePasswordAndHash", password, hash)
	m.AssertNotCalled(s.T(), "createHash", newPassword)

	s.Error(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotChangePasswordIncorrectPassword() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(false, errors.New("incorrect password"))
	m.On("createHash", newPassword).Return(hash, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(userId, email, hash))

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertCalled(s.T(), "comparePasswordAndHash", password, hash)
	m.AssertNotCalled(s.T(), "createHash", newPassword)

	s.Error(err)
	s.Equal(err, errors.New("old password is incorrect"))
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotChangePasswordHashFail() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)
	m.On("createHash", newPassword).Return(password, errors.New("password could not be hashed"))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(userId, email, hash))

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertCalled(s.T(), "comparePasswordAndHash", password, hash)
	m.AssertCalled(s.T(), "createHash", newPassword)

	s.Error(err)
	s.Equal(err, errors.New("failed to hash password"))
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotChangePasswordUpdateFail() {
	m := new(MockEncryption)

	m.On("comparePasswordAndHash", password, hash).Return(true, nil)
	m.On("createHash", newPassword).Return(hash, nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(userId, email, hash))

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "password"=$1,"updated_at"=$2 WHERE email = $3 AND "users"."deleted_at" IS NULL`)).
		WithArgs(hash, anyTime{}, email).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	err := ChangePassword(&models.UserDTO{Email: email, Password: password, NewPassword: newPassword}, s.db, m)

	m.AssertCalled(s.T(), "comparePasswordAndHash", password, hash)
	m.AssertCalled(s.T(), "createHash", newPassword)

	s.Error(err)
	s.Equal(err, errors.New("failed to update user password"))
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldGetUser() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email" FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow(email))

	_, err := GetUser(userId, s.db)

	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

// TASK REPO TESTS

func (s *Suite) TestShouldGetTasks() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."task_name","tasks"."description","tasks"."is_done" FROM "tasks" WHERE assignee = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "description", "is_done"}).
		AddRow(taskId, taskName, description, isDone))

	tasks, err := GetTasks(userId, s.db)

	s.NoError(err)
	s.Equal([]models.TaskResponse{{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}}, tasks)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotGetTasks() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."task_name","tasks"."description","tasks"."is_done" FROM "tasks" WHERE assignee = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(userId).
		WillReturnError(sqlmock.ErrCancelled)

	tasks, err := GetTasks(userId, s.db)

	s.Error(err)
	s.Equal([]models.TaskResponse(nil), tasks)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldAddTask() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(nil)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("created_at","updated_at","deleted_at","task_name","description","assignee","is_done") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, taskName, description, assignee, isDone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	s.mock.ExpectCommit()

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.NoError(err)
	s.Equal(taskId, newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldFailCopy() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(errors.New("cannot map to requested format"))

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.Error(err)
	s.Equal(errors.New("cannot map data"), err)
	s.Equal(uint(0), newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotAddTask() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(nil)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("created_at","updated_at","deleted_at","task_name","description","assignee","is_done") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, taskName, description, assignee, isDone).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.Error(err)
	s.Equal(errors.New("failed to create new task"), err)
	s.Equal(uint(0), newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldEditTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"is_done"=$2,"task_name"=$3,"updated_at"=$4 WHERE id = $5 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(description, isDone, taskName, anyTime{}, taskId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := EditTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.NoError(err)
	s.Equal(nil, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotEditTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"is_done"=$2,"task_name"=$3,"updated_at"=$4 WHERE id = $5 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(description, isDone, taskName, anyTime{}, taskId).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	err := EditTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.Error(err)
	s.Equal(errors.New("failed to update task"), err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldDeleteTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1`)).
		WithArgs(taskId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := DeleteTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.NoError(err)
	s.Equal(nil, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotDeleteTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1`)).
		WithArgs(taskId).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	err := DeleteTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.Error(err)
	s.Equal(errors.New("failed to create new task"), err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
