package forms

type TourForm struct {
	ID           string
	Name         string
	Slug         string
	Status       string
	Description  string
	Start        string
	End          string
	TicketingURL string
	Notes        string
	Errors       map[string]string
}

func (f *TourForm) Valid() bool {
	v := newValidator()

	v.Required("Name", f.Name)
	v.ValidSlug("Slug", f.Slug)
	v.ValidDateTime("Start", f.Start)
	v.ValidDateTime("End", f.End)
	v.ValidDateTimeRange("End", f.Start, f.End)
	v.ValidURL("TicketingURL", f.TicketingURL)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
