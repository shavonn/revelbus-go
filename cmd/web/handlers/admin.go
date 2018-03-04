package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/domain/models"
	"revelforce/internal/platform/flash"
	"strconv"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "admin-dashboard", &view.View{})
}

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	s := &models.Settings{
		ID: 1,
	}

	err := s.Fetch()
	if err != nil {
		if err == domain.ErrNotFound {
			view.Render(w, r, "settings", &view.View{
				Form:  new(models.SettingsForm),
				Title: "Settings",
			})
			return
		}
		view.ServerError(w, r, err)
		return
	}

	galleries, err := models.FetchGalleries()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.SettingsForm{
		ID:                strconv.Itoa(s.ID),
		ContactBlurb:      s.ContactBlurb.String,
		AboutBlurb:        s.AboutBlurb.String,
		AboutContent:      s.AboutContent.String,
		HomeGalleryID:     int(s.HomeGalleryID.Int64),
		HomeGalleryActive: s.HomeGalleryActive,
	}

	view.Render(w, r, "settings", &view.View{
		Title:     "Settings",
		Form:      f,
		Galleries: galleries,
	})
}

func PostSettings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.SettingsForm{
		ID:                r.PostForm.Get("id"),
		ContactBlurb:      r.PostForm.Get("contact_blurb"),
		AboutBlurb:        r.PostForm.Get("about_blurb"),
		AboutContent:      r.PostForm.Get("about_content"),
		HomeGalleryID:     utils.ToInt(r.PostForm.Get("home_gallery")),
		HomeGalleryActive: (len(r.Form["home_gallery_active"]) == 1),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		view.Render(w, r, "settings", v)
	}

	var msg string

	s := models.Settings{
		ID:                utils.ToInt(f.ID),
		ContactBlurb:      utils.NewNullStr(f.ContactBlurb),
		AboutBlurb:        utils.NewNullStr(f.AboutBlurb),
		AboutContent:      utils.NewNullStr(f.AboutContent),
		HomeGalleryID:     utils.NewNullInt(strconv.Itoa(f.HomeGalleryID)),
		HomeGalleryActive: f.HomeGalleryActive,
	}

	if f.ID != "" {
		err := s.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		msg = utils.MsgSuccessfullyUpdated
	} else {
		err := s.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/settings", http.StatusSeeOther)
}
