package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	trips, err := db.GetUpcomingTrips(3)
	if err != nil {
		serverError(w, r, err)
		return
	}

	slides, err := db.GetActiveSlides()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "home", &view{
		ActiveKey: "home",
		Trips:     trips,
		Slides:    slides,
	})
}

func contact(w http.ResponseWriter, r *http.Request) {
	render(w, r, "contact", &view{
		ActiveKey: "contact",
	})
}

func trips(w http.ResponseWriter, r *http.Request) {
	trips, err := db.GetUpcomingTrips(0)
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "trips", &view{
		ActiveKey: "trips",
		Title:     "Upcoming Trips",
		Trips:     trips,
	})
}

func trip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	t, err := db.GetBySlug(slug)
	if err != nil {
		serverError(w, r, err)
		return
	}

	trips, err := db.GetUpcomingTrips(2)
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "trip", &view{
		ActiveKey: "trip",
		Title:     t.Title,
		Trip:      t,
		Trips:     trips,
	})
}

func faq(w http.ResponseWriter, r *http.Request) {
	faqs, err := db.GetActiveFAQs()
	if err != nil {
		serverError(w, r, err)
		return
	}

	trips, err := db.GetUpcomingTrips(2)
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "faq", &view{
		ActiveKey:  "faq",
		Title:      "FAQ",
		FAQGrouped: faqs,
		Trips:      trips,
	})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "dashboard", &view{})
}
