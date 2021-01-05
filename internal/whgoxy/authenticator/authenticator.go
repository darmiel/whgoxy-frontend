package authenticator

import (
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/config"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var oauthConfig *oauth2.Config
var cookieSecret []byte

func New(router *mux.Router, conf config.OAuthConfig) {
	cookieSecret = []byte(conf.CookieSecret)

	oauthConfig = &oauth2.Config{
		RedirectURL:  conf.RedirectURL,
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthURL,
			TokenURL: conf.TokenURL,
		},
	}

	router.HandleFunc("/login", handleLoginRoute)
	router.HandleFunc("/callback", handleCallbackRoute)
	router.HandleFunc("/logout", handleLogoutRoute)
}
