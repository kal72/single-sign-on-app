package middleware

import (
	"skripsi-sso/config"
	"skripsi-sso/database/entities"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"net/http"
	"strings"
	"time"
)

var JwtDuration = time.Hour * 24

func CreateJwtToken(user entities.UserAuth) (string, error) {

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS512)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["authorized"] = true
	claims["user_id"] = user.Email
	claims["exp"] = time.Now().Add(JwtDuration).Unix()

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(config.App.JwtSecret))

	if err != nil {
		logrus.Errorf("Token, Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func Authorize(tokenString string) (string, error) {
	result := strings.Split(tokenString, "Bearer ")
	if len(result) == 1 {
		return "", errors.New("Invalid token")
	}

	tokenString = result[1]
	_, err := jwt.Parse(result[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.App.JwtSecret), nil
	})

	return result[1], err
}

func Verify(tokenString string) (interface{}, error) {
	result := strings.Split(tokenString, "Bearer ")
	if len(result) == 1 {
		return "", errors.New("Invalid token")
	}

	tokenString = result[1]
	data, err := jwt.Parse(result[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.App.JwtSecret), nil
	})

	return data, err
}

func GenerateToken(manager *manage.Manager, userId string, request *http.Request) oauth2.TokenInfo {
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     "clientqwerty1234Dashboard",
		ClientSecret: "DONT_DELETE_THIS",
		UserID:       userId,
		RedirectURI:  "DONT_DELETE_THIS",
		Scope:        "all",
		Request:      request,
	}

	ti, err := manager.GenerateAccessToken(oauth2.PasswordCredentials, tgr)
	if err != nil {
		logrus.WithField("tag", "token").Error(err.Error())
		return nil
	}

	return ti
}
