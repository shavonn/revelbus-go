package handlers

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	render(w, r, "home", &view{})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "dashboard", &view{})
}
