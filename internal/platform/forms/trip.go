package forms

type TripForm struct {
	ID           string
	Title        string
	Slug         string
	Status       string
	Blurb        string
	Description  string
	Start        string
	End          string
	Price        string
	TicketingURL string
	Notes        string
	Image        string
	Errors       map[string]string
}

func (f *TripForm) Valid() bool {
	v := newValidator()

	v.Required("Title", f.Title)
	v.ValidSlug("Slug", f.Slug)
	v.ValidDateTime("Start", f.Start)
	v.ValidDateTime("End", f.End)
	v.ValidDateTimeRange("End", f.Start, f.End)
	v.ValidURL("TicketingURL", f.TicketingURL)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
