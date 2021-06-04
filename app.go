package htmlhouse

import (
	"database/sql"
	"fmt"
	"github.com/writeas/web-core/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

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

	if !app.cfg.WFMode {
		err = app.initDatabase()
		if err != nil {
			return app, err
		}
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
	if app.cfg.AllowPublish {
		api.HandleFunc("/create", app.handler(createHouse)).Methods("POST").Name("create")
	}
	api.HandleFunc("/{house:[A-Za-z0-9.-]{8}}", app.handler(renovateHouse)).Methods("POST").Name("update")
	api.HandleFunc("/public", app.handler(getPublicHousesData)).Methods("GET").Name("browse-api")

	if app.cfg.WFMode {
		horseAPI := app.router.PathPrefix("/🐴/").Subrouter()
		horseAPI.HandleFunc("/stable", app.handler(redirectToStable))
	}

	admin := app.router.PathPrefix("/admin/").Subrouter()
	admin.HandleFunc("/unpublish", app.handler(banHouse)).Methods("POST").Name("unpublish")
	admin.HandleFunc("/republish", app.handler(unbanHouse)).Methods("POST").Name("republish")

	app.router.HandleFunc("/", app.handler(getEditor)).Methods("GET").Name("index")
	app.router.HandleFunc("/edit/{house:[A-Za-z0-9.-]{8}}.html", app.handler(getEditor)).Methods("GET").Name("edit")
	app.router.HandleFunc("/stats/{house:[A-Za-z0-9.-]{8}}.html", app.handler(viewHouseStats)).Methods("GET").Name("stats")
	if app.cfg.WFMode {
		app.router.HandleFunc("/about.html", app.handler(handleHorseFile))
		app.router.HandleFunc("/contact.html", app.handler(handleHorseFile))
		app.router.PathPrefix("/js/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.PathPrefix("/css/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.PathPrefix("/img/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.PathPrefix("/404.html").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.PathPrefix("/browse").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
		app.router.HandleFunc("/{house:[A-Za-z0-9.\\-/]+}", app.handler(getEditor)).Methods("GET").Name("edit")
	} else {
		app.router.HandleFunc("/{house:[A-Za-z0-9.-]{8}}.html", app.handler(getHouse)).Methods("GET").Name("get")
	}
	app.router.HandleFunc("/browse", app.handler(viewHouses)).Methods("GET").Name("browse")
	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
}

func renderNotFound(app *app, w http.ResponseWriter, isHorse bool) error {
	p := "404.html"
	if isHorse {
		p = "404-css.html"
	}
	page, err := ioutil.ReadFile(filepath.Join(app.cfg.StaticDir, p))
	if err != nil {
		page = []byte("<!DOCTYPE html><html><body>HTMLlot.</body></html>")
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s", page)
	return err
}

type EditorPage struct {
	ID      string
	Content string
	Public  bool

	AllowPublish bool

	SiteURL string
}

func getEditor(app *app, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	house := vars["house"]

	editor := "editor"
	if app.cfg.WFMode {
		editor = "wf-editor"
	}
	if house == "" {
		var defaultPage []byte
		var err error
		if app.cfg.WFMode {
			return handleHorseFile(app, w, r)
		} else {
			defaultPage, err = ioutil.ReadFile(app.cfg.StaticDir + "/default.html")
			if err != nil {
				fmt.Printf("\n%s\n", err)
				defaultPage = []byte("<!DOCTYPE html>\n<html>\n</html>")
			}
		}
		app.templates[editor].ExecuteTemplate(w, "editor", &EditorPage{
			Content:      string(defaultPage),
			AllowPublish: app.cfg.AllowPublish,
		})

		return nil
	}

	d := &EditorPage{
		ID:           house,
		AllowPublish: app.cfg.AllowPublish,
	}

	var err error
	if app.cfg.WFMode {
		if strings.IndexRune(house, '.') == -1 {
			return renderNotFound(app, w, true)
		}
		u, err := url.Parse("https://" + house)
		if err != nil {
			fmt.Fprintf(w, "Unable to parse '%s': %s", house, err)
			return err
		}
		log.Info("Fetching %+v", u)
		resp, err := http.Get(u.String())
		if err != nil {
			fmt.Fprintf(w, "Unable to fetch '%s': %s", u.String(), err)
			return err
		}
		defer resp.Body.Close()

		htmlBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		d.Content = string(htmlBytes)
		// Change internal links to stay on designer
		d.Content = strings.Replace(d.Content, `<a href="`+u.String(), `<a href="/`+u.Host+u.Path, -1)
		d.Content = strings.Replace(d.Content, `href="/css/write.css"`, `href="`+u.Scheme+`://`+u.Host+`/css/write.css"`, -1)
		d.Content = strings.Replace(d.Content, "<a ", `<a target="_parent" `, -1)
		d.SiteURL = u.String()
	} else {
		d.Content, err = getHouseHTML(app, house)
		if err != nil {
			return err
		}
		d.Public = isHousePublic(app, house)
	}

	app.templates[editor].ExecuteTemplate(w, "editor", d)
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
		if err.Status > 300 && err.Status <= 399 {
			w.Header().Set("Location", err.Message)
			w.WriteHeader(err.Status)
			return
		}
		impart.WriteError(w, err)
		return
	}

	impart.WriteError(w, impart.HTTPError{http.StatusInternalServerError, "This is an unhelpful error message for a miscellaneous internal error."})
}
