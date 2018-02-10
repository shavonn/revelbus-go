package forms

type VendorForm struct {
	ID      string
	Name    string
	Address string
	City    string
	State   string
	Zip     string
	Phone   string
	Email   string
	URL     string
	Notes   string
	Brand   string
	Errors  map[string]string
}

func (f *VendorForm) Valid() bool {
	v := newValidator()

	v.Required("Name", f.Name)
	v.ValidEmail("Email", f.Email)
	v.ValidURL("URL", f.URL)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
