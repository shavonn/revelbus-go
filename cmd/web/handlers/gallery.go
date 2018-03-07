package handlers

import (
	"net/http"
	"revelbus/cmd/web/utils"
	"revelbus/cmd/web/view"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/domain/models"
	"revelbus/internal/platform/flash"
	"strconv"

	"github.com/gorilla/mux"
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

	err := g.Fetch()
	if err != nil {
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	f := &models.GalleryForm{
		ID:     strconv.Itoa(g.ID),
		Name:   g.Name.String,
		Folder: g.Folder.String,
	}

	view.Render(w, r, "gallery", &view.View{
		Title:   g.Name.String,
		Form:    f,
		Gallery: g,
	})
}

func PostGallery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.GalleryForm{
		ID:     r.PostForm.Get("id"),
		Name:   r.PostForm.Get("name"),
		Folder: r.PostForm.Get("folder"),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if f.ID == "" {
			v.Title = "New Gallery"
		}

		view.Render(w, r, "gallery", v)
	}

	var msg string

	g := models.Gallery{
		ID:   utils.ToInt(f.ID),
		Name: utils.NewNullStr(f.Name),
	}

	if g.ID != 0 {
		err := g.Update()
		if err != nil {
			if err == domain.ErrNotFound {
				view.NotFound(w, r)
				return
			}
			view.ServerError(w, r, err)
			return
		}

		uploads, err := utils.UploadFile(w, r, "files", "uploads/files/"+g.Folder.String, true)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		for _, f := range uploads {
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
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	id := strconv.Itoa(g.ID)

	http.Redirect(w, r, "/admin/gallery?id="+id, http.StatusSeeOther)
}

func ListGalleries(w http.ResponseWriter, r *http.Request) {
	galleries, err := models.FetchGalleries()
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

	err := g.Fetch()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = utils.DeleteFolder("uploads/files/" + g.Folder.String)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = g.Delete()
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

	f := &models.File{
		ID: utils.ToInt(fid),
	}

	err = utils.DeleteFile(f)
	if err != nil {
		if err == domain.ErrCannotDelete {
			err = flash.Add(w, r, utils.MsgCannotRemove, "warning")
			if err != nil {
				view.ServerError(w, r, err)
				return
			}
		}
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
