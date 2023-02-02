package repos

import (
	"errors"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/models"
)

type TaskInfo struct {
	TaskName string `json:"task_name"`
	Description string `json:"description"`
	IsDone bool `json:"isDone"`
}

func GetTasks(userId uint) ([]TaskInfo, error) {
	var query []TaskInfo
	result := util.DB.Model(models.Task{}).Where("assignee = ?", userId).Find(&query)

	if result.Error != nil {
		return query, errors.New("failed to retrieve tasks")
	}
	return query, nil
}

func AddTask(data *models.Task) error {
	result := util.DB.Model(models.Task{}).Create(&data)

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}

func EditTask(data *models.Task) error {
	result := util.DB.Model(models.Task{}).Select("task_name", "description", "is_done").Where("assignee = ? AND task_name = ?", data.Assignee, data.TaskName).Updates(&data)

	if result.Error != nil {
		return errors.New("failed to update task")
	}
	return nil
}

func DeleteTask(data *models.Task) error {
	result := util.DB.Unscoped().Where("assignee = ? AND task_name = ?", data.Assignee, data.TaskName).Delete(&models.Task{})

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}