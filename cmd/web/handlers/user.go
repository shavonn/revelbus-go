package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"
	"strconv"

	"github.com/gorilla/mux"
)

func UserForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		view.Render(w, r, "user", &view.View{
			Form:  new(models.UserForm),
			Title: "New User",
		})
		return
	}

	u := &models.User{
		ID: utils.ToInt(id),
	}

	err := u.Get()
	if err == db.ErrNotFound {
		view.NotFound(w, r)
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.UserForm{
		ID:    id,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}

	view.Render(w, r, "user", &view.View{
		Form: f,
		User: u,
	})
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		ID:    r.PostForm.Get("id"),
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
		Role:  r.PostForm.Get("role"),
	}

	if !f.Valid() {
		v := &view.View{
			Form: f,
		}

		if id == "" {
			v.Title = "New User"
		}

		view.Render(w, r, "user", v)
	}

	var msg string

	u := models.User{
		Name:  f.Name,
		Email: f.Email,
		Role:  f.Role,
	}

	if id != "" {
		u.ID = utils.ToInt(id)

		err := u.Update()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		if len(r.Form["reset_password"]) == 1 {
			pw := utils.RandomString(14)

			err := u.UpdatePassword(pw)
			if err != nil {
				view.ServerError(w, r, err)
				return
			}

			email.NewPassword(u.Email, pw)
		}

		msg = utils.MsgSuccessfullyUpdated
	} else {
		pw := utils.RandomString(14)
		u.Password = pw

		err := u.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		email.NewPassword(u.Email, pw)

		id = strconv.Itoa(u.ID)
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/user?id="+id, http.StatusSeeOther)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	view.Render(w, r, "users", &view.View{
		Title: "Users",
		Users: &users,
	})
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	u := &models.User{
		ID: utils.ToInt(id),
	}

	err := u.Delete()
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyRemoved, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
