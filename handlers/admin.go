package handlers

import (
	"net/http"
	"revelforce-admin/internal/platform/forms"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "dashboard.html", &view{})
}

func tourForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, "tour.html", &view{
		Form: new(forms.TourForm),
	})
}
