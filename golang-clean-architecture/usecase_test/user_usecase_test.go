package usecase_test

import (
	"errors"
	"golang-clean-architecture/domain"
	"golang-clean-architecture/domain/mocks"
	"golang-clean-architecture/use_cases"
	"testing"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	mockRepo		*mocks.UserRepository
	useCase			domain.UserUseCase
	
}

func (suite *UserTestSuite) SetupSuite() {
	suite.mockRepo = new(mocks.UserRepository)
	suite.useCase = use_cases.NewUserUseCase(suite.mockRepo)
}

func (suite *UserTestSuite) TestUserRegister_Positive() {

    user := &domain.User{Email: "test@example.com", Password: "password123", Role: "admin"}
    suite.mockRepo.On("Register", user).Return(nil)
    suite.mockRepo.On("VerifyFirst", user).Return(nil)
    suite.mockRepo.On("UserExists", user).Return(nil)
    err := suite.useCase.Register(user)
    suite.NoError(err, "no error when creating a user")
    suite.mockRepo.AssertCalled(suite.T(), "Register", user)
}

func (suite *UserTestSuite) TestUserRegister_DatabaseError() {

	user := &domain.User{Email: "test@example.com", Password: "password123", Role: "admin"}
	suite.mockRepo.On("Register", user).Return(errors.New("database error"))
	suite.mockRepo.On("VerifyFirst", user).Return(nil)
	suite.mockRepo.On("UserExists", user).Return(nil)
	err := suite.useCase.Register(user)
	suite.Error(err, "error when creating a user")
	suite.Equal(err.Error(), "database error")
	suite.mockRepo.AssertCalled(suite.T(), "Register", user)
}

func (suite *UserTestSuite) TestUserRegister_UserAlreadyExists() {
	user := &domain.User{Email: "test@example.com", Password: "password123", Role: "admin"}
	suite.mockRepo.On("Register", user).Return(nil)
	suite.mockRepo.On("VerifyFirst", user).Return(nil)
	suite.mockRepo.On("UserExists", user).Return((errors.New("user already exists")))
	err := suite.useCase.Register(user)
	suite.Error(err, "expected error when user already exists")
	suite.Equal(err.Error(), "user already exists")
	suite.mockRepo.AssertCalled(suite.T(), "UserExists", user)
}

// func (suite *UserTestSuite) TestUserLogin_Positive() {
//     // Prepare the input data
//     userInfo := &domain.User{Email: "test@example.com", Password: "password"}
//     foundUser := domain.User{Email: "test@example.com", Password: "jjfvnfnfvjbiehei uwlxu"}

//     // Mock the repository call to GetUserByEmail
//     suite.mockRepo.On("GetUserByEmail", "test@example.com").Return(foundUser)

//     // Call the Login method
//     token, err := suite.useCaase.Login(userInfo)

//     // Assert that the expected token is returned
//     suite.NotEmpty(token, "Token should not be empty")
//     suite.NoError(err, "No error should occur when logging in")

//     // Verify that the repository method was called with the expected arguments
//     suite.mockRepo.AssertCalled(suite.T(), "GetUserByEmail", "test@example.com")
// }

func (suite *UserTestSuite) PromoteUser_Positive() {
	user_id := "valid_user_id"
	suite.mockRepo.On("PromoteUser", user_id).Return(nil)
	err := suite.useCase.PromoteUser(user_id)

	suite.NoError(err, "no error in promoting the user")
	suite.mockRepo.AssertCalled(suite.T(), "PromoteUser", user_id)
}

func (suite *UserTestSuite) PromoteUser_UserAlreadyanAdmin() {
	user_id := "admin_user_id"
	suite.mockRepo.On("PromoteUser", user_id).Return("user already an admin")
	err := suite.useCase.PromoteUser(user_id)

	suite.Error(err, "error in promoting the user")
	suite.Equal(err, "user already an admin")
	suite.mockRepo.AssertCalled(suite.T(), "PromoteUser", user_id)
}

func (suite *UserTestSuite) PromoteUser_NoUserWithSpecifiedID() {
	user_id := "valid_user_id"
	suite.mockRepo.On("PromoteUser", user_id).Return("no user with specified id")
	err := suite.useCase.PromoteUser(user_id)

	suite.Error(err, "error in promoting the user")
	suite.Equal(err, "no user with specified ID")
	suite.mockRepo.AssertCalled(suite.T(), "PromoteUser", user_id)
}


func TestUserTestSuite(t *testing.T) {
    suite.Run(t, new(UserTestSuite))
}