package utils

import (
	"net/url"
	"path/filepath"
	"revelforce/internal/platform/domain/models"
	"strconv"
	"strings"
	"time"
)

var dateFormat = "20060102T150400Z"

func GetCalendarLinks(t *models.Trip) map[string]string {
	m := make(map[string]string)
	address := ""

	if len(t.Venues) > 0 {
		for _, v := range t.Venues {
			if v.Primary {
				address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
			}
		}
	}

	m["google"] = google(t, address)
	m["yahoo"] = yahoo(t, address)
	m["ics"] = ics(t, address)
	return m
}

func stripSpaces(s string) string {
	return strings.Replace(s, " ", "+", -1)
}

func google(t *models.Trip, a string) string {
	uri := "https://www.google.com/calendar/render?action=TEMPLATE"
	uri = uri + "&text=" + stripSpaces(t.Title)
	uri = uri + "&dates=" + t.Start.Add(time.Hour*5).Format(dateFormat) + "/" + t.End.Add(time.Hour*4).Format(dateFormat)
	uri = uri + "&details=" + stripSpaces("For details, visit: http://www.revelbus.com/trip/"+t.Slug)

	if a != "" {
		uri = uri + "&location=" + stripSpaces(a)
	}

	uri = uri + "&ctz=America/New_York&sf=true&output=xml"

	return uri
}

func ics(t *models.Trip, a string) string {
	return filepath.Join("/assets/ical/", strconv.Itoa(t.ID)+".ics")
}

func yahoo(t *models.Trip, a string) string {
	uri := "https://calendar.yahoo.com/?v=60&view=d&type=20"
	uri = uri + "&title=" + url.QueryEscape(t.Title)
	uri = uri + "&st=" + t.Start.Add(time.Hour*5).Format(dateFormat)
	uri = uri + "&et=" + t.End.Format(dateFormat)
	uri = uri + "&desc=" + url.QueryEscape("For details, visit: http://www.revelbus.com/trip/"+t.Slug)

	if a != "" {
		uri = uri + "&in_loc=" + url.QueryEscape(a)
	}

	return uri
}
