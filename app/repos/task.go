package repos

import (
	"errors"
	"github.com/Yash294/TODO/util"
	"github.com/Yash294/TODO/app/models"
	"github.com/jinzhu/copier"
)

func GetTasks(userId uint) ([]models.TaskResponse, error) {
	var query []models.TaskResponse
	result := util.DB.Model(models.Task{}).Where("assignee = ?", userId).Find(&query)

	if result.Error != nil {
		return nil, errors.New("failed to retrieve tasks")
	}
	return query, nil
}

func AddTask(dataDTO *models.TaskDTO, userId uint) error {
	// convert DTO
	var dataRepo models.Task
	if err := copier.Copy(&dataRepo, &dataDTO); err != nil {
		return errors.New("cannot map data")
	}

	dataRepo.Assignee = userId

	result := util.DB.Model(models.Task{}).Create(&dataRepo)

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}

func EditTask(dataDTO *models.TaskDTO, userId uint) error {
	result := util.DB.Model(models.Task{}).Where("assignee = ? AND task_name = ?", userId, dataDTO.OldTaskName).Updates(map[string]interface{}{"task_name": dataDTO.TaskName, "description": dataDTO.Description, "is_done": dataDTO.IsDone})
	
	if result.Error != nil {
		return errors.New("failed to update task")
	}
	return nil
}

func DeleteTask(dataDTO *models.TaskDTO, userId uint) error {
	result := util.DB.Unscoped().Where("assignee = ? AND task_name = ?", userId, dataDTO.TaskName).Delete(&models.Task{})

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}