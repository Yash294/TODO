package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `json:"email" validate:"omitempty,max=64" gorm:"uniqueIndex"`
	Password string `json:"password" validate:"omitempty,min=16,max=64"`
}

type UserDTO struct {
	Email string `json:"email" validate:"omitempty,max=64"`
	Password    string `json:"password" validate:"omitempty,min=16,max=64"`
	NewPassword string `json:"newPassword" validate:"omitempty,min=16,max=64"`
}

type UserResponse struct {
	Email string `json:"email" validate:"omitempty,max=64"`
}