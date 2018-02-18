package handlers

import (
	"html/template"
	"net/http"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	trips, err := models.GetUpcomingTrips(3)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	slides, err := models.GetActiveSlides()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "home", &view.View{
		ActiveKey: "home",
		Trips:     trips,
		Slides:    slides,
	})
}

func About(w http.ResponseWriter, r *http.Request) {
	s := models.Settings{
		ID: 1,
	}

	err := s.Get()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "about", &view.View{
		ActiveKey: "about",
		Blurb:     s.AboutBlurb,
		Content:   template.HTML(s.AboutContent),
	})
}

func Contact(w http.ResponseWriter, r *http.Request) {
	s := models.Settings{
		ID: 1,
	}

	err := s.Get()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "contact", &view.View{
		ActiveKey: "contact",
		Blurb:     s.ContactBlurb,
	})
}

func ContactPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.ContactForm{
		Name:    r.PostForm.Get("name"),
		Phone:   r.PostForm.Get("phone"),
		Email:   r.PostForm.Get("email"),
		Message: r.PostForm.Get("message"),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		view.Render(w, r, "contact", v)
	}

	err = email.ContactEmail(f)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, "Your message has been sent!", "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/contact", http.StatusSeeOther)
}

func Trips(w http.ResponseWriter, r *http.Request) {
	trips, err := models.GetUpcomingTrips(0)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "trips", &view.View{
		ActiveKey: "trips",
		Title:     "Upcoming Trips",
		Trips:     trips,
	})
}

func Trip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	t, err := models.GetBySlug(slug)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	trips, err := models.GetUpcomingTrips(2)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "trip", &view.View{
		ActiveKey: "trip",
		Title:     t.Title,
		Trip:      t,
		Trips:     trips,
		Content:   template.HTML(t.Description),
	})
}

func Faq(w http.ResponseWriter, r *http.Request) {
	faqs, err := models.GetActiveFAQs()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	trips, err := models.GetUpcomingTrips(2)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "faq", &view.View{
		ActiveKey:  "faq",
		Title:      "FAQ",
		FAQGrouped: faqs,
		Trips:      trips,
	})
}
