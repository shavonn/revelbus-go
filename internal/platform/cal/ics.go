package cal

import (
	"bufio"
	"io"
	"net/http"
	"revelforce/internal/platform/domain/models"
	"strconv"
	"time"
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

type vEvent struct {
	uID         string
	dtStamp     time.Time
	dtStart     time.Time
	dtEnd       time.Time
	summary     string
	description string
	location    string
	tzID        string
	allDay      bool
}

type vComponent interface {
	encodeIcal(w io.Writer) error
}

// encode event
func (e *vEvent) encodeIcal(w io.Writer) error {

	var timeStampLayout, timeStampType, tzidTxt string

	if e.allDay {
		timeStampLayout = dateLayout
		timeStampType = "DATE"
	} else {
		timeStampLayout = dateTimeLayout
		timeStampType = "DATE-TIME"
		if len(e.tzID) == 0 || e.tzID == "UTC" {
			timeStampLayout = timeStampLayout + "Z"
		}
	}

	if len(e.tzID) != 0 && e.tzID != "UTC" {
		tzidTxt = "TZID=" + e.tzID + ";"
	}

	b := bufio.NewWriter(w)
	if _, err := b.WriteString("BEGIN:VEVENT\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTSTAMP:" + e.dtStamp.UTC().Format(dateFormat) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("UID:" + e.uID + "\r\n"); err != nil {
		return err
	}

	if len(e.tzID) != 0 && e.tzID != "UTC" {
		if _, err := b.WriteString("TZID:" + e.tzID + "\r\n"); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("SUMMARY:" + e.summary + "\r\n"); err != nil {
		return err
	}
	if e.description != "" {
		if _, err := b.WriteString("DESCRIPTION:" + e.description + "\r\n"); err != nil {
			return err
		}
	}
	if e.location != "" {
		if _, err := b.WriteString("LOCATION:" + e.location + "\r\n"); err != nil {
			return err
		}
	}
	if _, err := b.WriteString("DTSTART;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.dtStart.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("DTEND;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.dtEnd.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("END:VEVENT\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

// encode calendar
func (c *vCalendar) encode(w io.Writer) error {
	var b = bufio.NewWriter(w)

	if _, err := b.WriteString("BEGIN:VCALENDAR\r\n"); err != nil {
		return err
	}

	attrs := []map[string]string{
		{"VERSION:": c.version},
		{"PRODID:": c.prodID},
		{"URL:": c.url},
		{"NAME:": c.name},
		{"X-WR-CALNAME:": c.name},
		{"DESCRIPTION:": c.description},
		{"X-WR-CALDESC:": c.description},
		{"TIMEZONE-ID:": c.timezone},
		{"X-WR-TIMEZONE:": c.timezone},
		{"REFRESH-INTERVAL;VALUE=DURATION:": c.refreshInterval},
		{"X-PUBLISHED-TTL:": c.refreshInterval},
		{"COLOR:": c.color},
		{"CALSCALE:": c.calScale},
		{"METHOD:": c.method},
	}

	for _, item := range attrs {
		for k, v := range item {
			if len(v) == 0 {
				continue
			}
			if _, err := b.WriteString(k + v + "\r\n"); err != nil {
				return err
			}
		}
	}

	for _, component := range c.vComponent {
		if err := component.encodeIcal(b); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("END:VCALENDAR\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

// create new cal with my defaults
func newCal() *vCalendar {
	cal := &vCalendar{
		version:         "2.0",
		calScale:        "GREGORIAN",
		prodID:          "Revel Bus",
		name:            "Revel Bus Trips",
		url:             "http://www.revelbus.com/trips",
		method:          "PUBLISH",
		timezone:        "EDST",
		refreshInterval: "PT12H",
	}

	return cal
}

func GenerateICS(w http.ResponseWriter, t *models.Trip) error {
	cal := newCal()
	address := getAddress(t)

	e := &vEvent{
		uID:         "REVBUS" + strconv.Itoa(t.ID),
		dtStamp:     time.Now(),
		dtStart:     t.Start,
		dtEnd:       t.End,
		summary:     t.Title,
		location:    address,
		description: "For details, visit: http://www.revelbus.com/trip/" + t.Slug,
		tzID:        "EDST",
		allDay:      false,
	}

	cal.vComponent = append(cal.vComponent, e)

	w.Header().Set("Content-Type", "text/calendar")
	err := cal.encode(w)
	return err
}
