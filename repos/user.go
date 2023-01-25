package repos

import (
	"fmt"
	"errors"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/models"
	"gorm.io/gorm"
)

func Login(email string, password string) error {
	// query db to check if email and passwords match
	var query models.User
	result := database.DB.Model(models.User{}).Select("email", "password").Where("email = ? AND password = ?", email, password).First(&query)

	fmt.Println(query)
	fmt.Println(result.Error)

	// if error is not nil, check cause, otherwise return nil for success
	if result.Error != nil {
		// if record not found, user input is incorrect, throw an error
		// otherwise the error is unrelated to user so throw
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("email/password is incorrect")
		} else {
			return errors.New("failed to retrieve login info")
		}
	}
	return nil
}

func CreateUser(email string, password string) error {
	// create the user as expected
	var user = models.User{ Email: email, Password: password}
	result := database.DB.Create(&user)

	// NEED TO CHECK uniqueness
	// if unsuccessful, throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to create requested user")
	}
	return nil
}

func ChangePassword(email string, currentPassword string, newPassword string) error {
	// query db to see if user credentials exist
	var query models.User
	result := database.DB.Where("email = ? AND password = ?", email, currentPassword).First(&query)

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

	// otherwise now we can update the user's password
	result = database.DB.Model(models.User{}).Where("email = ?", email).Update("password", newPassword)

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}

func IsemailAvailable(email string) (bool, error) {
	// query db to see if email exists already
	var query string
	result := database.DB.Model(models.User{}).Select("email").Where("email = ?", email).First(&query)

	// if the error is not nil, check cause
	if result.Error != nil {
		// if no record found, email available so return true
		// otherwise return error
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return true, nil
		} else {
			return false, errors.New("failed to retrieve email")
		}
	}

	// email not available so return false
	return false, nil
}
