package repository_test

import (
	"context"
	"golang-clean-architecture/domain"
	"golang-clean-architecture/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskTestSuite struct {
	suite.Suite
	client			*mongo.Client
	db 				*mongo.Database
	collection		*mongo.Collection
	repo			domain.TaskRepository
}

func (suite	*TaskTestSuite)	SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		suite.T().Fatal(err)
	}

	db := client.Database("test")
	collection := db.Collection("tasks")
	repo := repository.NewTaskRepository(db, "tasks")
	suite.client = client
	suite.db = db
	suite.collection = collection
	suite.repo = repo
}

func (suite *TaskTestSuite) TearDownSuite() {
	err := suite.db.Drop(context.TODO())
	if err != nil {
		suite.T().Fatal(err)
	}

	err = suite.client.Disconnect(context.TODO())
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TaskTestSuite) TestPostTask() {
    task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

    err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")

    var insertedTask domain.Task
    err = suite.collection.FindOne(context.TODO(), bson.M{"_id": task.ID}).Decode(&insertedTask)
    suite.NoError(err, "no error retrieving a task")
    suite.Equal(task.Title, insertedTask.Title, "Titles should match")
    suite.Equal(task.Description, insertedTask.Description, "Descriptions should match")
    suite.Equal(task.Status, insertedTask.Status, "Statuses should match")
}

func (suite *TaskTestSuite) TestDeleteTask_Positive() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")
	err = suite.repo.DeleteTask(task.ID.Hex())
	suite.NoError(err, "no error while deleting a task")
}

func (suite *TaskTestSuite) TestDeleteTask_Negative() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")
	err = suite.repo.DeleteTask(primitive.NewObjectID().Hex())
	suite.Error(err, "error while deleting a task")
}

func (suite *TaskTestSuite) TestUpdateTask_Positive() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")

	modifiedTask := &domain.Task{
        Title:       "this task is modified",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err = suite.repo.UpdateTask(task.ID.Hex(), modifiedTask)
	suite.NoError(err, "no error while updating a task")
}

func (suite *TaskTestSuite) TestUpdateTask_Negative() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")

	modifiedTask := &domain.Task{
        Title:       "this task is modified",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

	err = suite.repo.UpdateTask(primitive.NewObjectID().Hex(), modifiedTask)
	suite.Error(err, "error while updating a task")
}

func (suite *TaskTestSuite) TestGetTask_Positive() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

    err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")
	_, err = suite.repo.GetTask(task.ID.Hex())
    suite.NoError(err, "no error retrieving a task")
}

func (suite *TaskTestSuite) TestGetTask_Negative() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

    err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")
	_, err = suite.repo.GetTask(primitive.NewObjectID().Hex())
    suite.Error(err, "error retrieving a task")
}

func (suite *TaskTestSuite) TestGetTasks() {
	task := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "this is yet another task",
        Description: "to be done by tomorrow",
        DueDate:     time.Now(),
        Status:      "pending",
    }

    err := suite.repo.PostTask(task)
    suite.NoError(err, "no error while inserting a task")
	_, err = suite.repo.GetTasks()
    suite.NoError(err, "no error retrieving a task")
}

func TestTaskTestSuite(t *testing.T) {
    suite.Run(t, new(TaskTestSuite))
}


