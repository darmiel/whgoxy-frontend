package web

import (
	"fmt"
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/authenticator"
	"github.com/darmiel/whgoxy-frontend/internal/whgoxy/config"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type WebServer struct {
	parser    *TemplateParser
	router    *mux.Router
	templates map[string]*template.Template
	conf      config.WebConfig
}

func NewWebServer(conf config.WebConfig, parser *TemplateParser, router *mux.Router) (ws *WebServer) {
	return &WebServer{
		router:    router,
		parser:    parser,
		templates: make(map[string]*template.Template),
		conf:      conf,
	}
}

func (ws *WebServer) readTemplates() {
	ws.templates["home"] = ws.parser.MustParseTemplate("home")
	ws.templates["err_404"] = ws.parser.MustParseTemplate("err_404")
}

func (ws *WebServer) updateRoutes() {
	router := ws.router

	// static dir
	staticDir := ws.conf.WebDir + "/static/"
	prefix := http.StripPrefix("/static", http.FileServer(http.Dir(staticDir)))
	router.PathPrefix("/static/").Handler(prefix)

	// routes
	router.HandleFunc("/", ws.homeRouteHandler)

	// 404
	router.NotFoundHandler = http.HandlerFunc(ws.error404)
}

func (ws *WebServer) Run() (err error) {
	ws.readTemplates()
	ws.updateRoutes()

	return http.ListenAndServe(ws.conf.Addr, ws.router)
}

func (ws *WebServer) Exec(name string, r *http.Request, w http.ResponseWriter, data map[string]interface{}) (err error) {
	if data == nil {
		data = make(map[string]interface{})
	}

	// add default data
	data["CurrentURL"] = r.URL.String()
	if user, ok := authenticator.GetUser(r); ok && user != nil {
		data["User"] = user.DiscordUser
		log.Println("OK user found:", user, ok)
	} else {
		// // debug user
		// // TODO: Remove me later
		// data["User"] = &discord.DiscordUser{
		// 	UserID:        "150347348088848384",
		// 	Username:      "d2a",
		// 	Avatar:        "408d6f884febd122f5252e2fc6d93c2e",
		// 	Discriminator: "1325",
		// 	PublicFlags:   256,
		// 	Flags:         256,
		// 	Locale:        "en-US",
		// 	MFAEnabled:    true,
		// }
		log.Println("ERR user not found:", user, ok)
	}

	// get template
	tpl, ok := ws.templates[name]
	if !ok {
		w.WriteHeader(404)
		_, _ = fmt.Fprint(w, "Template "+name+" not found.")
		return
	}
	return tpl.Execute(w, data)
}

func (ws *WebServer) MustExec(name string, w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	panic(ws.Exec(name, r, w, data))
}
