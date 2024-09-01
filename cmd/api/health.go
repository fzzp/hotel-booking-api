package main

import (
	"net/http"
	"strings"
)

func (app *application) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	app.SUCC(w, r, strings.ToUpper("successfully"))
}
