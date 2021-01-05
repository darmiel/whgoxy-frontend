package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type WebConfig struct {
	Dir  string
	Addr string
}

// "./web/static"
func NewWebConfig(dir string) (cfg *WebConfig) {
	return &WebConfig{
		Dir: dir,
	}
}

type WebServer struct {
	parser    *TemplateParser
	router    *mux.Router
	templates map[string]*template.Template
	conf      *WebConfig
}

func NewWebServer(conf *WebConfig, parser *TemplateParser, router *mux.Router) (ws *WebServer) {
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
	staticDir := ws.conf.Dir + "/static/"
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
	data["CurrentURL"] = r.URL.RequestURI()
	log.Println("Data:", data)

	// get template
	tpl, ok := ws.templates[name]
	if !ok {
		w.WriteHeader(404)
		_, _ = fmt.Fprint(w, "Template "+name+" not found.")
		return
	}

	// execut template
	return tpl.Execute(w, map[string]interface{}{
		"CurrentURL": "abc",
	})
}

func (ws *WebServer) MustExec(name string, w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	panic(ws.Exec(name, r, w, data))
}
