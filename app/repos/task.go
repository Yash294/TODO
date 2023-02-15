package repos

import (
	"errors"
	"github.com/Yash294/TODO/database"
	"github.com/Yash294/TODO/app/models"
	"github.com/jinzhu/copier"
)

func GetTasks(userId uint) ([]models.TaskResponse, error) {
	var query []models.TaskResponse
	result := database.DB.Model(models.Task{}).Where("assignee = ?", userId).Find(&query)

	if result.Error != nil {
		return nil, errors.New("failed to retrieve tasks")
	}
	return query, nil
}

func AddTask(dataDTO *models.TaskDTO, userId uint) (uint, error) {
	var dataRepo models.Task
	if err := copier.Copy(&dataRepo, &dataDTO); err != nil {
		return 0, errors.New("cannot map data")
	}

	dataRepo.Assignee = userId

	result := database.DB.Model(models.Task{}).Create(&dataRepo)

	if result.Error != nil {
		return 0, errors.New("failed to create new task")
	}

	var query uint
	result = database.DB.Model(models.Task{}).Select("id").Where("task_name = ?", dataRepo.TaskName).First(&query)
	
	if result.Error != nil {
		
		return 0, errors.New("failed to retrieve newly added task")
	}

	return query, nil
}

func EditTask(dataDTO *models.TaskDTO) error {
	result := database.DB.Model(models.Task{}).Where("id = ?", dataDTO.ID).Updates(map[string]interface{}{"task_name": dataDTO.TaskName, "description": dataDTO.Description, "is_done": dataDTO.IsDone})
	
	if result.Error != nil {
		return errors.New("failed to update task")
	}
	return nil
}

func DeleteTask(dataDTO *models.TaskDTO) error {
	result := database.DB.Unscoped().Where("id = ?", dataDTO.ID).Delete(&models.Task{})

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}