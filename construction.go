package htmlhouse

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

func renovateHouse(app *app, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	houseID := vars["house"]
	html := r.FormValue("html")
	if strings.TrimSpace(html) == "" {
		return impart.HTTPError{http.StatusBadRequest, "Supply something to publish."}
	}

	authHouseID, err := app.session.readToken(r)
	if err != nil {
		return err
	}

	if authHouseID != houseID {
		return impart.HTTPError{http.StatusUnauthorized, "Bad token for this ⌂ house ⌂."}
	}

	_, err = app.db.Exec("UPDATE houses SET html = ? WHERE id = ?", html, houseID)
	if err != nil {
		return err
	}

	if err = app.session.writeToken(w, houseID); err != nil {
		return err
	}

	resUser := newSessionInfo(houseID)

	return impart.WriteSuccess(w, resUser, http.StatusOK)
}

func getHouseHTML(app *app, houseID string) (string, error) {
	var html string
	err := app.db.QueryRow("SELECT html FROM houses WHERE id = ?", houseID).Scan(&html)
	switch {
	case err == sql.ErrNoRows:
		return "", impart.HTTPError{http.StatusNotFound, "Return to sender. Address unknown."}
	case err != nil:
		fmt.Printf("Couldn't fetch: %v\n", err)
		return "", err
	}

	return html, nil
}

var (
	htmlReg = regexp.MustCompile("<html( ?.*)>")
)

func getHouse(app *app, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	houseID := vars["house"]

	// Fetch HTML
	html, err := getHouseHTML(app, houseID)
	if err != nil {
		if err, ok := err.(impart.HTTPError); ok {
			if err.Status == http.StatusNotFound {
				page, err := ioutil.ReadFile(app.cfg.StaticDir + "/404.html")
				if err != nil {
					page = []byte("<!DOCTYPE html><html><body>HTMLlot.</body></html>")
				}
				fmt.Fprintf(w, "%s", page)
				return nil
			}
		}
		return err
	}

	// Add nofollow meta tag
	if strings.Index(html, "<head>") == -1 {
		html = htmlReg.ReplaceAllString(html, "<html$1><head></head>")
	}
	html = strings.Replace(html, "<head>", "<head><meta name=\"robots\" content=\"nofollow\" />", 1)

	// Add links back to HTMLhouse
	watermark := "<div style='position: absolute;top:16px;right:16px;'><a href='/'>&lt;&#8962;/&gt;</a> &middot; <a href='/edit/" + houseID + ".html'>edit</a></div>"
	if strings.Index(html, "</body>") == -1 {
		html = strings.Replace(html, "</html>", "</body></html>", 1)
	}
	html = strings.Replace(html, "</body>", fmt.Sprintf("%s</body>", watermark), 1)

	// Print HTML, with sanity check in case someone did something crazy
	if strings.Index(html, "<a href='/'>&lt;&#8962;/&gt;</a>") == -1 {
		fmt.Fprintf(w, "%s%s", html, watermark)
	} else {
		fmt.Fprintf(w, "%s", html)
	}

	if r.Method != "HEAD" && !bots.IsBot(r.UserAgent()) {
		app.db.Exec("UPDATE houses SET view_count = view_count + 1 WHERE id = ?", houseID)
	}

	return nil
}
