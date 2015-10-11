package htmlhouse

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/writeas/impart"
)

type app struct {
	cfg     *config
	router  *mux.Router
	db      *sql.DB
	session sessionManager
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

	app.router.HandleFunc("/{house:[A-Za-z0-9.-]+}.html", app.handler(getHouse)).Methods("GET").Name("get")
	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
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
