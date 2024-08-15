package infrastructure

import (
	"errors"
	"fmt"
	"golang-clean-architecture/domain"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userInfo *domain.User) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" : userInfo.Email,
		"role" : userInfo.Role,
		"exp" : time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error while generating token")
	}

	return token, nil
}