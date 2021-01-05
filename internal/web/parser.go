package web

import (
	"fmt"
	"html/template"
)

const rootTmpl = `{{ define "root" }} {{ template "base" . }} {{ end }}`

type TemplateParser struct {
}

func NewTemplateParser() (parser *TemplateParser) {
	return &TemplateParser{}
}

func (parser *TemplateParser) ParseTemplate(name string) (tpl *template.Template, err error) {
	root, err := template.New("root").Parse(rootTmpl)
	if err != nil {
		return nil, err
	}

	tmplDir := fmt.Sprintf("%s/template", "web")
	componentsDir := tmplDir + "/components"

	basePath := fmt.Sprintf("%s/base.gohtml", tmplDir)
	tmplPath := fmt.Sprintf("%s/%s.gohtml", tmplDir, name)

	return root.ParseFiles([]string{
		// components
		fmt.Sprintf("%s/head.gohtml", componentsDir),
		fmt.Sprintf("%s/navbar.gohtml", componentsDir),
		fmt.Sprintf("%s/footer.gohtml", componentsDir),

		// base
		basePath,

		// template
		tmplPath,
	}...)
}

// MustParseTemplate calls ParseTemplate(...) and panics on an error.
func (parser *TemplateParser) MustParseTemplate(name string) *template.Template {
	tmpl, err := parser.ParseTemplate(name)
	if err != nil {
		panic(err)
	}

	return tmpl
}
