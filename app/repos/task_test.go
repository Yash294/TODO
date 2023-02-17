package repos

import (
	"errors"
	"regexp"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Yash294/TODO/app/models"
)

var (
	taskId = uint(1)
	taskName = "todo"
	description = "finish the todo app"
	assignee = userId
	isDone = false
)

// TASK REPO TESTS

func (s *Suite) TestShouldGetTasks() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."task_name","tasks"."description","tasks"."is_done" FROM "tasks" WHERE assignee = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "description", "is_done"}).
		AddRow(taskId, taskName, description, isDone))

	tasks, err := GetTasks(userId, s.db)

	s.NoError(err)
	s.Equal([]models.TaskResponse{{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}}, tasks)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotGetTasks() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT "tasks"."id","tasks"."task_name","tasks"."description","tasks"."is_done" FROM "tasks" WHERE assignee = $1 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(userId).
		WillReturnError(sqlmock.ErrCancelled)

	tasks, err := GetTasks(userId, s.db)

	s.Error(err)
	s.Equal([]models.TaskResponse(nil), tasks)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldAddTask() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(nil)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("created_at","updated_at","deleted_at","task_name","description","assignee","is_done") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, taskName, description, assignee, isDone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	s.mock.ExpectCommit()

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.NoError(err)
	s.Equal(taskId, newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldFailCopy() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(errors.New("cannot map to requested format"))

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.Error(err)
	s.Equal(errors.New("cannot map data"), err)
	s.Equal(uint(0), newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotAddTask() {
	n := new(MockCopier)

	n.On("copy", &models.Task{TaskName: taskName, Description: description, Assignee: userId, IsDone: isDone}, &models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}).Return(nil)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks" ("created_at","updated_at","deleted_at","task_name","description","assignee","is_done") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(anyTime{}, anyTime{}, nil, taskName, description, assignee, isDone).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	newTaskId, err := AddTask(&models.TaskDTO{TaskName: taskName, Description: description, IsDone: isDone}, userId, s.db, n)

	s.Error(err)
	s.Equal(errors.New("failed to create new task"), err)
	s.Equal(uint(0), newTaskId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldEditTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"is_done"=$2,"task_name"=$3,"updated_at"=$4 WHERE id = $5 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(description, isDone, taskName, anyTime{}, taskId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := EditTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.NoError(err)
	s.Equal(nil, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotEditTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"is_done"=$2,"task_name"=$3,"updated_at"=$4 WHERE id = $5 AND "tasks"."deleted_at" IS NULL`)).
		WithArgs(description, isDone, taskName, anyTime{}, taskId).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	err := EditTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.Error(err)
	s.Equal(errors.New("failed to update task"), err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldDeleteTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1`)).
		WithArgs(taskId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := DeleteTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.NoError(err)
	s.Equal(nil, err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *Suite) TestShouldNotDeleteTask() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1`)).
		WithArgs(taskId).
		WillReturnError(sqlmock.ErrCancelled)
	s.mock.ExpectRollback()

	err := DeleteTask(&models.TaskDTO{ID: taskId, TaskName: taskName, Description: description, IsDone: isDone}, s.db)

	s.Error(err)
	s.Equal(errors.New("failed to create new task"), err)
	s.NoError(s.mock.ExpectationsWereMet())
}