package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" validate:"omitempty,min=5,max=16,alphanum"`
	Password string `json:"password" validate:"omitempty,min=5,max=16,alphanum"`
	Tasks    []Task `json:"tasks" gorm:"foreignKey:TaskName"`
}
