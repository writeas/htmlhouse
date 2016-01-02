package htmlhouse

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/writeas/impart"
)

type app struct {
	cfg       *config
	router    *mux.Router
	db        *sql.DB
	session   sessionManager
	templates map[string]*template.Template
}

func newApp() (*app, error) {
	var err error

	app := &app{}

	app.cfg, err = newConfig()
	if err != nil {
		return app, err
	}

	app.session, err = newSessionManager(app.cfg)
	if err != nil {
		return app, err
	}

	err = app.initDatabase()
	if err != nil {
		return app, err
	}

	app.initTemplates()

	app.initRouter()

	return app, nil
}

func (app *app) close() {
	app.db.Close()
}

func (app *app) initRouter() {
	app.router = mux.NewRouter()

	api := app.router.PathPrefix("/⌂/").Subrouter()
	api.HandleFunc("/create", app.handler(createHouse)).Methods("POST").Name("create")
	api.HandleFunc("/{house:[A-Za-z0-9.-]{8}}", app.handler(renovateHouse)).Methods("POST").Name("update")

	admin := app.router.PathPrefix("/admin/").Subrouter()
	admin.HandleFunc("/unpublish", app.handler(banHouse)).Methods("POST").Name("unpublish")

	app.router.HandleFunc("/", app.handler(getEditor)).Methods("GET").Name("index")
	app.router.HandleFunc("/edit/{house:[A-Za-z0-9.-]{8}}.html", app.handler(getEditor)).Methods("GET").Name("edit")
	app.router.HandleFunc("/stats/{house:[A-Za-z0-9.-]{8}}.html", app.handler(viewHouseStats)).Methods("GET").Name("stats")
	app.router.HandleFunc("/{house:[A-Za-z0-9.-]{8}}.html", app.handler(getHouse)).Methods("GET").Name("get")
	app.router.HandleFunc("/browse", app.handler(viewHouses)).Methods("GET").Name("browse")
	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
}

type EditorPage struct {
	ID      string
	Content string
	Public  bool
}

func getEditor(app *app, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	house := vars["house"]

	if house == "" {
		defaultPage, err := ioutil.ReadFile(app.cfg.StaticDir + "/default.html")
		if err != nil {
			fmt.Printf("\n%s\n", err)
			defaultPage = []byte("<!DOCTYPE html>\n<html>\n</html>")
		}
		app.templates["editor"].ExecuteTemplate(w, "editor", &EditorPage{"", string(defaultPage), false})

		return nil
	}

	html, err := getHouseHTML(app, house)
	if err != nil {
		return err
	}

	app.templates["editor"].ExecuteTemplate(w, "editor", &EditorPage{house, html, isHousePublic(app, house)})
	return nil
}

type handlerFunc func(app *app, w http.ResponseWriter, r *http.Request) error

func (app *app) handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			return h(app, w, r)
		}())
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if err, ok := err.(impart.HTTPError); ok {
		impart.WriteError(w, err)
		return
	}

	impart.WriteError(w, impart.HTTPError{http.StatusInternalServerError, "This is an unhelpful error message for a miscellaneous internal error."})
}
