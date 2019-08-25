package sso

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	oauth2server "github.com/go-oauth2/gin-server"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"github.com/vgarvardt/go-pg-adapter/pgxadapter"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"net/http"
	"skripsi-sso/config"
	"skripsi-sso/internal/single-sign-on/web"
	"skripsi-sso/utils"
	"time"
)

func OauthServerInit() *manage.Manager {
	logrus.Info("initialize oAuth2 server")
	pgxConnConfig, _ := pgx.ParseURI(config.Databases[0].ConnectionString)
	configPool := pgx.ConnPoolConfig{ConnConfig: pgxConnConfig, MaxConnections: 20}
	poolConn, err := pgx.NewConnPool(configPool)
	if err != nil {
		logrus.Error(err.Error())
	}
	//pgxConn, _ := pgx.Connect(pgxConnConfig)

	tokenCfg := &manage.Config{
		AccessTokenExp:    config.JWT_EXPIRED,
		RefreshTokenExp:   config.JWT_REFRESH_EXPIRED,
		IsGenerateRefresh: true,
	}

	manager := manage.NewDefaultManager()
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(config.App.JwtSecret), jwt.SigningMethodHS512))
	manager.SetAuthorizeCodeTokenCfg(tokenCfg)
	manager.SetAuthorizeCodeExp(time.Minute * 5)
	manager.SetPasswordTokenCfg(tokenCfg)
	manager.SetClientTokenCfg(manage.DefaultClientTokenCfg)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	adapter := pgxadapter.NewConnPool(poolConn)
	clientStore, _ := pg.NewClientStore(adapter)
	tokenStore, _ := pg.NewTokenStore(adapter)
	defer tokenStore.Close()

	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// Initialize the oauth2 service
	oauth2server.InitServer(manager)
	oauth2server.SetAllowGetAccessRequest(true)
	oauth2server.SetClientInfoHandler(server.ClientFormHandler)
	oauth2server.SetUserAuthorizationHandler(web.UserAuthorizeHandler)
	oauth2server.SetAllowedGrantType(oauth2.AuthorizationCode, oauth2.ClientCredentials, oauth2.Refreshing)
	oauth2server.SetClientScopeHandler(web.ClientScopeHandler)

	//insert default client oauth2
	adapter.Exec(`insert into oauth2_clients values('clientqwerty1234Dashboard', 'DONT_DELETE_THIS', 'DONT_DELETE_THIS', '{"secret":"DONT_DELETE_THIS"}') on conflict (id) do nothing`)

	return manager
}

func Init(r *gin.Engine, cm *utils.CacheManager, manager *manage.Manager) {
	r.Use(static.Serve("/resources", static.LocalFile("./internal/single-sign-on/static/resources", false)))
	r.LoadHTMLFiles(
		"./internal/single-sign-on/static/login.html",
		"./internal/single-sign-on/static/auth.html",
		"./internal/single-sign-on/static/register.html",
		"./internal/single-sign-on/static/profile.html",
		"./internal/single-sign-on/static/oauth2client.html",
		"./internal/single-sign-on/static/oauth2clientgenerate.html",
	)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/profile")
	})

	web.Oauth2DownloadRouter(r)
	web.AuthorizeHandler(r)
	web.LoginPage(r, cm, manager)
	web.RegisterPage(r)
	web.ProfilePage(r)
	web.RegisterOauthPage(r)
}
