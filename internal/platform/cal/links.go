package cal

import (
	"net/url"
	"path/filepath"
	"revelforce/internal/platform/domain/models"
)

func GetCalendarLinks(t *models.Trip) map[string]string {
	m := make(map[string]string)
	e := tripToVEvent(t)

	m["google"] = e.google()
	m["yahoo"] = e.yahoo()
	m["ics"] = e.ics()
	return m
}

func (ve *vEvent) google() string {
	uri := "https://www.google.com/calendar/render?action=TEMPLATE"
	uri = uri + "&text=" + stripSpaces(ve.summary)
	uri = uri + "&dates=" + ve.dtStart.Format(dateFormat) + "/" + ve.dtEnd.Format(dateFormat)
	uri = uri + "&details=" + stripSpaces("For details, visit: http://www.revelbus.com/trip/"+ve.slug)

	if ve.location != "" {
		uri = uri + "&location=" + stripSpaces(ve.location)
	}

	uri = uri + "&ctz=America/New_York&sf=true&output=xml"

	return uri
}

func (ve *vEvent) ics() string {
	return filepath.Join("/ical/", ve.slug+".ics")
}

func (ve *vEvent) yahoo() string {
	uri := "https://calendar.yahoo.com/?v=60&view=d&type=20"
	uri = uri + "&title=" + url.QueryEscape(ve.summary)
	uri = uri + "&st=" + ve.dtStart.Format(dateFormat)
	uri = uri + "&et=" + ve.dtEnd.Format(dateFormat)
	uri = uri + "&desc=" + url.QueryEscape("For details, visit: http://www.revelbus.com/trip/"+ve.slug)

	if ve.location != "" {
		uri = uri + "&in_loc=" + url.QueryEscape(ve.location)
	}

	return uri
}
