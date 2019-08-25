package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"OPTIONS", "POST", "PUT", "GET", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Accept", "Accept-Language", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "x-token"},
		MaxAge:           24 * time.Hour,
	})
}
