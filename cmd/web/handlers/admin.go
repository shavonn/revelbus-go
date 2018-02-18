package handlers

import (
	"net/http"
	"revelforce/cmd/web/view"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "admin-dashboard", &view.View{})
}
