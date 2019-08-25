package helper

import (
	"skripsi-sso/middleware"
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func GenerateClientId() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(bytes)
}

func GenerateClientSecret() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func StoreSession(token string, c *gin.Context) {
	session := sessions.Default(c)
	session.Set("s-token", token)
	session.Save()
}

func CheckSession(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	sToken := session.Get("s-token")
	if sToken != nil {
		_, err := middleware.Authorize("Bearer " + sToken.(string))
		if err != nil {
			session.Clear()
			return "", false
		}

		return sToken.(string), true
	}

	return "", false
}
