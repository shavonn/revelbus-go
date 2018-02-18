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
	return t.Format("Mon, Jan 2, 2006 at 3:04 PM")
}

func month(s time.Time, e time.Time) string {
	if s.Month() == e.Month() {
		return s.Format("Jan")
	}
	return s.Format("Jan") + " - " + e.Format("Jan")
}

func day(s time.Time, e time.Time) string {
	if s.Month() == e.Month() && s.Day() != e.Day() {
		return s.Format("2") + " - " + e.Format("2")
	}
	return s.Format("2")
}

func blurb(s string) string {
	if len(s) > 105 {
		return s[:105]
	}
	return s
}

type view struct {
	ActiveKey   string
	Err         appError
	FAQ         *db.FAQ
	FAQs        *db.FAQs
	FAQGrouped  *db.GroupedFAQ
	Flash       flash.Msg
	Form        forms.Form
	HeaderStyle string
	Me          *db.User
	Path        string
	Slide       *db.Slide
	Slides      *db.Slides
	Title       string
	Token       string
	Trip        *db.Trip
	Trips       *db.Trips
	Vendor      *db.Vendor
	Vendors     *db.Vendors
	User        *db.User
	Users       *db.Users
}

type appError struct {
	Code    int
	Message string
}

func render(w http.ResponseWriter, r *http.Request, tpl string, v *view) {
	v.Path = r.URL.Path
	v.Token = nosurf.Token(r)
	v.HeaderStyle = getHeaderStyle(tpl)

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
		"month":     month,
		"day":       day,
		"blurb":     blurb,
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

func getHeaderStyle(t string) string {
	switch t {
	case "trips":
		return "game_guys"
	case "faq":
		return "golfers"
	case "contact":
		return "swimmers"
	default:
		return ""
	}
}
