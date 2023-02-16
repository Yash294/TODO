package repos

import (
	"errors"
	"github.com/Yash294/TODO/app/models"
	"gorm.io/gorm"
)

func GetTasks(userId uint, db *gorm.DB) ([]models.TaskResponse, error) {
	var query []models.TaskResponse
	result := db.Model(models.Task{}).Where("assignee = ?", userId).Find(&query)

	if result.Error != nil {
		return nil, errors.New("failed to retrieve tasks")
	}
	return query, nil
}

func AddTask(dataDTO *models.TaskDTO, userId uint, db *gorm.DB, copy Copier) (uint, error) {
	var dataRepo = &models.Task {
		TaskName: "todo",
		Description: "finish the todo app",
		Assignee: uint(1),
		IsDone: false,
	}

	if err := copy.copy(dataRepo, dataDTO); err != nil {
		return 0, errors.New("cannot map data")
	}

	dataRepo.Assignee = userId

	result := db.Model(models.Task{}).Create(&dataRepo)

	if result.Error != nil {
		return 0, errors.New("failed to create new task")
	}

	return dataRepo.ID, nil
}

func EditTask(dataDTO *models.TaskDTO, db *gorm.DB) error {
	result := db.Model(models.Task{}).Where("id = ?", dataDTO.ID).Updates(map[string]interface{}{"task_name": dataDTO.TaskName, "description": dataDTO.Description, "is_done": dataDTO.IsDone})
	
	if result.Error != nil {
		return errors.New("failed to update task")
	}
	return nil
}

func DeleteTask(dataDTO *models.TaskDTO, db *gorm.DB) error {
	result := db.Unscoped().Where("id = ?", dataDTO.ID).Delete(&models.Task{})

	if result.Error != nil {
		return errors.New("failed to create new task")
	}
	return nil
}