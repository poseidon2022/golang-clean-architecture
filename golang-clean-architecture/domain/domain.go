package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id"`
	Email    string       		  `json:"email" bson:"email"`
	Password string       		  `json:"password" bson:"password"`
	Role     string				  `json:"role" bson:"role"`
}

type Task struct {
	ID         	primitive.ObjectID   `json:"id" bson:"_id"`
	Title       string    			 `json:"title" bson:"title"`
	Description string    			 `json:"description" bson:"description"`
	DueDate     time.Time 			 `json:"due_date" bson:"due_date"`
	Status      string    			 `json:"status" bson:"status"`
}

type AuthenticatedUser struct {
	Role		string
	Email		string
}

type TaskRepository interface {
	GetTasks()							([]*Task, error)
	GetTask(string)						(Task, error)
	PostTask(*Task)						error
	DeleteTask(string)					error
	UpdateTask(string, *Task)			error
}

type TaskUseCase interface {
	GetTasks()							([]*Task, error)
	GetTask(string)						(Task, error)
	PostTask(Task)						error
	DeleteTask(string)					error
	UpdateTask(string, *Task)			error
}

type UserRepository interface {
	Register(*User)						error
	VerifyFirst(*User)					error
	UserExists(*User)					error
	GetUserByEmail(string)				User
	PromoteUser(string)					error
}

type UserUseCase interface {
	Register(*User)						error
	Login(*User)						(string, error)
	PromoteUser(string)					error
}
