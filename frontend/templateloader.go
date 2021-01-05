package frontend

import (
	"html/template"
	"log"
	"os"
	"path"
	"strings"
)

// Funcmap
var funcs = template.FuncMap{
	"Upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"Lower": func(s string) string {
		return strings.ToLower(s)
	},
}

func New() {
	name := path.Base("./templates/index.gohtml")
	log.Println(name)
	tpl, err := template.New(name).
		Funcs(funcs).
		ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/index.gohtml")
	if err != nil {
		log.Fatalln("Fatal:", err)
		return
	}
	data := map[string]string{
		"Title": "Das ist ein Titel!",
	}
	log.Println("Executing template with data:", data)
	log.Println("Result:", tpl.Execute(os.Stdout, data))
}
