package repos

import (
	"errors"
	"strings"
	"github.com/Yash294/TODO/app/models"
	"github.com/Yash294/TODO/database"
	"github.com/alexedwards/argon2id"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type EncryptionStruct struct {}

type Encryption interface {
	createHash(password string) (string, error)
	comparePasswordAndHash(password string, hash string) (bool, error)
}

var EncryptionInstance Encryption = new(EncryptionStruct)

func (encrypt *EncryptionStruct) createHash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (encrypt *EncryptionStruct) comparePasswordAndHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

type CopierStruct struct {}

type Copier interface {
	copy(toValue interface{}, fromValue interface {}) error
}

var CopierInstance Copier = new(CopierStruct)

func (cop *CopierStruct) copy(toValue interface{}, fromValue interface {}) error {
	return copier.Copy(toValue, fromValue)
}

func Login(dataDTO *models.UserDTO, db *gorm.DB, encryption Encryption) (uint, error) {
	// query db to check if email and passwords match
	var query models.User
	result := db.Where("email = ?", strings.ToLower(dataDTO.Email)).First(&query)

	// if error is not nil, check cause, otherwise return nil for success
	if result.Error != nil {
		// if record not found, user input is incorrect, throw an error
		// otherwise the error is unrelated to user so throw
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("email does not exist")
		} else {
			return 0, errors.New("failed to retrieve login info")
		}
	}

	if match, err := encryption.comparePasswordAndHash(dataDTO.Password, query.Password); !match || err != nil {
		return 0, errors.New("password is incorrect")
	}

	return query.ID, nil
}

func CreateUser(dataDTO *models.UserDTO, db *gorm.DB, encryption Encryption, copy Copier) error {
	// convert DTO
	var dataRepo = &models.User{
		Email: "hello@gmail.com",
		Password: "password",
	}

	if err := copy.copy(dataRepo, dataDTO); err != nil {
		return errors.New("cannot map data")
	}

	hash, err := encryption.createHash(dataDTO.Password)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	dataRepo.Password = hash
	dataRepo.Email = strings.ToLower(dataRepo.Email)

	// create the user as expected
	result := db.Create(&dataRepo)

	// if unsuccessful, throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to create requested user")
	}
	return nil
}

func ChangePassword(dataDTO *models.UserDTO) error {
	// query db to see if user credentials exist
	var query models.User
	result := database.DB.Where("email = ?", strings.ToLower(dataDTO.Email)).First(&query)

	// is the error is not nil check cause
	if result.Error != nil {
		// if no record found, that means user input is incorrect
		// so throw an error. If other error, throw it
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("email is incorrect")
		} else {
			return errors.New("failed to retrieve login info")
		}
	}

	if match, err := argon2id.ComparePasswordAndHash(dataDTO.Password, query.Password); !match || err != nil {
		return errors.New("old password is incorrect")
	}

	hash, err := argon2id.CreateHash(dataDTO.NewPassword, argon2id.DefaultParams)
	
	if err != nil {
		return errors.New("failed to hash password")
	}

	// otherwise now we can update the user's password
	result = database.DB.Model(models.User{}).Where("email = ?", strings.ToLower(dataDTO.Email)).Update("password", hash)

	// if update not successful, then throw an error, otherwise return nil
	if result.Error != nil {
		return errors.New("failed to update user password")
	}
	return nil
}

func GetUser(userId uint) (models.UserResponse, error) {
	var query models.UserResponse
	result := database.DB.Model(models.User{}).Select("email").Where("id = ?", userId).First(&query)

	if result.Error != nil {
		return models.UserResponse{}, errors.New("failed to retrieve email")
	}

	query.Email = strings.Split(query.Email, "@")[0]
	
	return query, nil
}