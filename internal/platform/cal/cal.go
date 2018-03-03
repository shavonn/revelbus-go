package cal

import (
	"io"
	"revelforce/internal/platform/domain/models"
	"strconv"
	"strings"
	"time"
)

const (
	dateLayout     = "20060102"
	dateTimeLayout = "20060102T150405"
)

type vCalendar struct {
	version         string
	prodID          string
	url             string
	name            string
	description     string
	timezone        string
	refreshInterval string
	color           string
	calScale        string
	method          string

	vComponent []vComponent
}
type vComponent interface {
	encodeIcal(w io.Writer) error
}

type vEvent struct {
	uID         string
	dtStamp     time.Time
	dtStart     time.Time
	dtEnd       time.Time
	slug        string
	summary     string
	description string
	location    string
	tzID        string
	allDay      bool
}

func tripToVEvent(t *models.Trip) *vEvent {
	address := ""

	if len(t.Venues) > 0 {
		for _, v := range t.Venues {
			if v.Primary {
				address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
			}
		}

		if address == "" {
			v := t.Venues[0]
			address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
		}
	}

	return &vEvent{
		uID:         "REVBUS" + strconv.Itoa(t.ID),
		dtStamp:     time.Now(),
		dtStart:     t.Start,
		dtEnd:       t.End,
		summary:     t.Title.String,
		location:    address,
		description: "For details, visit: http://www.revelbus.com/trip/" + t.Slug.String,
		tzID:        "EDST",
		allDay:      false,

		slug: t.Slug.String,
	}
}

func stripSpaces(s string) string {
	return strings.Replace(s, " ", "+", -1)
}
