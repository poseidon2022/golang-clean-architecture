package infrastructure

import (
	"errors"
	"golang-clean-architecture/domain"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, errors.New("error while hashing password")
	}
	return hashedPassword, nil
}

func ComparePasswords(existingUser *domain.User, userInfo *domain.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userInfo.Password))
	if err != nil {
		return errors.New("passwords don't match")
	}
	return nil
}