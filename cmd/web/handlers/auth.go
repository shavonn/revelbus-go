package handlers

import (
	"net/http"
	"revelforce/cmd/web/utils"
	"revelforce/cmd/web/view"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/db/models"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"

	"github.com/gorilla/mux"
)

func SignupForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "signup", &view.View{
		Form:  new(models.UserForm),
		Title: "Signup",
	})
}

func PostSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		Name:            r.PostForm.Get("name"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
		Role:            "user",
	}

	if !f.ValidSignup() {
		view.Render(w, r, "signup", &view.View{
			Form: f,
		})
		return
	}

	u := models.User{
		Name:     f.Name,
		Email:    f.Email,
		Password: f.Password,
		Role:     f.Role,
	}

	err = u.Create()
	if err != nil {
		if err == db.ErrDuplicateEmail {
			f.Errors["Email"] = "E-mail address is already in use"
			view.Render(w, r, "signup", &view.View{
				Form:  f,
				Title: "Signup",
			})
			return
		}
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgSuccessfulSignup, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "login", &view.View{
		Form:  new(models.UserForm),
		Title: "Login",
	})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !f.ValidLogin() {
		view.Render(w, r, "login", &view.View{
			Form: f,
		})
		return
	}

	u := &models.User{
		Email: f.Email,
	}

	err = u.VerifyUser(f.Password)
	if err != nil {
		if err == db.ErrInvalidCredentials {
			err = flash.Add(w, r, utils.MsgUnsuccessfulLogin, "danger")
			if err != nil {
				view.ServerError(w, r, err)
				return
			}

			view.Render(w, r, "login", &view.View{
				Form:  f,
				Title: "Login",
			})
			return
		}
		view.ServerError(w, r, err)
		return
	}

	err = utils.SetUserSession(w, r, u)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/u", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := utils.RemoveUserSession(w, r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func ForgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, r, "forgot", &view.View{
		Form:  new(models.UserForm),
		Title: "Forgot Password",
	})
}

func PostForgotPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		Email: r.PostForm.Get("email"),
	}

	if !f.ValidForgot() {
		view.Render(w, r, "forgot", &view.View{
			Form:  f,
			Title: "Forgot Password",
		})
		return
	}

	u := models.User{
		Email: f.Email,
	}

	err = u.Fetch()
	if err == nil {
		rh := utils.RandomString(20)

		err = u.SetRecover(rh)
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		email.RecoverAccount(u.Email, rh)
	}

	err = flash.Add(w, r, utils.MsgRecoverySent, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]
	hash := vars["hash"]

	if email == "" || hash == "" {
		err := flash.Add(w, r, utils.MsgInvalidRecovery, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/forgot", http.StatusSeeOther)
	}

	u := &models.User{
		Email: email,
	}

	err := u.CheckRecover(hash)
	if err != nil {
		err := flash.Add(w, r, utils.MsgInvalidRecovery, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/forgot", http.StatusSeeOther)
	}

	view.Render(w, r, "reset", &view.View{
		Form: &models.UserForm{
			Email:        email,
			RecoveryHash: hash,
		},
		Title: "Reset Password",
	})
}

func PostPasswordReset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		view.ClientError(w, r, http.StatusBadRequest)
		return
	}

	f := &models.UserForm{
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
		RecoveryHash:    r.PostForm.Get("recovery_hash"),
	}

	if !f.ValidPassword() || !f.ValidForgot() {
		view.Render(w, r, "reset", &view.View{
			Form:  f,
			Title: "Reset Password",
		})
		return
	}

	u := &models.User{
		Email: f.Email,
	}

	err = u.Recover(f.RecoveryHash, f.Password)
	if err != nil {
		if err == db.ErrInvalidCredentials {
			err = flash.Add(w, r, utils.MsgInvalidCredentials, "danger")
			if err != nil {
				view.ServerError(w, r, err)
				return
			}

			view.Render(w, r, "reset", &view.View{
				Form:  f,
				Title: "Reset Password",
			})
			return
		}
		view.ServerError(w, r, err)
		return
	}

	err = flash.Add(w, r, utils.MsgPasswordResetSuccessful, "success")
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}
