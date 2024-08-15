package repository

import (
	"context"
	"errors"
	"fmt"
	"golang-clean-architecture/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	Database		 *mongo.Database
	Collection 		 string
}

func NewTaskRepository(db *mongo.Database, collection string)  domain.TaskRepository {
	return &TaskRepository {
		Database : db,
		Collection: collection,
	}
}

func (tr *TaskRepository) GetTasks() ([]*domain.Task, error) {
	var tasks []*domain.Task
	collection := tr.Database.Collection(tr.Collection)
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("first")
		return nil, errors.New("error while fetching tasks")
	}
	
	for cur.Next(context.TODO()) {
		var task domain.Task
		err := cur.Decode(&task)
		if err != nil {
			fmt.Println("second")
			return nil, errors.New("error while fetching tasks")
		}
		tasks = append(tasks, &task)
	}

	fmt.Println(tasks)

	if cur.Err() != nil {
		fmt.Println("third")
		return nil, errors.New("error while fetching tasks")
	}

	return tasks, nil
}

func (tr *TaskRepository) GetTask(taskID string) (domain.Task, error) {
	
	collection := tr.Database.Collection(tr.Collection)
	processedID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return domain.Task{}, errors.New("invalid task id")
	}
	filter := bson.D{{Key : "_id", Value : processedID}}
	var task domain.Task
	err = collection.FindOne(context.TODO(), filter).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return domain.Task{}, errors.New("there is no task with the specified id")
	}
	if err != nil {
		return domain.Task{}, errors.New("internal server error")
	}

	return task, nil
}

func (tr *TaskRepository) PostTask(task *domain.Task) error {
	collection := tr.Database.Collection(tr.Collection)
	task.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		return errors.New("error while trying to insert data")
	}
	return nil
}

func (tr *TaskRepository) DeleteTask(task_id string) error {
	processedID, err := primitive.ObjectIDFromHex(task_id)
	if err != nil {
		return errors.New("invalid task id")
	}

	filter := bson.D{{Key : "_id", Value : processedID}}

	collection := tr.Database.Collection(tr.Collection)
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if deleteResult.DeletedCount == 0{
		return errors.New("task with the specified id not found")
	}
	if err != nil {
		return errors.New("internal sever error")
	}
	return nil
}

func (tr *TaskRepository) UpdateTask(id string, modified *domain.Task) error {

	processedID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task id")
	}

	filter := bson.D{{Key : "_id", Value : processedID}}
	collection := tr.Database.Collection(tr.Collection)

	update := bson.D{}

	if modified.Title != "" {
		update = append(update, bson.E{
			Key : "$set", Value : bson.D{{
				Key : "title", Value : modified.Title,
			}}})
	}

	if modified.Status != "" {
		update = append(update, bson.E{
			Key : "$set", Value : bson.D{{
				Key : "status", Value : modified.Status,
			}}})
	}

	if modified.Description != "" {
		update = append(update, bson.E{
			Key : "$set", Value : bson.D{{
				Key : "description", Value : modified.Description,
			}}})
	}

	updatedResult, err := collection.UpdateOne(context.TODO(), filter, update)

	if updatedResult.MatchedCount == 0 {
		return errors.New("document with the specified id not found")
	}
	if err != nil || updatedResult.ModifiedCount == 0 {
		return errors.New("internal sever error")
	}

	return nil
}
