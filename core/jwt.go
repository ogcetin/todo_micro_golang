package core

import (
	"fmt"
	"strings"
	"time"
	"todomicro/helpers"
	"todomicro/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

var signKey = []byte("gunardindanbatarkentepelerinvedavaktidirteletabilerin")

type JWT struct {
	Phone    string
	Password string
}

var loggedUser models.User

func GetLoggedUser() models.User {
	return loggedUser
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}
		err := VerifyToken(tokenString)

		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CreateToken(j JWT) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone":    j.Phone,
		"password": j.Password,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(headerString string) error {
	// TODO check Bearer keyword

	authorizationParts := strings.Split(headerString, " ")
	token, err := jwt.Parse(authorizationParts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return signKey, nil
	})

	if err != nil {
		return err
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	var phone = claims["phone"].(string)
	var pass = claims["password"].(string)

	loggedUser, err = helpers.UserCheck(phone, pass)
	if err != nil {
		return err
	}

	return nil
}
