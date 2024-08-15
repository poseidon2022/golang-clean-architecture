package use_cases

import (
	"golang-clean-architecture/domain"
	"strings"
	"errors"
)

type TaskUseCase struct {
	Repository 		domain.TaskRepository
}

func NewTaskUseCase(tr domain.TaskRepository) domain.TaskUseCase {
	return &TaskUseCase {
		Repository: tr,
	}
}

func (tu *TaskUseCase) GetTasks() ([]*domain.Task, error) {
	tasks, err := tu.Repository.GetTasks()
	return tasks, err
}

func (tu *TaskUseCase) GetTask(taskId string) (domain.Task, error) {
	task, err := tu.Repository.GetTask(taskId)
	return task, err
}

func (tu *TaskUseCase) PostTask(task domain.Task) error {

	task.Description = strings.TrimSpace(task.Description)
	task.Title = strings.TrimSpace(task.Title)
	task.Status = strings.TrimSpace(task.Status)

	if task.Description == "" || task.Status == "" || task.Title == "" {
		return errors.New("required fields are missing")
	}
	err := tu.Repository.PostTask(&task)
	return err
}

func (tu *TaskUseCase) DeleteTask(taskID string) error {
	err := tu.Repository.DeleteTask(taskID)
	return err
}

func (tu *TaskUseCase) UpdateTask(taskID string, modifiedTask *domain.Task) error {
	err := tu.Repository.UpdateTask(taskID, modifiedTask)
	return err
}