package forms

type SlideForm struct {
	ID     string
	Title  string
	Blurb  string
	Style  string
	Active bool
	Order  string
	Errors map[string]string
}

func (f *SlideForm) Valid() bool {
	v := newValidator()

	v.Required("Title", f.Title)
	v.Required("Blurb", f.Blurb)
	v.Required("Style", f.Style)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
