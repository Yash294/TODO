package repos

import (
	"strings"
	"errors"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/models"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/copier"
)

type PasswordResetInfo struct {
	Email    string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

func Login(dataDTO *models.UserDTO) (uint, error) {
	// query db to check if email and passwords match
	var query models.User
	result := util.DB.Where("email = ?", strings.ToLower(dataDTO.Email)).First(&query)

	// if error is not nil, check cause, otherwise return nil for success
	if result.Error != nil {
		// if record not found, user input is incorrect, throw an error
		// otherwise the error is unrelated to user so throw
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("email/password is incorrect")
		} else {
			return 0, errors.New("failed to retrieve login info")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(query.Password), []byte(dataDTO.Password)); err != nil {
		return 0, errors.New("passwords do not match")
	}

	return query.ID, nil
}

func CreateUser(dataDTO *models.UserDTO) error {
	// convert DTO
	var dataRepo models.User
	if err := copier.Copy(&dataRepo, &dataDTO); err != nil {
		return errors.New("cannot map data")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(dataRepo.Password), 14)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	dataRepo.Password = string(bytes)
	dataRepo.Email = strings.ToLower(dataRepo.Email)

	// create the user as expected
	result := util.DB.Create(&dataRepo)

	// if unsuccessful, throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to create requested user")
	}
	return nil
}

func ChangePassword(dataDTO *models.UserDTO) error {
	// query db to see if user credentials exist
	var query models.User
	result := util.DB.Where("email = ?", strings.ToLower(dataDTO.Email)).First(&query)

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

	if err := bcrypt.CompareHashAndPassword([]byte(query.Password), []byte(dataDTO.Password)); err != nil {
		return errors.New("passwords do not match")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(dataDTO.NewPassword), 14)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	// otherwise now we can update the user's password
	result = util.DB.Model(models.User{}).Where("email = ?", strings.ToLower(dataDTO.Email)).Update("password", string(bytes))

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}