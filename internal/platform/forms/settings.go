package forms

type SettingsForm struct {
	ID           string
	ContactBlurb string
	AboutBlurb   string
	AboutContent string
	Errors       map[string]string
}

func (f *SettingsForm) Valid() bool {
	v := newValidator()

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
