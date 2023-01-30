package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskName string `json:"taskName" validate:"omitempty,ascii" gorm:"uniqueIndex:task"`
	Description string `json:"description"`
	Assignee string `json:"assignee" gorm:"uniqueIndex:task"`
	IsDone   bool   `json:"isDone"`
}
