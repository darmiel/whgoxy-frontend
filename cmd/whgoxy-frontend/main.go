package main

import (
	"context"
	"fmt"
	"github.com/darmiel/whgoxy-frontend/authenticator"
	"github.com/darmiel/whgoxy-frontend/discord"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var (
	oauthConfig *oauth2.Config
)

func init() {
	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1337/callback",
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}

func secret(user, realm string) string {
	if user == "john" {
		// password is "hello"
		return "b98e16cbc3d01734b264adba7baa3bf9"
	}
	return ""
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var html = `
<html>
	<body>
		<a href="/login">Login</a>
	</body>
</html>
`
		_, _ = fmt.Fprint(w, html)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := oauthConfig.AuthCodeURL("test-123")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		token, err := oauthConfig.Exchange(context.TODO(), code)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}
		// store user id,
		// http.Redirect(w, r, "/?token=" + token.AccessToken, http.StatusTemporaryRedirect)

		user, err := discord.NewUserByToken(token)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Invalid token or parse error: %s", err.Error())
			return
		}

		// set cookie
		authenticator.LoginUser(w, &authenticator.User{
			Token:       token,
			DiscordUser: user,
		})

		// redirect
		http.Redirect(w, r, "/user/me", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/user/me", func(w http.ResponseWriter, r *http.Request) {
		user, ok := authenticator.GetUser(r)
		if !ok {
			// _, _ = fmt.Fprint(w, "Not authenticated.")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		html := `
<html>
	<body>
		<!-- Username -->
		<h1>Hallo, %s!</h1>
		<!-- User ID -->
		<small>%s</small>
		<!-- Avatar -->
		<img src="%s" alt="avatar">
	</body>
</html>
`
		_, _ = fmt.Fprintf(w, html, user.DiscordUser.Username, user.DiscordUser.UserID, user.DiscordUser.GetAvatarUrl())
	})

	if err := http.ListenAndServe(":1337", nil); err != nil {
		log.Fatalln("Error serving:", err.Error())
	}
}
