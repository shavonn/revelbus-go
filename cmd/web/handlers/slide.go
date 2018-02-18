package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strconv"

	"github.com/gorilla/mux"
)

func slideForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		render(w, r, "slide", &view{
			Form:  new(forms.SlideForm),
			Title: "New Slide",
		})
		return
	}

	s := &db.Slide{
		ID: toInt(id),
	}

	err := s.Get()
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.SlideForm{
		ID:     strconv.Itoa(s.ID),
		Title:  s.Title,
		Blurb:  s.Blurb,
		Style:  s.Style,
		Order:  strconv.Itoa(s.Order),
		Active: s.Active,
	}

	render(w, r, "slide", &view{
		Title: s.Title,
		Form:  f,
		Slide: s,
	})
}

func postSlide(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.SlideForm{
		ID:     r.PostForm.Get("id"),
		Title:  r.PostForm.Get("title"),
		Blurb:  r.PostForm.Get("blurb"),
		Style:  r.PostForm.Get("style"),
		Order:  r.PostForm.Get("order"),
		Active: (id == "" || ((len(r.Form["active"]) == 1) && id != "")),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		if id == "" {
			v.Title = "New Slide"
		}

		render(w, r, "slide", v)
	}

	var msg string

	s := db.Slide{
		ID:     toInt(f.ID),
		Title:  f.Title,
		Blurb:  f.Blurb,
		Style:  f.Style,
		Order:  toInt(f.Order),
		Active: f.Active,
	}

	if id != "" {
		s.ID = toInt(id)

		err := s.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		msg = MsgSuccessfullyUpdated
	} else {
		err := s.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		id = strconv.Itoa(s.ID)
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/slide?id="+id, http.StatusSeeOther)
}

func listSlides(w http.ResponseWriter, r *http.Request) {
	slides, err := db.GetSlides()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "slides", &view{
		Title:  "Slides",
		Slides: slides,
	})
}

func removeSlide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s := db.Slide{
		ID: toInt(id),
	}

	err := s.Delete()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyRemoved, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/slides", http.StatusSeeOther)
}
