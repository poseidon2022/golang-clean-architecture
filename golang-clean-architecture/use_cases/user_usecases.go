package use_cases

import (
	"errors"
	"fmt"
	"golang-clean-architecture/domain"
	"golang-clean-architecture/infrastructure"
	"strings"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	Repository	domain.UserRepository
}

func NewUserUseCase(ur domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{
		Repository : ur,
	}
}

func (user  *UserUseCase) Register(newUser *domain.User) error {
	
	if newUser.Email == "" || newUser.Password == "" {
		return errors.New("required field missing")
	}

	err := user.Repository.VerifyFirst(newUser)
	if err != nil {
		if err.Error() == "a user is found on db" {
			newUser.Role = "user"
		}
		if err.Error() == "internal server error" {
			return err
		}
	} else {
		newUser.Role = "admin"
	}

	err = user.Repository.UserExists(newUser)
	if err != nil {
		return err
	}

	hashedPassword, err := infrastructure.HashPassword(newUser.Password)
	if err != nil {
		return errors.New("internal server error")
	}

	newUser.Password = string(hashedPassword)
	err = user.Repository.Register(newUser)
	if err != nil {
		return err
	}

	return nil
}


func (user *UserUseCase) Login(userInfo *domain.User) (string, error){
	userInfo.Password = strings.TrimSpace(userInfo.Password)
	userInfo.Email = strings.TrimSpace(userInfo.Email)
	if userInfo.Password == "" || userInfo.Email == "" {
		return "", errors.New("required fields are missing")
	}
	foundUser := user.Repository.GetUserByEmail(userInfo.Email)
	if foundUser == (domain.User{}) {
		return "", errors.New("invalid credentials")
	}
	
	validateUser := infrastructure.ComparePasswords(&foundUser, userInfo)
	if validateUser != nil {
		fmt.Println("here")
		return "", errors.New("invalid credentials")
	}

	token, err := infrastructure.GenerateToken(&foundUser)
	if err != nil {
		fmt.Println("here")
		return "", errors.New("internal server error")
	}

	return token, nil
}

func (user *UserUseCase) PromoteUser(userID string) error {
	err := user.Repository.PromoteUser(userID)
	return err
}