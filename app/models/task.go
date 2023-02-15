package models

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	TaskName string `json:"taskName" validate:"omitempty,ascii" gorm:"uniqueIndex:task"`
	Description string `json:"description"`
	Assignee uint `json:"assignee" gorm:"uniqueIndex:task"`
	IsDone   bool   `json:"isDone"`
}

type TaskDTO struct {
	ID uint `json:"id"`
	TaskName string `json:"taskName" `
	Description string `json:"description"`
	IsDone bool `json:"isDone"`
}

type TaskResponse struct {
	ID uint `json:"id"`
	TaskName string `json:"taskName" validate:"omitempty,ascii"`
	Description string `json:"description"`
	IsDone   bool   `json:"isDone"`
}
