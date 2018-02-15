package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"

	"golang.org/x/crypto/bcrypt"
)

func profileForm(w http.ResponseWriter, r *http.Request) {
	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	f := &forms.UserForm{
		Name:  u.Name,
		Email: u.Email,
	}

	render(w, r, "profile", &view{
		Form:  f,
		Title: "Update Profile",
	})
}

func postProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
	}

	if !f.Valid() {
		render(w, r, "profile", &view{
			Form:  f,
			Title: "Update Profile",
		})
		return
	}

	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	u.Name = f.Name
	u.Email = f.Email

	err = u.Update()
	if err == db.ErrDuplicateEmail {
		f.Errors["Email"] = "E-mail address is already in use"
		render(w, r, "profile", &view{
			Form:  f,
			Title: "Update Profile",
		})
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = setUserSession(w, r, u)
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyUpdated, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/u/profile", http.StatusSeeOther)
}

func passwordForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, "password", &view{
		Form:  new(forms.UserForm),
		Title: "Update Password",
	})
}

func postPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		OldPassword:     r.PostForm.Get("old_password"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}

	if !f.ValidPasswordUpdate() {
		render(w, r, "password", &view{
			Form:  f,
			Title: "Update Password",
		})
		return
	}

	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = u.VerifyAndUpdatePassword(f.OldPassword, f.Password)
	if err == db.ErrInvalidCredentials || err == bcrypt.ErrHashTooShort {
		err = flash.Add(w, r, MsgInvalidCredentials, "danger")
		if err != nil {
			serverError(w, r, err)
			return
		}

		render(w, r, "password", &view{
			Form:  f,
			Title: "Update Password",
		})
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfullyUpdated, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/u/password", http.StatusSeeOther)
}
