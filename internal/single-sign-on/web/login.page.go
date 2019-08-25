package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"github.com/sirupsen/logrus"
	"gopkg.in/oauth2.v3/manage"
	"net/http"
	"net/url"
	"skripsi-sso/database"
	"skripsi-sso/database/repositories/user"
	"skripsi-sso/utils"
	"strings"
)

type formdata struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
	Remember bool   `form:"remember"`
}

var cache *utils.CacheManager
var manager *manage.Manager

func LoginPage(r *gin.Engine, cm *utils.CacheManager, m *manage.Manager) {
	cache = cm
	manager = m
	router := r.Group("/auth")
	{
		router.GET("/login", login)
		router.POST("/login", login)
		router.GET("/logout", logout)
	}
}

func login(c *gin.Context) {
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	client_id := c.Query("client_id")
	return_to := c.Query("return_to")
	//if return_to == "" || client_id == "" {
	//	//return_to = "/auth/login"
	//	c.JSON(http.StatusInternalServerError, gin.H{"message": "Page not found"})
	//	return
	//}

	_, ok := store.Get("x-user-login")
	if ok {
		//c.Redirect(http.StatusFound, "http://localhost:9999")
		c.JSON(http.StatusOK, gin.H{"message": "page not found"})
		return
	}

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	if c.Request.Method == http.MethodPost {
		var param formdata
		err := c.ShouldBind(&param)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": err.Error()})
			return
		}

		if return_to != "" || client_id != "" {
			return_to, _ = url.QueryUnescape(return_to)
			urlParse, errParse := url.Parse(return_to)
			if errParse != nil {
				c.HTML(http.StatusOK, "login.html", gin.H{"error": true, "message": "parameter return_to failed"})
				return
			}

			urlParse, errParse = url.Parse(urlParse.Query().Get("redirect_uri"))
			if errParse != nil {
				c.HTML(http.StatusOK, "login.html", gin.H{"error": true, "message": "parameter return_to failed"})
				return
			}

			//domain := urlParse.Scheme + "://" + urlParse.Host
		}

		repo := user.NewUserRepo(database.DB)
		userAuth, _ := repo.FindUserByEmailForAuthSso(param.Email)
		if userAuth == nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": true, "message": "email or password is incorrect"})
			return
		}

		if utils.ComparePassword(userAuth.Password, param.Password) {

			store.Set("x-user-login", userAuth.Email)
			_ = store.Save()
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": true, "message": "email or password is incorrect"})
			return
		}

		if return_to == "" {
			if userAuth.RoleId == 1 {
				return_to = "/oauth2client"
			} else {
				return_to = "/profile"
			}
		}

		logrus.Info(return_to)
		c.Redirect(http.StatusFound, return_to)
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
	return
}

func logout(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != "" {
		token = strings.Replace(token, "Bearer ", "", -1)
		logrus.WithField("tag", "sso").Info("logout token : ", token)

		errRemove := manager.RemoveAccessToken(token)
		if errRemove != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": errRemove.Error()})
			return
		}
	} else {
		err := session.Destroy(nil, c.Writer, c.Request)
		if err != nil {
			logrus.WithField("tag", "sso").Info("error : ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.Redirect(http.StatusFound, "/auth/login")
}
