package repos

import (
	"strings"
	"errors"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/models"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetInfo struct {
	Email    string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func Login(data *models.User) (models.User, error) {
	// query db to check if email and passwords match
	var query models.User
	result := util.DB.Where("email = ?", strings.ToLower(data.Email)).First(&query)

	// if error is not nil, check cause, otherwise return nil for success
	if result.Error != nil {
		// if record not found, user input is incorrect, throw an error
		// otherwise the error is unrelated to user so throw
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return query, errors.New("email/password is incorrect")
		} else {
			return query, errors.New("failed to retrieve login info")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(query.Password), []byte(data.Password)); err != nil {
		return query, errors.New("passwords do not match")
	}

	return query, nil
}

func CreateUser(data *models.User) error {

	bytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	data.Password = string(bytes)
	data.Email = strings.ToLower(data.Email)

	// create the user as expected
	result := util.DB.Create(&data)

	// if unsuccessful, throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to create requested user")
	}
	return nil
}

func ChangePassword(data *PasswordResetInfo) error {
	// query db to see if user credentials exist
	var query models.User
	result := util.DB.Where("email = ?", strings.ToLower(data.Email)).First(&query)

	// is the error is not nil check cause
	if result.Error != nil {
		// if no record found, that means user input is incorrect
		// so throw an error. If other error, throw it
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("email/password is incorrect")
		} else {
			return errors.New("failed to retrieve login info")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(query.Password), []byte(data.Password)); err != nil {
		return errors.New("passwords do not match")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 14)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	data.NewPassword = string(bytes)

	// otherwise now we can update the user's password
	result = util.DB.Model(models.User{}).Where("email = ?", strings.ToLower(data.Email)).Update("password", data.NewPassword)

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}