package main

import (
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/authenticator"
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/config"
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/web"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	// config
	if err := config.Load(); err != nil {
		log.Fatalln("Fatal:", err.Error())
		return
	}

	// auth
	router := mux.NewRouter()
	authenticator.New(router, config.ConfigOAuth)
	parser := web.NewTemplateParser()

	ws := web.NewWebServer(config.ConfigWeb, parser, router)

	if err := ws.Run(); err != nil {
		panic(err)
	}
}
