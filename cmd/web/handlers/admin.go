package handlers

import (
	"net/http"
)

func adminDashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "admin-dashboard", &view{})
}
