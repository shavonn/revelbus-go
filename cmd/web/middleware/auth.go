package middleware

import (
	"net/http"
	"revelbus/cmd/web/utils"
	"revelbus/cmd/web/view"
	"revelbus/internal/platform/flash"
)

func RequireLogin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := utils.IsAuthenticated(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if u == nil {
		err = flash.Add(w, r, utils.MsgMustBeLoggedIn, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, "/auth/login", 302)
		return
	}
	next(w, r)
}

func RequireAdmin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := utils.IsAuthenticated(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if u == nil {
		err = flash.Add(w, r, utils.MsgMustBeLoggedIn, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, "/auth/login", 302)
		return
	} else if u.Role.String != "admin" {
		err = flash.Add(w, r, utils.MsgMustBeAdmin, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, "/u", 302)
		return
	}
	next(w, r)
}

func RequireGuest(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	u, err := utils.IsAuthenticated(r)
	if err != nil {
		view.ServerError(w, r, err)
		return
	}

	if u != nil {
		err = flash.Add(w, r, utils.MsgAlreadyAuthenticated, "warning")
		if err != nil {
			view.ServerError(w, r, err)
			return
		}

		http.Redirect(w, r, "/", 302)
		return
	}
	next(w, r)
}
