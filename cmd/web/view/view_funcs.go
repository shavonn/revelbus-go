package view

import "time"

func humanDate(t time.Time) string {
	return t.Format("Mon, Jan 2, 2006 at 3:04 PM")
}

func seoDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func notTrip(s string) bool {
	return (s != "trip")
}

func getShortMonth(s time.Time, e time.Time) string {
	if s.Month() == e.Month() {
		return s.Format("Jan")
	}
	return s.Format("Jan") + " - " + e.Format("Jan")
}

func getDateRange(s time.Time, e time.Time) string {
	if s.Month() == e.Month() && s.Day() != e.Day() {
		return s.Format("2") + " - " + e.Format("2")
	}
	return s.Format("2")
}

func numToMonth(m string) string {
	switch m {
	case "01":
		return "January"
	case "02":
		return "February"
	case "03":
		return "March"
	case "04":
		return "April"
	case "05":
		return "May"
	case "06":
		return "June"
	case "07":
		return "July"
	case "08":
		return "August"
	case "09":
		return "September"
	case "10":
		return "October"
	case "11":
		return "November"
	case "12":
		return "December"
	default:
		return "00"
	}
}

func blurb(s string) string {
	if len(s) > 105 {
		return s[:105]
	}
	return s
}
