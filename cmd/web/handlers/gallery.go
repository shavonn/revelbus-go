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
	"github.com/gosimple/slug"
)

func GalleryForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "gallery", &view.View{
			Form:  new(models.GalleryForm),
			Title: "New Gallery",
		})
		return
	}

	g := &models.Gallery{
		ID: utils.ToInt(id),
	}

	err := g.Get()
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.GalleryForm{
		ID:   strconv.Itoa(g.ID),
		Name: g.Name,
	}

	view.Render(w, r, "gallery", &view.View{
		Title:   g.Name,
		Form:    f,
		Gallery: g,
	})
}

func PostGallery(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.GalleryForm{
		ID:   r.PostForm.Get("id"),
		Name: r.PostForm.Get("name"),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if id == "" {
			v.Title = "New Gallery"
		}

		view.Render(w, r, "gallery", v)
	}

	var msg string

	g := models.Gallery{
		ID:   utils.ToInt(f.ID),
		Name: f.Name,
	}

	if id != "" {
		err := g.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		fldr := slug.Make(g.Name)

		uploads, err := utils.UploadFile(w, r, "files", "uploads/files/"+fldr)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		for _, upload := range uploads {
			f := models.File{
				Name: upload,
			}

			err := f.Create()
			if err != nil {
				view.ServerError(w, r, err)
				return
			}

			g.AttachImage(strconv.Itoa(f.ID))
			if err != nil {
				view.ServerError(w, r, err)
				return
			}
		}

		msg = utils.MsgSuccessfullyUpdated
	} else {
		err := g.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		id = strconv.Itoa(g.ID)
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/gallery?id="+id, http.StatusSeeOther)
}

func ListGalleries(w http.ResponseWriter, r *http.Request) {
	galleries, err := models.GetGalleries()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "galleries", &view.View{
		Title:     "Galleries",
		Galleries: galleries,
	})
}

func RemoveGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	g := models.Gallery{
		ID: utils.ToInt(id),
	}

	err := g.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/galleries", http.StatusSeeOther)
}

func DetachImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fid := vars["fid"]

	g := models.Gallery{
		ID: utils.ToInt(id),
	}

	err := g.DetachImage(fid)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemovedImage, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/gallery/?id="+id, http.StatusSeeOther)
}
