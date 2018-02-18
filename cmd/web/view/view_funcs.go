package view

import "time"

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
