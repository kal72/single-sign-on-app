package api

import (
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
)

func TokenRouter(r *gin.RouterGroup) {
	router := r.Group("/oauth2")
	{
		router.GET("/token", exchangeToken)
		router.POST("/token", exchangeToken)
	}
}

func exchangeToken(c *gin.Context) {
	ginserver.HandleTokenRequest(c)
}
