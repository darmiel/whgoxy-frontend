package authenticator

import (
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/config"
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/discord"
	"github.com/dchest/authcookie"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

type User struct {
	DiscordUser *discord.DiscordUser
	Token       *oauth2.Token
}

var (
	authenticatedUsers = make(map[string]*User)
)

func GetLoginCookie(r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(config.ConfigOAuth.CookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func GetUser(r *http.Request) (u *User, ok bool) {
	// check if user sent a login cookie
	value, ok := GetLoginCookie(r)
	if !ok {
		return nil, false
	}
	log.Println("Checking if cookie is valid...")

	// check if cookie is valid
	if login := authcookie.Login(value, cookieSecret); login != "" {
		u, ok = authenticatedUsers[login]
		return u, ok
	} else {
		return nil, false
	}
}

func GetUserOrDie(r *http.Request, w http.ResponseWriter) (u *User, die bool) {
	u, ok := GetUser(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, true
	}
	return u, false
}

func LoginUser(w http.ResponseWriter, u *User) {
	log.Println("Logging in user", u.DiscordUser.Username, "...")

	// generate cookie
	cookie := authcookie.NewSinceNow(
		u.DiscordUser.UserID,
		8*time.Hour,
		cookieSecret,
	)

	// add cookie
	http.SetCookie(w, &http.Cookie{
		Name:    config.ConfigOAuth.CookieName,
		Value:   cookie,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// add to map
	authenticatedUsers[u.DiscordUser.UserID] = u
}

func LogoutUser(w http.ResponseWriter, u *User) {
	log.Println("Logging out user", u.DiscordUser.Username, "...")

	http.SetCookie(w, &http.Cookie{
		Name:  config.ConfigOAuth.CookieName,
		Value: "",
	})

	delete(authenticatedUsers, u.DiscordUser.UserID)
}
