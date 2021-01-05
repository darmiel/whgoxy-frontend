package web

import (
	"fmt"
	"net/http"
)

func (ws *WebServer) homeRouteHandler(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("home", writer, request, nil)
}

func (ws *WebServer) error404(writer http.ResponseWriter, request *http.Request) {
	ws.MustExec("err_404", writer, request, nil)
}

func (*WebServer) redirectHttps(w http.ResponseWriter, req *http.Request) {
	target := fmt.Sprintf("https://%s%s", req.Host, req.URL.Path)
	http.Redirect(w, req, target, http.StatusPermanentRedirect)
}
