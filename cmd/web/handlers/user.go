package handlers

import (
	"net/http"
	"revelbus/cmd/web/utils"
	"revelbus/cmd/web/view"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/domain/models"
	"revelbus/internal/platform/emails"
	"revelbus/internal/platform/flash"
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

	err := u.Fetch()
	if err != nil {
		if err == domain.ErrNotFound {
			view.NotFound(w, r)
			return
		}
		view.ServerError(w, r, err)
		return
	}

	f := &models.UserForm{
		ID:    id,
		Name:  u.Name.String,
		Email: u.Email.String,
		Role:  u.Role.String,
	}

	view.Render(w, r, "user", &view.View{
		Title: f.Name,
		Form:  f,
	})
}

func PostUser(w http.ResponseWriter, r *http.Request) {
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

		if f.ID == "" {
			v.Title = "New User"
		}

		view.Render(w, r, "user", v)
	}

	var msg string

	u := models.User{
		ID:    utils.ToInt(f.ID),
		Name:  utils.NewNullStr(f.Name),
		Email: utils.NewNullStr(f.Email),
		Role:  utils.NewNullStr(f.Role),
	}

	if u.ID != 0 {
		err := u.Update()
		if err != nil {
			if err == domain.ErrNotFound {
				view.NotFound(w, r)
				return
			}
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

			emails.NewPassword(u.Email.String, pw)
		}

		msg = utils.MsgSuccessfullyUpdated
	} else {
		pw := utils.RandomString(14)
		u.Password = utils.NewNullStr(pw)

		err := u.Create()
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		emails.NewPassword(u.Email.String, pw)
		msg = utils.MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	id := strconv.Itoa(u.ID)

	http.Redirect(w, r, "/admin/user?id="+id, http.StatusSeeOther)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.FetchUsers()
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
