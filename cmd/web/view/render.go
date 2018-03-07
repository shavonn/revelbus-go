package view

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"revelbus/cmd/web/utils"
	"revelbus/internal/platform/domain/models"
	"revelbus/internal/platform/flash"
	"revelbus/internal/platform/forms"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/spf13/viper"
)

type View struct {
	ActiveKey    string
	Blurb        string
	Content      template.HTML
	Err          appError
	FAQs         *models.FAQs
	Files        *models.Files
	FAQGrouped   *models.GroupedFAQs
	Gallery      *models.Gallery
	Galleries    *models.Galleries
	TripsGrouped *models.GroupedTrips
	Flash        flash.Msg
	Form         forms.Form
	HeaderStyle  string
	Me           *models.User
	Path         string
	Slides       *models.Slides
	Title        string
	Token        string
	Trip         *models.Trip
	Trips        *models.Trips
	Vendors      *models.Vendors
	Users        *models.Users
}

type appError struct {
	Code    int
	Message string
}

func Render(w http.ResponseWriter, r *http.Request, tpl string, v *View) {
	v.Path = r.URL.Path
	v.Token = nosurf.Token(r)
	v.HeaderStyle = getHeaderStyle(tpl)

	flash, err := flash.Fetch(w, r)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	v.Flash = flash

	u, err := utils.IsAuthenticated(r)
	if err != nil {
		ServerError(w, r, err)
		return
	}
	v.Me = u

	t, err := parseTemplates()
	if err != nil {
		ServerError(w, r, err)
		return
	}

	err = t.ExecuteTemplate(w, tpl, v)
	if err != nil {
		ServerError(w, r, err)
		return
	}
}

func parseTemplates() (*template.Template, error) {
	fm := template.FuncMap{
		"humanDate":     humanDate,
		"getShortMonth": getShortMonth,
		"getDateRange":  getDateRange,
		"numToMonth":    numToMonth,
		"blurb":         blurb,
		"seoDate":       seoDate,
		"notTrip":       notTrip,
	}
	templ := template.New("").Funcs(fm)
	err := filepath.Walk(viper.GetString("files.tpl"), func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			if _, err = templ.ParseFiles(path); err != nil {
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
	case "about":
		return "wine_gals"
	case "contact":
		return "swimmers"
	case "faq":
		return "golfers"
	case "trips":
		return "game_guys"
	default:
		return ""
	}
}
