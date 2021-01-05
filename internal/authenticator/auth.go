package authenticator

import (
	"github.com/darmiel/whgoxy-frontend/internal/discord"
	"github.com/dchest/authcookie"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	DiscordUser *discord.DiscordUser
	Token       *oauth2.Token
}

var (
	authenticatedUsers = make(map[string]*User)
)

var (
	cookieHost   = os.Getenv("COOKIE_HOST")
	cookieSecret = []byte(os.Getenv("COOKIE_SECRET"))
	cookieName   = "whgoxy-authenticated"
)

func GetLoginCookie(r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(cookieName)
	log.Println("r.Cookie Result:", cookie, err)
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
	// generate cookie
	cookie := authcookie.NewSinceNow(u.DiscordUser.UserID, 8*time.Hour, cookieSecret)
	log.Println("Generated cookie:", cookie)

	// add cookie
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   cookie,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// add to map
	authenticatedUsers[u.DiscordUser.UserID] = u
}
