package htmlhouse

import (
	"fmt"
	"html/template"
)

const (
	templatesDir = "templates/"
)

func initTemplate(app *app, name string) {
	fmt.Printf("Loading %s%s.html\n", templatesDir, name)

	app.templates[name] = template.Must(template.New("").ParseFiles(templatesDir + name + ".html"))
}

func (app *app) initTemplates() {
	app.templates = map[string]*template.Template{}

	// Initialize dynamic pages
	initTemplate(app, "editor")
	initTemplate(app, "stats")
	initTemplate(app, "browse")
}
