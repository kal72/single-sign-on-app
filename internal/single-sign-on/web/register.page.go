package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"skripsi-sso/database"
	"skripsi-sso/database/entities"
	"skripsi-sso/database/repositories/user"
	"skripsi-sso/utils"
)

func RegisterPage(r *gin.Engine) {
	r.GET("/register", register)
	r.POST("/register", register)
}

func register(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var param entities.User
		err := c.ShouldBind(&param)
		if err != nil {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": err.Error()})
			return
		}

		if param.Password != param.ConfirmPassword {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": "Password does not match!"})
			return
		}

		userRepo := user.NewUserRepo(database.DB)
		param.RoleId = 2
		param.Password = utils.GenerateHashAndSalt(param.Password)
		err = userRepo.SaveUser(&param)
		if err != nil {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": err.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/auth/login")
		return
	}

	c.HTML(http.StatusOK, "register.html", nil)
	return
}
