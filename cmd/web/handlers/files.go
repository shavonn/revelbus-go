package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/flash"

	"github.com/gorilla/mux"
)

func UploadForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "upload", &view.View{
		Title: "Upload",
	})
}

func PostUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	fldr := r.PostForm.Get("fldr")

	_, err = utils.UploadFile(w, r, "files", "uploads/files/"+fldr, false)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := models.GetFiles()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "files", &view.View{
		Title: "Files",
		Files: files,
	})
}
func RemoveFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	f := &models.File{
		ID: utils.ToInt(id),
	}

	err := utils.DeleteFile(f)
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err == db.ErrCannotDelete {
		err = flash.Add(w, r, utils.MsgCannotRemove, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	} else {
		err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
}
