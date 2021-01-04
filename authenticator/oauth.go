package authenticator

import (
	"context"
	"fmt"
	"github.com/darmiel/whgoxy-frontend/discord"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

func New(dashboardURL string) {
	oauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:1337/callback",
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := oauthConfig.AuthCodeURL("test-123")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// get discord code from query params
		code := r.URL.Query().Get("code")

		// request discord token
		token, err := oauthConfig.Exchange(context.TODO(), code)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		// parse user information
		// (id, username, ...)
		user, err := discord.NewUserByToken(token)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Invalid token or parse error: %s", err.Error())
			return
		}

		// set cookie
		LoginUser(w, &User{
			Token:       token,
			DiscordUser: user,
		})

		// redirect
		http.Redirect(w, r, dashboardURL, http.StatusTemporaryRedirect) // TODO: Change this to dashboard
	})
}
