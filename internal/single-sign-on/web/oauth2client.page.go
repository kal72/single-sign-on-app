package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"skripsi-sso/database"
	"skripsi-sso/internal/single-sign-on/helper"
	"skripsi-sso/internal/single-sign-on/model"
	"skripsi-sso/utils/builder"
)

func RegisterOauthPage(r *gin.Engine) {
	r.GET("/oauth2client", registerOauth)
	r.POST("/oauth2client", registerOauth)
}

func registerOauth(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var param model.Oauth2Clients
		err := c.ShouldBind(&param)
		if err != nil {
			c.HTML(http.StatusOK, "oauth2client.html", gin.H{"error": err.Error()})
			return
		}

		param.ID = helper.GenerateClientId()
		param.Secret = helper.GenerateClientSecret()
		data, errJson := json.Marshal(&param)
		if errJson != nil {
			logrus.Error(errJson.Error())
			c.JSON(http.StatusBadRequest, builder.Response(false, err.Error(), nil))
			return
		}

		param.Data = data
		errSave := database.DB[0].Create(&param).Error
		if errSave != nil {
			logrus.Error(errSave.Error())
			c.JSON(http.StatusInternalServerError, builder.Response(false, "error save", nil))
			return
		}

		c.HTML(http.StatusOK, "oauth2clientgenerate.html", gin.H{"client": param})
		return
	}

	oauthclient := model.Oauth2Clients{}
	c.HTML(http.StatusOK, "oauth2client.html", gin.H{"client": oauthclient})
	return
}
