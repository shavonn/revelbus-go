package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strconv"
)

func settingsForm(w http.ResponseWriter, r *http.Request) {
	s := &db.Settings{
		ID: 1,
	}

	err := s.Get()
	if err == db.ErrNotFound {
		render(w, r, "settings", &view{
			Form:  new(forms.SettingsForm),
			Title: "Settings",
		})
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.SettingsForm{
		ID:           strconv.Itoa(s.ID),
		ContactBlurb: s.ContactBlurb,
		AboutBlurb:   s.AboutBlurb,
		AboutContent: s.AboutContent,
	}

	render(w, r, "settings", &view{
		Title: "Settings",
		Form:  f,
	})
}

func postSettings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.SettingsForm{
		ID:           r.PostForm.Get("id"),
		ContactBlurb: r.PostForm.Get("contact_blurb"),
		AboutBlurb:   r.PostForm.Get("about_blurb"),
		AboutContent: r.PostForm.Get("about_content"),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		render(w, r, "settings", v)
	}

	var msg string

	s := db.Settings{
		ID:           toInt(f.ID),
		ContactBlurb: f.ContactBlurb,
		AboutBlurb:   f.AboutBlurb,
		AboutContent: f.AboutContent,
	}

	if f.ID != "" {
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
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
}
