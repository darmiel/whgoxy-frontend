package main

import (
	"fmt"
	"github.com/darmiel/whgoxy-frontend/authenticator"
	"github.com/darmiel/whgoxy-frontend/frontend"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	_, _ = fmt.Fprint(w, `<html><body><a href="/login">Login</a></body></html>`)
	// })

	router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if u, die := authenticator.GetUserOrDie(r, w); die {
			return
		} else {
			_, _ = fmt.Fprintf(w, "You are logged in! (%s)", u.DiscordUser.Username)
		}
	})

	// auth
	authenticator.New(router, "/dashboard")

	//// Template
	// static content
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	frontend.New()
	////

	if err := http.ListenAndServe(":1337", router); err != nil {
		log.Fatalln("Error serving:", err.Error())
	}
}
