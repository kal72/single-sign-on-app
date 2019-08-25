package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"net/http"
	"skripsi-sso/database"
	"skripsi-sso/database/repositories/user"
)

func ProfilePage(r *gin.Engine) {
	r.GET("/profile", profile)
}

func profile(c *gin.Context) {
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	data, ok := store.Get("x-user-login")
	if !ok {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}

	userRepo := user.NewUserRepo(database.DB)
	user, err := userRepo.FindUserByEmailForAuthSso(data.(string))

	c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})
	return
}
