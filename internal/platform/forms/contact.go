package forms

type ContactForm struct {
	Name    string
	Phone   string
	Email   string
	Message string
	Errors  map[string]string
}

func (f *ContactForm) Valid() bool {
	v := newValidator()

	v.Required("Name", f.Name)
	v.Required("Email", f.Email)
	v.Required("Message", f.Message)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
