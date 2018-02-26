package utils

import (
	"os"
	"path/filepath"
	"revelforce/internal/platform/db/models"
	"strconv"
	"time"

	"github.com/soh335/ical"
	"github.com/spf13/viper"
)

func NewCal() *ical.VCalendar {
	cal := ical.NewBasicVCalendar()
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
	cal := NewCal()

	e := &ical.VEvent{
		UID:         strconv.Itoa(t.ID),
		DTSTAMP:     time.Now(),
		DTSTART:     t.Start,
		DTEND:       t.End,
		SUMMARY:     t.Title,
		DESCRIPTION: t.Blurb,
		TZID:        "America/New_York",
		AllDay:      false,
	}

	cal.VComponent = append(cal.VComponent, e)

	ics, err := os.Create(filepath.Join(viper.GetString("files.static"), "ical/", strconv.Itoa(t.ID)+".ics"))
	if err != nil {
		return err
	}
	defer ics.Close()

	err = cal.Encode(ics)
	return err
}
