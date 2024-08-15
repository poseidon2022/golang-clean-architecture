package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang-clean-architecture/delivery/controllers"
	"golang-clean-architecture/domain"
	"golang-clean-architecture/domain/mocks"
	"golang-clean-architecture/infrastructure"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	router 				*gin.Engine
	userController		*controllers.UserController
	taskController		*controllers.TaskController
	mockUserUseCase 	*mocks.UserUseCase
	mockTaskUseCase		*mocks.TaskUseCase
	TaskGroup			[]*domain.Task
	SingleTask			domain.Task
}

func (suite *ControllerTestSuite) SetupTest() {
	suite.router = gin.Default()
	suite.mockUserUseCase = new(mocks.UserUseCase)
	suite.mockTaskUseCase = new(mocks.TaskUseCase)
	suite.userController = &controllers.UserController{
		UserUseCase : suite.mockUserUseCase,
	}
	suite.taskController = &controllers.TaskController {
		TaskUseCase : suite.mockTaskUseCase,
	}

	suite.TaskGroup = []*domain.Task{
		{
			Title : "Title 1",
			Description : "this is title 1",
			DueDate : time.Now(),
			Status : "pending",
		},
		{
			Title : "Title 2",
			Description : "this is title 2",
			DueDate : time.Now(),
			Status : "pending",
		},
	}

	suite.SingleTask = domain.Task{Title : "Title 1", Description : "this is title 1",Status : "pending",}
	suite.router.POST("/register", suite.userController.Register())
	suite.router.POST("/login", suite.userController.Login())
	suite.router.PUT("/promote/:id", infrastructure.AuthMiddleWare(), suite.userController.PromoteUser())
	suite.router.GET("/tasks", infrastructure.AuthMiddleWare(), suite.taskController.GetTasks())
	suite.router.POST("/tasks", infrastructure.AuthMiddleWare(), suite.taskController.PostTask())
	suite.router.DELETE("/tasks/:id", infrastructure.AuthMiddleWare(), suite.taskController.DeleteTask())
	suite.router.PUT("/tasks/:id", infrastructure.AuthMiddleWare(), suite.taskController.UpdateTask())
	suite.router.GET("/tasks/:id", infrastructure.AuthMiddleWare(), suite.taskController.GetTask())
}


func (suite *ControllerTestSuite) GenerateToken(email string, role string) (string, error) {

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" :email,
		"role" : role,
		"exp" : time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("error while generating token")
	}

	return token, nil
}
func (suite *ControllerTestSuite) TestRegisterSuccess() {
    user := domain.User{
        Email:    "newuser@example.com",
        Password: "password123",
    }

    suite.mockUserUseCase.On("Register", &user).Return(nil)
    body, err := json.Marshal(user)
    suite.NoError(err, "error while marshalling user data")

    req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
    suite.NoError(err, "error while creating new request")
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusOK, recorder.Code)

    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err, "error while unmarshalling response body")
    suite.Equal(gin.H{"message": "User registered Successfully"}, responseBody)

    suite.mockUserUseCase.AssertCalled(suite.T(), "Register", &user)
}

func (suite *ControllerTestSuite) TestRegisterUserAlreadyExists() {
    user := domain.User{
        Email:    "existing@example.com",
        Password: "password123",
    }

	suite.mockUserUseCase.On("Register", &user).Return(errors.New("user already exists"))

    body, _ := json.Marshal(user)
    req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusInternalServerError, recorder.Code)

    var responseBody gin.H
    err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)
    suite.Equal(gin.H{"error": "user already exists"}, responseBody)

	suite.mockUserUseCase.AssertCalled(suite.T(), "Register", &user)
}

func (suite *ControllerTestSuite) TestLoginSuccess() {
    user := domain.User{
		Email:    "newuser@example.com",
		Password: "password123",
	}
	
	suite.mockUserUseCase.On("Login", &user).Return("a_mock_token", nil)
	
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	recorder := httptest.NewRecorder()
	suite.router.ServeHTTP(recorder, req)
	
	suite.Equal(http.StatusOK, recorder.Code)
	
	var responseBody gin.H
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	suite.NoError(err)
	expectedMessage := "logged in successfully, here is your token: a_mock_token"
	suite.Equal(gin.H{"message": expectedMessage}, responseBody)
	
	suite.mockUserUseCase.AssertCalled(suite.T(), "Login", &user)
}

func (suite *ControllerTestSuite) TestPromoteUserSuccess() {
    // Create a new HTTP POST request to the /promote/12345 endpoint without a request body
    req, err := http.NewRequest(http.MethodPut, "/promote/12345", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)
    suite.mockUserUseCase.On("PromoteUser", "12345").Return(nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusOK, recorder.Code)

    // Unmarshal the response body to a gin.H map
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected success message and new role
    suite.Equal("user with the given ID promoted to admin", responseBody["message"])
    // Assert that the PromoteUser method was called with the correct user ID
    suite.mockUserUseCase.AssertCalled(suite.T(), "PromoteUser", "12345")
}

func (suite *ControllerTestSuite) TestPromote_UserUnauthorized() {
    // Create a new HTTP PUT request to the /promote/12345 endpoint without a request body
    req, err := http.NewRequest(http.MethodPut, "/promote/12345", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "user")
    suite.NoError(err)
    suite.mockUserUseCase.On("PromoteUser", "12345").Return(nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusForbidden
    suite.Equal(http.StatusForbidden, recorder.Code)

    // Unmarshal the response body to a gin.H map
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected unauthorized error
    suite.Equal("You are not authorized to promote another user", responseBody["error"])

    // Assert that the PromoteUser method was not called with the correct user ID
    suite.mockUserUseCase.AssertNotCalled(suite.T(), "PromoteUser", "12345")
}


func (suite *ControllerTestSuite) TestPromote_UserAlreadyanAdmin() {
    // Create a new HTTP PUT request to the /promote/12345 endpoint without a request body
    req, err := http.NewRequest(http.MethodPut, "/promote/12345", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)

    // Mock the PromoteUser method to return an error indicating the user is already an admin
    suite.mockUserUseCase.On("PromoteUser", "12345").Return(errors.New("user is already an admin"))

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusForbidden
    suite.Equal(http.StatusBadRequest, recorder.Code)

    // Unmarshal the response body to a gin.H map
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected error message
    suite.Equal("user is already an admin", responseBody["error"])

    // Assert that the PromoteUser method was called with the correct user ID
    suite.mockUserUseCase.AssertCalled(suite.T(), "PromoteUser", "12345")
}

func (suite *ControllerTestSuite) TestGetTasksSuccess() {
    // Create a new HTTP GET request to the /tasks endpoint without a request body
    returnedTask := suite.TaskGroup
    req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "user")
    suite.NoError(err)

    // Mock the GetTasks method to return the expected tasks
    suite.mockTaskUseCase.On("GetTasks").Return(returnedTask, nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusOK, recorder.Code)

    // Unmarshal the response body to a slice of domain.Task
    var responseBody []*domain.Task
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected tasks
    suite.Equal(returnedTask[0].Title, responseBody[0].Title)
	suite.Equal(returnedTask[0].Status, responseBody[0].Status)

    // Assert that the GetTasks method was called
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTasks")
}

func (suite *ControllerTestSuite) TestGetTasks_ErrorFetchingTasks() {
    // Create a new HTTP GET request to the /tasks endpoint without a request body
    req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "user")
    suite.NoError(err)

    // Mock the GetTasks method to return the expected tasks
    suite.mockTaskUseCase.On("GetTasks").Return(nil, errors.New("error while fetching tasks"))

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusInternalServerError, recorder.Code)

    // Unmarshal the response body to a slice of domain.Task
    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected tasks
    suite.Equal("internal server error", responseBody["error"])

    // Assert that the GetTasks method was called
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTasks")
}

func (suite *ControllerTestSuite) TestPostTaskSuccess() {
    task := suite.SingleTask
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)

    suite.mockTaskUseCase.On("PostTask", task).Return(nil)
    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)
    suite.NoError(err, "no error while creating new request")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusOK, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")
    suite.Equal(map[string]interface{}{"message": "task added successfully"}, responseBody)

    suite.mockTaskUseCase.AssertCalled(suite.T(), "PostTask", task)
}

func (suite *ControllerTestSuite) TestPostTask_UserNotAuthorized() {
    task := suite.SingleTask
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "user")
    suite.NoError(err)

    suite.mockTaskUseCase.On("PostTask", task).Return(nil)
    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)
    suite.NoError(err, "no error while creating new request")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusForbidden, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")
    suite.Equal(map[string]interface{}{"error": "You are not authorized to post a task"}, responseBody)

    suite.mockTaskUseCase.AssertNotCalled(suite.T(), "PostTask", task)
}

func (suite *ControllerTestSuite) TestPostTask_ErrorWhileAddingTask() {
    task := suite.SingleTask
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)

    suite.mockTaskUseCase.On("PostTask", task).Return(errors.New("error while trying to insert data"))
    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)
    suite.NoError(err, "no error while creating new request")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusInternalServerError, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")
    suite.Equal(map[string]interface{}{"error": "internal server error"}, responseBody)

    suite.mockTaskUseCase.AssertCalled(suite.T(), "PostTask", task)
}

func (suite *ControllerTestSuite) TestDeleteTaskSuccess() {

    req, err := http.NewRequest(http.MethodDelete, "/tasks/12345", nil)
    suite.NoError(err)

    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)
    suite.mockTaskUseCase.On("DeleteTask", "12345").Return(nil)

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)


    suite.Equal(http.StatusOK, recorder.Code)

    var responseBody gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    suite.Equal("task deleted successfully", responseBody["message"])
    suite.mockTaskUseCase.AssertCalled(suite.T(), "DeleteTask", "12345")
}

func (suite *ControllerTestSuite) TestUpdateTaskSuccess() {
    task := suite.SingleTask
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "admin")
    suite.NoError(err)

    suite.mockTaskUseCase.On("UpdateTask", "12345", &task).Return(nil)
    body, err := json.Marshal(task)
    suite.NoError(err, "no error while marshalling task data")

    req, err := http.NewRequest(http.MethodPut, "/tasks/12345", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)
    suite.NoError(err, "no error while creating new request")

    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    suite.Equal(http.StatusOK, recorder.Code)

    // Log the response body
    responseBodyBytes := recorder.Body.Bytes()
    fmt.Println("Response Body:", string(responseBodyBytes))

    var responseBody map[string]interface{}
    err = json.Unmarshal(responseBodyBytes, &responseBody)
    suite.NoError(err, "no error while unmarshalling response body")
    suite.Equal(map[string]interface{}{"message": "task updated successfully"}, responseBody)

    suite.mockTaskUseCase.AssertCalled(suite.T(), "UpdateTask", "12345", &task)
}

func (suite *ControllerTestSuite) TestGetTaskByIDSuccess() {
    // Create a new HTTP GET request to the /tasks endpoint without a request body
    returnedTask := suite.SingleTask
    req, err := http.NewRequest(http.MethodGet, "/tasks/12345", nil)
    suite.NoError(err)

    // Generate a token for the user
    token, err := suite.GenerateToken("kidusm3l@gmail.com", "user")
    suite.NoError(err)

    // Mock the GetTasks method to return the expected tasks
    suite.mockTaskUseCase.On("GetTask", "12345").Return(returnedTask, nil)

    // Set the Content-Type and Authorization headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + token)

    // Create a new HTTP test recorder
    recorder := httptest.NewRecorder()
    suite.router.ServeHTTP(recorder, req)

    // Assert that the response status code is http.StatusOK
    suite.Equal(http.StatusOK, recorder.Code)

    // Unmarshal the response body to a slice of domain.Task
    var responseBody domain.Task
    err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
    suite.NoError(err)

    // Assert that the response body contains the expected tasks
    suite.Equal(returnedTask.Title, responseBody.Title)
	suite.Equal(returnedTask.Status, responseBody.Status)

    // Assert that the GetTasks method was called
    suite.mockTaskUseCase.AssertCalled(suite.T(), "GetTask", "12345")
}


func TestControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

