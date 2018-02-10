package handlers

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	render(w, r, "index.html", &view{})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "dashboard.html", &view{})
}
