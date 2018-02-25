package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/flash"
	"strconv"

	"github.com/gorilla/mux"
)

func SlideForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "slide", &view.View{
			Form:  new(models.SlideForm),
			Title: "New Slide",
		})
		return
	}

	s := &models.Slide{
		ID: utils.ToInt(id),
	}

	err := s.Fetch()
	if err != nil {
		if err == db.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	f := &models.SlideForm{
		ID:     strconv.Itoa(s.ID),
		Title:  s.Title,
		Blurb:  s.Blurb,
		Style:  s.Style,
		Order:  s.Order,
		Active: s.Active,
	}

	view.Render(w, r, "slide", &view.View{
		Title: s.Title,
		Form:  f,
	})
}

func PostSlide(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.SlideForm{
		ID:     r.PostForm.Get("id"),
		Title:  r.PostForm.Get("title"),
		Blurb:  r.PostForm.Get("blurb"),
		Style:  r.PostForm.Get("style"),
		Order:  utils.ToInt(r.PostForm.Get("order")),
		Active: (len(r.Form["active"]) == 1),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if f.ID == "" {
			v.Title = "New Slide"
		}

		view.Render(w, r, "slide", v)
	}

	var msg string

	s := models.Slide{
		ID:     utils.ToInt(f.ID),
		Title:  f.Title,
		Blurb:  f.Blurb,
		Style:  f.Style,
		Order:  f.Order,
		Active: f.Active,
	}

	if s.ID != 0 {
		err := s.Update()
		if err != nil {
			if err == db.ErrNotFound {
				view.NotFound(w, r)
				return
			}
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

	id := strconv.Itoa(s.ID)

	http.Redirect(w, r, "/admin/slide?id="+id, http.StatusSeeOther)
}

func ListSlides(w http.ResponseWriter, r *http.Request) {
	slides, err := models.FetchSlides()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "slides", &view.View{
		Title:  "Slides",
		Slides: slides,
	})
}

func RemoveSlide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s := models.Slide{
		ID: utils.ToInt(id),
	}

	err := s.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/slides", http.StatusSeeOther)
}
