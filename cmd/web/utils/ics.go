package utils

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"revelforce/internal/platform/domain/models"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

const (
	stampLayout    = "20060102T150405Z"
	dateLayout     = "20060102"
	dateTimeLayout = "20060102T150405"
)

type vCalendar struct {
	VERSION string
	PRODID  string
	URL     string

	NAME         string
	X_WR_CALNAME string
	DESCRIPTION  string
	X_WR_CALDESC string

	TIMEZONE_ID   string
	X_WR_TIMEZONE string

	REFRESH_INTERVAL string
	X_PUBLISHED_TTL  string

	COLOR    string
	CALSCALE string
	METHOD   string

	vComponent []vComponent
}

type vEvent struct {
	UID         string
	DTSTAMP     time.Time
	DTSTART     time.Time
	DTEND       time.Time
	SUMMARY     string
	DESCRIPTION string
	LOCATION    string
	TZID        string

	AllDay bool
}

type vComponent interface {
	encodeIcal(w io.Writer) error
}

func (e *vEvent) encodeIcal(w io.Writer) error {

	var timeStampLayout, timeStampType, tzidTxt string

	if e.AllDay {
		timeStampLayout = dateLayout
		timeStampType = "DATE"
	} else {
		timeStampLayout = dateTimeLayout
		timeStampType = "DATE-TIME"
		if len(e.TZID) == 0 || e.TZID == "UTC" {
			timeStampLayout = timeStampLayout + "Z"
		}
	}

	if len(e.TZID) != 0 && e.TZID != "UTC" {
		tzidTxt = "TZID=" + e.TZID + ";"
	}

	b := bufio.NewWriter(w)
	if _, err := b.WriteString("BEGIN:VEVENT\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("DTSTAMP:" + e.DTSTAMP.UTC().Format(stampLayout) + "\r\n"); err != nil {
		return err
	}
	if _, err := b.WriteString("UID:" + e.UID + "\r\n"); err != nil {
		return err
	}

	if len(e.TZID) != 0 && e.TZID != "UTC" {
		if _, err := b.WriteString("TZID:" + e.TZID + "\r\n"); err != nil {
			return err
		}
	}

	if _, err := b.WriteString("SUMMARY:" + e.SUMMARY + "\r\n"); err != nil {
		return err
	}
	if e.DESCRIPTION != "" {
		if _, err := b.WriteString("DESCRIPTION:" + e.DESCRIPTION + "\r\n"); err != nil {
			return err
		}
	}
	if e.LOCATION != "" {
		if _, err := b.WriteString("LOCATION:" + e.LOCATION + "\r\n"); err != nil {
			return err
		}
	}
	if _, err := b.WriteString("DTSTART;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.DTSTART.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("DTEND;" + tzidTxt + "VALUE=" + timeStampType + ":" + e.DTEND.Format(timeStampLayout) + "\r\n"); err != nil {
		return err
	}

	if _, err := b.WriteString("END:VEVENT\r\n"); err != nil {
		return err
	}

	return b.Flush()
}

func (c *vCalendar) encode(w io.Writer) error {
	var b = bufio.NewWriter(w)

	if _, err := b.WriteString("BEGIN:VCALENDAR\r\n"); err != nil {
		return err
	}

	attrs := []map[string]string{
		{"VERSION:": c.VERSION},
		{"PRODID:": c.PRODID},
		{"URL:": c.URL},
		{"NAME:": c.NAME},
		{"X-WR-CALNAME:": c.X_WR_CALNAME},
		{"DESCRIPTION:": c.DESCRIPTION},
		{"X-WR-CALDESC:": c.X_WR_CALDESC},
		{"TIMEZONE-ID:": c.TIMEZONE_ID},
		{"X-WR-TIMEZONE:": c.X_WR_TIMEZONE},
		{"REFRESH-INTERVAL;VALUE=DURATION:": c.REFRESH_INTERVAL},
		{"X-PUBLISHED-TTL:": c.X_PUBLISHED_TTL},
		{"COLOR:": c.COLOR},
		{"CALSCALE:": c.CALSCALE},
		{"METHOD:": c.METHOD},
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

func newCal() *vCalendar {
	cal := &vCalendar{
		VERSION:  "2.0",
		CALSCALE: "GREGORIAN",
	}
	cal.PRODID = "Revel Bus"
	cal.URL = "http://www.revelbus.com/trips"
	cal.METHOD = "PUBLISH"

	cal.NAME = "Revel Bus Trips"
	cal.X_WR_CALNAME = "Revel Bus Trips"

	cal.TIMEZONE_ID = "America/New_York"
	cal.X_WR_TIMEZONE = "America/New_York"
	return cal
}

func CreateICS(t *models.Trip) error {
	cal := newCal()
	address := ""

	if len(t.Venues) > 0 {
		for _, v := range t.Venues {
			if v.Primary {
				address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
			}
		}
	}

	e := &vEvent{
		UID:         "REVBUS" + strconv.Itoa(t.ID),
		DTSTAMP:     time.Now(),
		DTSTART:     t.Start,
		DTEND:       t.End,
		SUMMARY:     t.Title,
		LOCATION:    address,
		DESCRIPTION: "For details, visit: http://www.revelbus.com/trip/" + t.Slug,
		TZID:        "America/New_York",
		AllDay:      false,
	}

	cal.vComponent = append(cal.vComponent, e)

	ics, err := os.Create(filepath.Join(viper.GetString("files.static"), "ical/", strconv.Itoa(t.ID)+".ics"))
	if err != nil {
		return err
	}
	defer ics.Close()

	err = cal.encode(ics)
	return err
}
