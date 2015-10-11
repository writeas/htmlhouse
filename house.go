package htmlhouse

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"log"
)

func Serve() {
	app, err := newApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.close()

	n := negroni.Classic()
	n.UseHandler(app.router)
	n.Run(fmt.Sprintf(":%d", app.cfg.ServerPort))
}
