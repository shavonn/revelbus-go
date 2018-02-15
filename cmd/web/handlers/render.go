package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/flash"
	"revelforce/internal/platform/forms"
	"strings"
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

	t, err := parseTemplates()
	if err != nil {
		serverError(w, r, err)
		return
	}

	err = t.ExecuteTemplate(w, tpl, v)
	if err != nil {
		serverError(w, r, err)
		return
	}
}

func parseTemplates() (*template.Template, error) {
	fm := template.FuncMap{
		"humanDate": humanDate,
	}
	templ := template.New("").Funcs(fm)
	err := filepath.Walk(viper.GetString("files.tpl"), func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	return templ, nil
}
