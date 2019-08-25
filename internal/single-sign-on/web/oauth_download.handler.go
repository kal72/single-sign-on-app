package web

import (
	"skripsi-sso/database"
	model2 "skripsi-sso/internal/single-sign-on/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Oauth2DownloadRouter(r *gin.Engine) {
	router := r.Group("/oauth2")
	{
		router.GET("/client/download", clientOauthDownloadJson)
	}
}

func clientOauthDownloadJson(c *gin.Context) {
	clientID := c.Query("client_id")
	clientSecret := c.Query("client_secret")
	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "client_id or client_secret cannot be empty"})
		return
	}

	var oauth2Clients model2.Oauth2Clients
	err := database.DB[0].Where("id = ? AND secret = ?", clientID, clientSecret).First(&oauth2Clients).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	//jByte, _ := json.Marshal(oauth2Clients)
	c.Header("Content-Disposition", "attachment; filename=client.credentials.json")
	c.Data(http.StatusOK, "application/octet-stream", oauth2Clients.Data)
}
