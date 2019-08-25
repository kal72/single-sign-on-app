package web

import (
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
)

func AuthorizeHandler(r *gin.Engine) {
	router := r.Group("/oauth2")
	{
		router.GET("/authorize", func(c *gin.Context) {
			ginserver.HandleAuthorizeRequest(c)
		})
	}
}
