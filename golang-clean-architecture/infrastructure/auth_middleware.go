package infrastructure

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang-clean-architecture/domain"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "authorization header not found"})
			c.Abort()
			return 
		}
		headerSlice := strings.Split(authHeader, " ")
		if len(headerSlice) != 2 ||	strings.ToLower(headerSlice[0]) != "bearer" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "bearer token not found"})
			c.Abort()
			return
		}

		token, err :=  jwt.Parse(headerSlice[1], func(token *jwt.Token) (interface{}, error) {
			if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("incompatible tokenization method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})	
		
		if err != nil || !token.Valid {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error" : "internal server error"})
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "invalid token"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error" : "invalid token"})
			c.Abort()
			return
		}

		c.Set("AuthorizedUser", &domain.AuthenticatedUser{
			Role : role,
			Email : email,
		})
		c.Next()
	}
}