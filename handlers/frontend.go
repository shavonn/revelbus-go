package handlers

import (
	"net/http"
	"revelforce-admin/internal/platform/flash"
)

func index(w http.ResponseWriter, r *http.Request) {

	err := flash.Add(w, r, "it's me", "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "index.html", &view{})
}
