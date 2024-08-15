package usecase_test

import (
	"errors"
	"golang-clean-architecture/domain"
	"golang-clean-architecture/domain/mocks"
	"golang-clean-architecture/use_cases"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskTestSuite struct {
    suite.Suite
    taskmockRepo *mocks.TaskRepository
    taskuseCase  domain.TaskUseCase
}

func (suite *TaskTestSuite) SetupTest() {
    suite.taskmockRepo = new(mocks.TaskRepository)
    suite.taskuseCase = use_cases.NewTaskUseCase(suite.taskmockRepo)
}

func (suite *TaskTestSuite) TestGetTasks_Positive() {

    tasks := []*domain.Task{
        {ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"},
        {ID: primitive.NewObjectID(), Title: "Task 2", Description: "Description 2", Status: "pending"},
    }
    suite.taskmockRepo.On("GetTasks").Return(tasks, nil)
    result, err := suite.taskuseCase.GetTasks()
    suite.NoError(err, "no error when getting tasks")
    suite.Equal(tasks, result, "tasks should match")
    suite.taskmockRepo.AssertCalled(suite.T(), "GetTasks")
}

func (suite *TaskTestSuite) TestGetTasks_ErrorFetchingTasks() {

    tasks := []*domain.Task{
        {ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"},
        {ID: primitive.NewObjectID(), Title: "Task 2", Description: "Description 2", Status: "pending"},
    }
    suite.taskmockRepo.On("GetTasks").Return(tasks, errors.New("error while fetching tasks"))
    _, err := suite.taskuseCase.GetTasks()
    suite.Error(err, "error when getting tasks")
    suite.Equal(err.Error(), "error while fetching tasks")
    suite.taskmockRepo.AssertCalled(suite.T(), "GetTasks")
}

func (suite *TaskTestSuite) TestGetTask_Positive() {

    task := domain.Task{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"}
    suite.taskmockRepo.On("GetTask", task.ID.Hex()).Return(task, nil)
    result, err := suite.taskuseCase.GetTask(task.ID.Hex())
    suite.NoError(err, "no error when fetching a task")
    suite.Equal(task, result, "tasks should match")
    suite.taskmockRepo.AssertCalled(suite.T(), "GetTask", task.ID.Hex())
}

func (suite *TaskTestSuite) TestGetTask_InvalidTaskID() {

    task := domain.Task{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"}
    suite.taskmockRepo.On("GetTask", task.ID.Hex()).Return(task, errors.New("invalid task ID"))
	_, err := suite.taskuseCase.GetTask(task.ID.Hex())
    suite.Error(err, "error while fetching a task")
    suite.Equal(err.Error(), "invalid task ID")
    suite.taskmockRepo.AssertCalled(suite.T(), "GetTask", task.ID.Hex())
}

func (suite *TaskTestSuite) TestGetTask_NoTaskWithSpecifiedID() {

    task := domain.Task{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"}
    suite.taskmockRepo.On("GetTask", task.ID.Hex()).Return(task, errors.New("no task with specified ID"))
	_, err := suite.taskuseCase.GetTask(task.ID.Hex())
    suite.Error(err, "error while fetching a task")
    suite.Equal(err.Error(), "no task with specified ID")
    suite.taskmockRepo.AssertCalled(suite.T(), "GetTask", task.ID.Hex())
}

func (suite *TaskTestSuite) TestDeleteTask_Positive() {

    task := domain.Task{ID: primitive.NewObjectID(), Title: "Task 1", Description: "Description 1", Status: "pending"}
    suite.taskmockRepo.On("DeleteTask", task.ID.Hex()).Return(nil)
    err := suite.taskuseCase.DeleteTask(task.ID.Hex())
    suite.NoError(err, "no error when deleting a task")
    suite.taskmockRepo.AssertCalled(suite.T(), "DeleteTask", task.ID.Hex())
}

func (suite *TaskTestSuite) TestDeleteTask_InvalidTaskID() {

	task_id := primitive.NewObjectID().Hex()
    suite.taskmockRepo.On("DeleteTask", task_id).Return(errors.New("invalid task ID"))
	err := suite.taskuseCase.DeleteTask(task_id)
    suite.Error(err, "error while deleting a task")
    suite.Equal(err.Error(), "invalid task ID")
    suite.taskmockRepo.AssertCalled(suite.T(), "DeleteTask", task_id)
}

func (suite *TaskTestSuite) TestDeleteTask_TaskNotFound() {

	task_id := primitive.NewObjectID().Hex()
    suite.taskmockRepo.On("DeleteTask", task_id).Return(errors.New("task with specified id not found"))
	err := suite.taskuseCase.DeleteTask(task_id)
    suite.Error(err, "error while deleting a task")
    suite.Equal(err.Error(), "task with specified id not found")
    suite.taskmockRepo.AssertCalled(suite.T(), "DeleteTask", task_id)
}

func (suite *TaskTestSuite) TestUpdateTask_Positive() {

    modifiedTask := domain.Task{Title: "Task 2", Description: "Description 1", Status: "pending"}
	task_id := primitive.NewObjectID().Hex()
    suite.taskmockRepo.On("UpdateTask", task_id, &modifiedTask).Return(nil)
	err := suite.taskuseCase.UpdateTask(task_id, &modifiedTask)
    suite.NoError(err, "no error while fetching a task")
    suite.taskmockRepo.AssertCalled(suite.T(), "UpdateTask", task_id, &modifiedTask)

}

func (suite *TaskTestSuite) TestUpdateTask_TaskNotFound() {

	modifiedTask := domain.Task{Title: "Task 2", Description: "Description 1", Status: "pending"}
	task_id := primitive.NewObjectID().Hex()
    suite.taskmockRepo.On("UpdateTask", task_id, &modifiedTask).Return(errors.New("task with specified id not found"))
	err := suite.taskuseCase.UpdateTask(task_id, &modifiedTask)
    suite.Error(err, "error while updating a task")
    suite.Equal(err.Error(), "task with specified id not found")
    suite.taskmockRepo.AssertCalled(suite.T(), "UpdateTask", task_id, &modifiedTask)
}

func (suite *TaskTestSuite) TestPostTask_Positive() {
	insertedTask := domain.Task{Title: "Task 2", Description: "Description 1", DueDate : time.Now(), Status: "pending"}
	suite.taskmockRepo.On("PostTask", &insertedTask).Return(nil)
	err := suite.taskuseCase.PostTask(insertedTask)
	suite.NoError(err,  "no error while posting a task")
	suite.taskmockRepo.AssertCalled(suite.T(), "PostTask", &insertedTask)
}

func (suite *TaskTestSuite) TestPostTask_RequiredFieldsMissing() {
	insertedTask := domain.Task{
		Title : "",
		Description: " nicnokd",
		DueDate : time.Now(),
		Status : "on going",
	}
	suite.taskmockRepo.On("PostTask", &insertedTask).Return(errors.New("required field missing"))
	err := suite.taskuseCase.PostTask(insertedTask)
	suite.Error(err, "error while posting a task")
	suite.Equal(err.Error(), "required fields are missing")
	suite.taskmockRepo.AssertNotCalled(suite.T(), "PostTask", &insertedTask)
}

func TestTaskTestSuite(t *testing.T) {
    suite.Run(t, new(TaskTestSuite))
}