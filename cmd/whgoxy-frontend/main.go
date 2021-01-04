package main

import (
	"fmt"
	"github.com/darmiel/whgoxy-frontend/authenticator"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `<html><body><a href="/login">Login</a></body></html>`)
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if u, die := authenticator.GetUserOrDie(r, w); die {
			return
		} else {
			_, _ = fmt.Fprintf(w, "You are logged in! (%s)", u.DiscordUser.Username)
		}
	})

	// auth
	authenticator.New("/dashboard")

	if err := http.ListenAndServe(":1337", nil); err != nil {
		log.Fatalln("Error serving:", err.Error())
	}
}
