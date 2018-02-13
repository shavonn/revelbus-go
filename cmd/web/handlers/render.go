package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"time"

	"github.com/justinas/nosurf"
	"github.com/spf13/viper"
)

func humanDate(t time.Time) string {
	return t.Format("Jan 2, 2006 at 3:04 PM")
}

type view struct {
	ActiveKey string
	Err       appError
	Flash     flash.Msg
	Form      forms.Form
	Me        *db.User
	Path      string
	Title     string
	Token     string
	Trip      *db.Trip
	Trips     *db.Trips
	Vendor    *db.Vendor
	Vendors   *db.Vendors
	User      *db.User
	Users     *db.Users
}

type appError struct {
	Code    int
	Message string
}

func render(w http.ResponseWriter, r *http.Request, tpl string, v *view) {
	v.Path = r.URL.Path
	v.Token = nosurf.Token(r)

	flash, err := flash.Fetch(w, r)
	if err != nil {
		serverError(w, r, err)
		return
	}
	v.Flash = flash

	u, err := loggedIn(r)
	if err != nil {
		serverError(w, r, err)
		return
	}
	v.Me = u

	f := []string{
		filepath.Join(viper.GetString("files.tpl"), "app.html"),
		filepath.Join(viper.GetString("files.tpl"), "trip-nav.html"),
		filepath.Join(viper.GetString("files.tpl"), tpl),
	}

	fm := template.FuncMap{
		"humanDate": humanDate,
	}

	ts, err := template.New("").Funcs(fm).ParseFiles(f...)
	if err != nil {
		serverError(w, r, err)
	}

	err = ts.ExecuteTemplate(w, "app", v)
	if err != nil {
		serverError(w, r, err)
	}
}
