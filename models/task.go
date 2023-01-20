package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskName string `json:"taskName" validate:"omitempty,ascii"`
	Assignee string `json:"assignee"`
	IsDone   bool   `json:"isDone"`
}
