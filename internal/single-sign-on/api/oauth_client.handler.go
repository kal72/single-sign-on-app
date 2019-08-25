package api

import (
	"skripsi-sso/database"
	"skripsi-sso/internal/single-sign-on/helper"
	"skripsi-sso/internal/single-sign-on/model"
	"skripsi-sso/utils/builder"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Oauth2ClientRouter(r *gin.RouterGroup) {
	router := r.Group("/oauth2")
	{
		router.GET("/client", findOauthClient)
		router.POST("/client/register", newGenerateOauthClient)
	}
}

func newGenerateOauthClient(c *gin.Context) {
	logOauth := logrus.WithField("tag", "oauth2")
	var param model.Oauth2Clients
	err := c.BindJSON(&param)
	if err != nil {
		logOauth.Error(err.Error())
		c.JSON(http.StatusBadRequest, builder.Response(false, err.Error(), nil))
		return
	}

	param.ID = helper.GenerateClientId()
	param.Secret = helper.GenerateClientSecret()
	data, errJson := json.Marshal(&param)
	if errJson != nil {
		logOauth.Error(errJson.Error())
		c.JSON(http.StatusBadRequest, builder.Response(false, err.Error(), nil))
		return
	}

	param.Data = data
	errSave := database.DB[0].Create(&param).Error
	if errSave != nil {
		logOauth.Error(errSave.Error())
		c.JSON(http.StatusInternalServerError, builder.Response(false, "error save", nil))
		return
	}

	param.Data = nil
	c.JSON(http.StatusOK, builder.Response(true, "ok", param))
}

func findOauthClient(c *gin.Context) {
	logOauth := logrus.WithField("tag", "oauth2")

	var results []model.Oauth2Clients
	errFind := database.DB[0].Find(&results).Error
	if errFind != nil {
		logOauth.Error(errFind)
		c.JSON(http.StatusOK, builder.Response(false, errFind.Error(), results))
		return
	}

	c.JSON(http.StatusOK, builder.Response(true, "ok", results))
}
