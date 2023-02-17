package repos

import (
	"errors"
	"regexp"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Yash294/TODO/app/models"
	"gorm.io/gorm"
)

// USER REPO TESTS

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