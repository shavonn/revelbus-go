package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"

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

func about(w http.ResponseWriter, r *http.Request) {
	s := db.Settings{
		ID: 1,
	}

	err := s.Get()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "about", &view{
		ActiveKey: "about",
		Blurb:     s.AboutBlurb,
		Content:   s.AboutContent,
	})
}

func contact(w http.ResponseWriter, r *http.Request) {
	s := db.Settings{
		ID: 1,
	}

	err := s.Get()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "contact", &view{
		ActiveKey: "contact",
		Blurb:     s.ContactBlurb,
	})
}

func contactPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.ContactForm{
		Name:    r.PostForm.Get("name"),
		Phone:   r.PostForm.Get("phone"),
		Email:   r.PostForm.Get("email"),
		Message: r.PostForm.Get("message"),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		render(w, r, "contact", v)
	}

	err = email.ContactEmail(f)
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, "Your message has been sent!", "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/contact", http.StatusSeeOther)
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

func userDashboard(w http.ResponseWriter, r *http.Request) {
	render(w, r, "user-dashboard", &view{})
}
