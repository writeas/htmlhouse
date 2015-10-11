package htmlhouse

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/writeas/impart"
	"github.com/writeas/nerds/store"
)

func createHouse(app *app, w http.ResponseWriter, r *http.Request) error {
	html := r.FormValue("html")
	if strings.TrimSpace(html) == "" {
		return impart.HTTPError{http.StatusBadRequest, "Supply something to publish."}
	}

	houseID := store.GenerateFriendlyRandomString(8)

	_, err := app.db.Exec("INSERT INTO houses (id, html) VALUES (?, ?)", houseID, html)
	if err != nil {
		return err
	}

	if err = app.session.writeToken(w, houseID); err != nil {
		return err
	}

	resUser := newSessionInfo(houseID)

	return impart.WriteSuccess(w, resUser, http.StatusCreated)
}

func getHouse(app *app, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	houseID := vars["house"]

	// Fetch HTML
	var html string
	err := app.db.QueryRow("SELECT html FROM houses WHERE id = ?", houseID).Scan(&html)
	switch {
	case err == sql.ErrNoRows:
		return impart.HTTPError{http.StatusNotFound, "Return to sender. Address unknown."}
	case err != nil:
		fmt.Printf("Couldn't fetch: %v\n", err)
		return err
	}

	// Print HTML
	fmt.Fprintf(w, "%s", html)
	return nil
}
