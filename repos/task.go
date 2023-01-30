package repos

import (
	//"fmt"
	"errors"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/models"
	//"gorm.io/gorm"
)

type TaskInfo struct {
	TaskName string `json:"task_name"`
	Description string `json:"description"`
	IsDone bool `json:"isDone"`
}

func GetTasks(data *models.Task) ([]TaskInfo, error) {
	var query []TaskInfo
	result := database.DB.Model(models.Task{}).Select("task_name", "description", "is_done").Where("assignee = ?", data.Assignee).Find(&query)

	if result.Error != nil {
		return query, errors.New("failed to retrieve tasks")
	}
	return query, nil
}

func AddTask(data *models.Task) error {
	result := database.DB.Model(models.Task{}).Create(&data)

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}

func DeleteTask(data *models.Task) error {
	result := database.DB.Model(models.Task{}).Delete(&data)

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}

func EditTask(data *models.Task) error {
	result := database.DB.Model(models.Task{}).Updates(&data)

	if result.Error != nil {
		return errors.New("failed to update task")
	}
	return nil
}