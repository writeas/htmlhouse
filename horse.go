package htmlhouse

import (
	"fmt"
	"github.com/writeas/impart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func handleHorseFile(app *app, w http.ResponseWriter, r *http.Request) error {
	p := r.URL.Path
	if r.URL.Path == "/" {
		p = "index.html"
	}
	p = strings.Replace(p, ".html", "-css.html", 1)
	p = filepath.Join(app.cfg.StaticDir, p)
	http.ServeFile(w, r, p)
	return nil
}

func redirectToStable(app *app, w http.ResponseWriter, r *http.Request) error {
	u, err := url.Parse(r.FormValue("url"))
	if err != nil {
		fmt.Fprintf(w, "Unable to parse URL: %s", err)
		return nil
	}
	return impart.HTTPError{http.StatusFound, "/" + u.Host + u.Path}
}
