package forms

type UserForm struct {
	ID              string
	Name            string
	Email           string
	OldPassword     string
	Password        string
	ConfirmPassword string
	Role            string
	RecoveryHash    string
	Errors          map[string]string
}

func (f *UserForm) Valid() bool {
	v := newValidator()

	v.Required("Name", f.Name)
	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidSignup() bool {
	v := newValidator()

	v.Required("Name", f.Name)
	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)
	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidLogin() bool {
	v := newValidator()

	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)
	v.Required("Password", f.Password)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidForgot() bool {
	v := newValidator()

	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidPasswordUpdate() bool {
	v := newValidator()

	v.Required("OldPassword", f.OldPassword)
	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidPassword() bool {
	v := newValidator()

	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}
