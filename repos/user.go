package repos

import (
	"errors"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/models"
	"gorm.io/gorm"
)

type PasswordResetInfo struct {
	Email    string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func Login(data *models.User) (models.User, error) {
	// query db to check if email and passwords match
	var query models.User
	result := util.DB.Model(models.User{}).Select("id").Where("email = ? AND password = ?", data.Email, data.Password).First(&query)

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
	return query, nil
}

func CreateUser(data *models.User) error {
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
	result := util.DB.Where("email = ? AND password = ?", data.Email, data.Password).First(&query)

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
	result = util.DB.Model(models.User{}).Where("email = ?", data.Email).Update("password", data.NewPassword)

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}