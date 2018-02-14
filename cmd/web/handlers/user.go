package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strconv"

	"github.com/gorilla/mux"
)

func userForm(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		render(w, r, "user", &view{
			Form:  new(forms.UserForm),
			Title: "New User",
		})
		return
	}

	u := &db.User{
		ID: toInt(id),
	}

	err := u.Get()
	if err == db.ErrNotFound {
		notFound(w, r)
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.UserForm{
		ID:    id,
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}

	render(w, r, "user", &view{
		Form: f,
		User: u,
	})
}

func postUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		ID:    r.PostForm.Get("id"),
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
		Role:  r.PostForm.Get("role"),
	}

	if !f.Valid() {
		v := &view{
			Form: f,
		}

		if id == "" {
			v.Title = "New User"
		}

		render(w, r, "user", v)
	}

	var msg string

	u := db.User{
		Name:  f.Name,
		Email: f.Email,
		Role:  f.Role,
	}

	if id != "" {
		u.ID = toInt(id)

		err := u.Update()
		if err != nil {
			serverError(w, r, err)
			return
		}

		if len(r.Form["reset_password"]) == 1 {
			pw := randomString(14)

			err := u.UpdatePassword(pw)
			if err != nil {
				serverError(w, r, err)
				return
			}

			email.NewPassword(u.Email, pw)
		}

		msg = MsgSuccessfullyUpdated
	} else {
		pw := randomString(14)
		u.Password = pw

		err := u.Create()
		if err != nil {
			serverError(w, r, err)
			return
		}

		email.NewPassword(u.Email, pw)

		id = strconv.Itoa(u.ID)
		msg = MsgSuccessfullyCreated
	}

	err = flash.Add(w, r, msg, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/user?id="+id, http.StatusSeeOther)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers()
	if err != nil {
		serverError(w, r, err)
		return
	}

	render(w, r, "users", &view{
		Title: "Users",
		Users: &users,
	})
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	u := &db.User{
		ID: toInt(id),
	}

	err := u.Delete()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyRemoved, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
