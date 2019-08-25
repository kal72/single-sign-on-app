package web

import (
	"github.com/go-session/session"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

func UserAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Info("user authorize")
	for name, value := range r.Form {
		logrus.Info(name, " ", value)
	}

	value, ok := store.Get("x-user-login")
	if !ok {
		clientId := r.URL.Query().Get("client_id")
		return_to := r.URL.RequestURI()
		w.Header().Set("Location", "/auth/login?client_id="+clientId+"&return_to="+url.QueryEscape(return_to))
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = value.(string)
	return
}

func ClientScopeHandler(clientID, scope string) (allowed bool, err error) {
	logrus.Info("Client Scope: ", scope)

	valid := true
	scopes := strings.Split(scope, ",")
	if len(scopes) > 3 {
		allowed = false
		return
	}

	for _, s := range scopes {
		if s != "all" && s != "user_info" && s != "api" {
			valid = false
		}
	}

	allowed = valid
	return
}
