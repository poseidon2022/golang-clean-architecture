package controllers

import (
	"fmt"
	"golang-clean-architecture/domain"
	"net/http"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUseCase		domain.UserUseCase
}

type TaskController struct {
	TaskUseCase 	domain.TaskUseCase
}

func (uc *UserController) Register() gin.HandlerFunc {

	return func(c *gin.Context) {
		var newUser domain.User
		if err := c.BindJSON(&newUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error":"invalid signup format"})
			return
		}
		
		err := uc.UserUseCase.Register(&newUser)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message" : "User registered Successfully"})
	}
}

func (uc *UserController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user domain.User
		if err := c.BindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid user format"})
			return
		}
	
		token, err := uc.UserUseCase.Login(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message" : fmt.Sprintf("logged in successfully, here is your token: %v", token)})
	}
}

func (uc *UserController) PromoteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		AuthUser, ok := c.Get("AuthorizedUser")
		if !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error":"Authorization error"})
			return
		}

		AuthorizedUser := AuthUser.(*domain.AuthenticatedUser)

		if AuthorizedUser.Role != "admin" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error":"You are not authorized to promote another user"})
			return
		}

		userID := c.Param("id")
		err := uc.UserUseCase.PromoteUser(userID)

		if err != nil {

			if err.Error() == "invalid user ID" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
				return 
			}
			if err.Error() == "no user with the specified id found" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
				return
			}
			if err.Error() == "user is already an admin" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
				return
			}
			if err.Error() == "internal server error" {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
				return
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message" : "user with the given ID promoted to admin"})
	}
}

func (tc *TaskController) GetTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("AuthorizedUser")
		if !ok  {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error" : "You are not Authenticated to perform this task"})
			return
		}

		tasks, err := tc.TaskUseCase.GetTasks()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : "internal server error"})
			return
		}
		c.IndentedJSON(http.StatusOK, tasks)
	}
}

func (tc *TaskController) GetTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		task_id := c.Param("id")
		task, err := tc.TaskUseCase.GetTask(task_id)
		if err!=nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, task)
	}
}

func (tc *TaskController) PostTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthUser, ok := c.Get("AuthorizedUser")
		if !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error":"Authorization error"})
			return
		}

		AuthorizedUser := AuthUser.(*domain.AuthenticatedUser)
		
		if AuthorizedUser.Role != "admin" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error":"You are not authorized to post a task"})
			return
		}

		var task domain.Task
		if err := c.BindJSON(&task); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "invalid input format"})
			return
		}

		err := tc.TaskUseCase.PostTask(task)
		if err != nil {
			if err.Error() == "error while trying to insert data" {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : "internal server error"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
				return	
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message" : "task added successfully"})
	}
}

func (tc *TaskController) DeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {

		AuthUser, ok := c.Get("AuthorizedUser")
		if !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error":"Authorization error"})
			return
		}

		AuthorizedUser := AuthUser.(*domain.AuthenticatedUser)
		
		if AuthorizedUser.Role != "admin" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error":"You are not authorized to delete a task"})
			return
		}

		task_id := c.Param("id")
		err := tc.TaskUseCase.DeleteTask(task_id)
		if err!=nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message" : "task deleted successfully"})
	}
}

func (tc *TaskController) UpdateTask() gin.HandlerFunc{
	return func(c *gin.Context) {

		AuthUser, ok := c.Get("AuthorizedUser")
		if !ok {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error":"Authorization error"})
			return
		}

		AuthorizedUser := AuthUser.(*domain.AuthenticatedUser)
		
		if AuthorizedUser.Role != "admin" {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error":"You are not authorized to update a task"})
			return
		}

		task_id := c.Param("id")
		var updatedTask domain.Task
		if err := c.BindJSON(&updatedTask); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "invalid input format"})
			return
		}

		err := tc.TaskUseCase.UpdateTask(task_id, &updatedTask)
		if err != nil {
			if err.Error() == "error while trying to delete data" {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : "internal server error"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
				return	
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message" : "task updated successfully"})
	}
}