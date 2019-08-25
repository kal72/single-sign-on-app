package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"os"
	"skripsi-sso/config"
	"skripsi-sso/database"
	sso "skripsi-sso/internal/single-sign-on"
	"skripsi-sso/internal/single-sign-on/api"
	"skripsi-sso/middleware"
	"skripsi-sso/utils"
)

var Env = "heroku"

func main() {
	var port string
	arguments := os.Args
	if len(arguments) > 1 {
		port = arguments[1]
	}

	config.Init(Env)
	database.Init()
	cacheManager := utils.InitCache()

	router := gin.New()
	router.Use(ginlogrus.Logger(logrus.StandardLogger()), gin.Recovery())
	router.Use(middleware.Cors())

	oauthManager := sso.OauthServerInit()
	sso.Init(router, cacheManager, oauthManager)
	v1 := router.Group("/api/v1")
	{
		api.TokenRouter(v1)
		api.Oauth2ClientRouter(v1)
		api.UserInfoRouter(oauthManager, v1)
	}

	router.Run(":" + port)
}
