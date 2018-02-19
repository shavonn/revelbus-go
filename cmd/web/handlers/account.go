package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/flash"

	"golang.org/x/crypto/bcrypt"
)

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "user-dashboard", &view.View{})
}

func ProfileForm(w http.ResponseWriter, r *http.Request) {
	u, err := utils.LoggedIn(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	f := &models.UserForm{
		Name:  u.Name,
		Email: u.Email,
	}

	view.Render(w, r, "profile", &view.View{
		Form:  f,
		Title: "Update Profile",
	})
}

func PostProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
	}

	if !f.Valid() {
		view.Render(w, r, "profile", &view.View{
			Form:  f,
			Title: "Update Profile",
		})
		return
	}

	u, err := utils.LoggedIn(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	u.Name = f.Name
	u.Email = f.Email

	err = u.Update()
	if err == db.ErrDuplicateEmail {
		f.Errors["Email"] = "E-mail address is already in use"
		view.Render(w, r, "profile", &view.View{
			Form:  f,
			Title: "Update Profile",
		})
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = utils.SetUserSession(w, r, u)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyUpdated, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/u/profile", http.StatusSeeOther)
}

func PasswordForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "password", &view.View{
		Form:  new(models.UserForm),
		Title: "Update Password",
	})
}

func PostPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		OldPassword:     r.PostForm.Get("old_password"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
	}

	if !f.ValidPasswordUpdate() {
		view.Render(w, r, "password", &view.View{
			Form:  f,
			Title: "Update Password",
		})
		return
	}

	u, err := utils.LoggedIn(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = u.VerifyAndUpdatePassword(f.OldPassword, f.Password)
	if err == db.ErrInvalidCredentials || err == bcrypt.ErrHashTooShort {
		err = flash.Add(w, r, utils.MsgInvalidCredentials, "danger")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		view.Render(w, r, "password", &view.View{
			Form:  f,
			Title: "Update Password",
		})
		return
	} else if err != nil {
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfullyUpdated, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/u/password", http.StatusSeeOther)
}
