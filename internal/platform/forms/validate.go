package forms

import (
	"net/url"
	"regexp"
	"time"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var rxSlug = regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$")

type validator struct {
	Errors map[string]string
}

func NewValidator() *validator {
	return &validator{
		Errors: make(map[string]string),
	}
}

func (v *validator) Required(k string, i string) {
	if len(i) == 0 {
		v.Errors[k] = k + " is required."
	}
}

func (v *validator) ValidEmail(k string, i string) {
	if i != "" {
		if !rxEmail.MatchString(i) {
			v.Errors[k] = "Please enter a valid email address."
		}
	}
}

func (v *validator) ValidDateTime(k string, i string) {
	if i != "" {
		if _, err := time.Parse("2006-01-02 15:04", i); err != nil {
			v.Errors[k] = "Please enter a valid date/time."
		}
	}
}

func (v *validator) ValidDateTimeRange(k string, s string, e string) {
	if s != "" && e != "" {
		s, _ := time.Parse("2006-01-02 15:04", s)
		e, _ := time.Parse("2006-01-02 15:04", e)
		hs := e.Sub(s).Hours()
		if hs < 0 {
			v.Errors[k] = "Please enter a valid date/time range."
		}
	}
}

func (v *validator) ValidSlug(k string, i string) {
	if i != "" {
		if !rxSlug.MatchString(i) {
			v.Errors[k] = "Please enter a valid slug."
		}
	}
}

func (v *validator) ValidURL(k string, i string) {
	if i != "" {
		if _, err := url.ParseRequestURI(i); err != nil {
			v.Errors[k] = "Please enter a valid URL."
		}
	}
}
