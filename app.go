package htmlhouse

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/writeas/impart"
)

type app struct {
	cfg    *config
	router *mux.Router
}

func newApp() (*app, error) {
	var err error

	app := &app{}

	app.cfg, err = newConfig()
	if err != nil {
		return app, err
	}

	app.initRouter()

	return app, nil
}

func (app *app) initRouter() {
	app.router = mux.NewRouter()

	/*
		api := app.router.PathPrefix("/api/").Subrouter()

		anon := api.PathPrefix("/anon/").Subrouter()
		anon.HandleFunc("/signup", app.handler(signupAnonymous, authLevelNone)).Methods("POST").Name("signupAnon")

		auth := api.PathPrefix("/auth/").Subrouter()
		//	auth.HandleFunc("/", app.handler(login, authLevelNone)).Methods("POST").Name("login")
		auth.HandleFunc("/", app.handler(logout, authLevelUser)).Methods("DELETE").Name("logout")
		auth.HandleFunc("/signup", app.handler(signup, authLevelNone)).Methods("POST").Name("signup")

		users := api.PathPrefix("/users/").Subrouter()
		users.HandleFunc("/{username:[A-Za-z0-9.-]+}", app.handler(userByUsername, authLevelNone)).Methods("GET").Name("user")
	*/

	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.StaticDir)))
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func (app *app) handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			return h(w, r)
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
