package handlers

import (
	"net/http"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/email"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func signupForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, "signup.html", &view{
		Form:  new(forms.UserForm),
		Title: "Signup",
	})
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		Name:            r.PostForm.Get("name"),
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
		Role:            "user",
	}

	if !f.ValidSignup() {
		render(w, r, "signup.html", &view{
			Form: f,
		})
		return
	}

	u := db.User{
		Name:     f.Name,
		Email:    f.Email,
		Password: f.Password,
		Role:     f.Role,
	}

	err = u.Create()
	if err == db.ErrDuplicateEmail {
		f.Errors["Email"] = "E-mail address is already in use"
		render(w, r, "signup.html", &view{
			Form:  f,
			Title: "Signup",
		})
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgSuccessfulSignup, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, "login.html", &view{
		Form:  new(forms.UserForm),
		Title: "Login",
	})
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !f.ValidLogin() {
		render(w, r, "login.html", &view{
			Form: f,
		})
		return
	}

	u := &db.User{
		Email: f.Email,
	}

	err = u.VerifyUser(f.Password)
	if err == db.ErrInvalidCredentials {
		err = flash.Add(w, r, MsgUnsuccessfulLogin, "danger")
		if err != nil {
			serverError(w, r, err)
			return
		}

		render(w, r, "login.html", &view{
			Form:  f,
			Title: "Login",
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	err := removeUserSession(w, r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", 303)
}

func forgotPasswordForm(w http.ResponseWriter, r *http.Request) {
	render(w, r, "forgot.html", &view{
		Form:  new(forms.UserForm),
		Title: "Forgot Password",
	})
}

func postForgotPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		Email: r.PostForm.Get("email"),
	}

	if !f.ValidForgot() {
		render(w, r, "forgot.html", &view{
			Form:  f,
			Title: "Forgot Password",
		})
		return
	}

	u := db.User{
		Email: f.Email,
	}

	err = u.Get()
	if err == nil {
		rh := randomString(20)

		err = u.SetRecover(rh)
		if err != nil {
			serverError(w, r, err)
			return
		}

		email.RecoverAccount(u.Email, rh)
	}

	err = flash.Add(w, r, MsgRecoverySent, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetPasswordForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]
	hash := vars["hash"]

	if email == "" || hash == "" {
		err := flash.Add(w, r, MsgInvalidRecovery, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/forgot", http.StatusSeeOther)
	}

	u := &db.User{
		Email: email,
	}

	err := u.CheckRecover(hash)
	if err != nil {
		err := flash.Add(w, r, MsgInvalidRecovery, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/forgot", http.StatusSeeOther)
	}

	render(w, r, "reset.html", &view{
		Form: &forms.UserForm{
			Email:        email,
			RecoveryHash: hash,
		},
		Title: "Reset Password",
	})
}

func postPasswordReset(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		clientError(w, r, http.StatusBadRequest)
		return
	}

	f := &forms.UserForm{
		Email:           r.PostForm.Get("email"),
		Password:        r.PostForm.Get("password"),
		ConfirmPassword: r.PostForm.Get("confirm_password"),
		RecoveryHash:    r.PostForm.Get("recovery_hash"),
	}

	if !f.ValidPassword() || !f.ValidForgot() {
		render(w, r, "reset.html", &view{
			Form:  f,
			Title: "Reset Password",
		})
		return
	}

	u := &db.User{
		Email: f.Email,
	}

	err = u.Recover(f.RecoveryHash, f.Password)
	if err == db.ErrInvalidCredentials || err == bcrypt.ErrHashTooShort {
		err = flash.Add(w, r, MsgInvalidCredentials, "danger")
		if err != nil {
			serverError(w, r, err)
			return
		}

		render(w, r, "reset.html", &view{
			Form:  f,
			Title: "Reset Password",
		})
		return
	} else if err != nil {
		serverError(w, r, err)
		return
	}

	err = flash.Add(w, r, MsgPasswordResetSuccessful, "success")
	if err != nil {
		serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func requireLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if u == nil {
		err = flash.Add(w, r, MsgMustBeLoggedIn, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/auth/login", 302)
		return
	}
	next(w, r)
}

func requireAdmin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if u == nil {
		err = flash.Add(w, r, MsgMustBeLoggedIn, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/auth/login", 302)
		return
	} else if u.Role != "admin" {
		err = flash.Add(w, r, MsgMustBeAdmin, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/u", 302)
		return
	}
	next(w, r)
}

func requireGuest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}

	if u != nil {
		err = flash.Add(w, r, MsgAlreadyAuthenticated, "warning")
		if err != nil {
			serverError(w, r, err)
			return
		}

		http.Redirect(w, r, "/", 302)
		return
	}
	next(w, r)
}
