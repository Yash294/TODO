package repos

import (
	"errors"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/models"
	"gorm.io/gorm"
)

func Login(username string, password string) error {
	// query db to check if username and passwords match
	var query models.User
	result := database.DB.Model(models.User{}).Select("username", "password").Where("username = ? AND password = ?", username, password).First(&query)

	// if error is not nil, check cause, otherwise return nil for success
	if result.Error != nil {
		// if record not found, user input is incorrect, throw an error
		// otherwise the error is unrelated to user so throw
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("username/password is incorrect")
		} else {
			return errors.New("failed to retrieve login info")
		}
	}
	return nil
}

func CreateUser(username string, password string) error {
	// create the user as expected
	var user = models.User{Username: username, Password: password}
	result := database.DB.Create(&user)

	// if unsuccessful, throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to create requested user")
	}
	return nil
}

func ChangePassword(username string, currentPassword string, newPassword string) error {
	// query db to see if user credentials exist
	var query models.User
	result := database.DB.Where("username = ? AND password = ?", username, currentPassword).First(&query)

	// is the error is not nil check cause
	if result.Error != nil {
		// if no record found, that means user input is incorrect
		// so throw an error. If other error, throw it
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("username/password is incorrect")
		} else {
			return errors.New("failed to retrieve login info")
		}
	}

	// otherwise now we can update the user's password
	result = database.DB.Model(models.User{}).Where("username = ?", username).Update("password", newPassword)

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}

func IsUsernameAvailable(username string) (bool, error) {
	// query db to see if username exists already
	var query string
	result := database.DB.Model(models.User{}).Select("username").Where("username = ?", username).First(&query)

	// if the error is not nil, check cause
	if result.Error != nil {
		// if no record found, username available so return true
		// otherwise return error
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true, nil
		} else {
			return false, errors.New("failed to retrieve username")
		}
	}

	// username not available so return false
	return false, nil
}
