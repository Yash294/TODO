package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `json:"email" validate:"omitempty,max=64" gorm:"uniqueIndex"`
	Password string `json:"password" validate:"omitempty,min=8,max=64"`
	Tasks    []Task `json:"tasks" gorm:"foreignKey:TaskName"`
}
