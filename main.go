package main

import (
	"github.com/darmiel/whgoxy-frontend/internal/authenticator"
	"github.com/darmiel/whgoxy-frontend/internal/web"
	"github.com/gorilla/mux"
)

func main() {

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	_, _ = fmt.Fprint(w, `<html><body><a href="/login">Login</a></body></html>`)
	// })

	// router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
	// 	if u, die := authenticator.GetUserOrDie(r, w); die {
	// 		return
	// 	} else {
	// 		_, _ = fmt.Fprintf(w, "You are logged in! (%s)", u.DiscordUser.Username)
	// 	}
	// })

	// auth
	router := mux.NewRouter()
	authenticator.New(router, "/dashboard")
	parser := web.NewTemplateParser()

	ws := web.NewWebServer(&web.WebConfig{
		Dir:  "./web",
		Addr: ":1337",
	}, parser, router)

	panic(ws.Run())
}
