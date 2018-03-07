package handlers

import (
	"net/http"
	"revelbus/cmd/web/utils"
	"revelbus/cmd/web/view"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/domain/models"
	"revelbus/internal/platform/flash"

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

	if _, err = utils.UploadFile(w, r, "files", "uploads/files/"+fldr, false); err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
}

func ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := models.FetchFiles()
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

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/files", http.StatusSeeOther)
}
