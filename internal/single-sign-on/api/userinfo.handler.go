package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/oauth2.v3/manage"
	"skripsi-sso/database"
	"skripsi-sso/database/repositories/user"
	"strings"
)

var mManager *manage.Manager

func UserInfoRouter(manager *manage.Manager, r *gin.RouterGroup) {
	mManager = manager
	r.GET("/userinfo", userinfo)
}

func userinfo(c *gin.Context) {
	sToken := c.GetHeader("Authorization")
	logrus.Debug(sToken)
	result := strings.Split(sToken, "Bearer ")
	ti, err := mManager.LoadAccessToken(result[1])
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	userRepo := user.NewUserRepo(database.DB)
	userInfo, err := userRepo.FindUserByEmailForAuthSso(ti.GetUserID())
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	logrus.Debug(userInfo)
	c.JSON(200, gin.H{"token": ti, "user": userInfo})
}
