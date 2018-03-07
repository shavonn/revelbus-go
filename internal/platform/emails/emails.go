package emails

import (
	"revelbus/internal/platform/forms"
	"revelbus/pkg/email"
)

func NewPassword(e string, pw string) error {
	m := email.Email{
		To: []string{
			e,
		},
		Subject: "Your New Password",
		Text:    "Your new password is: " + pw,
		HTML:    "<p>Your new password is: " + pw + "</p>",
	}

	err := email.Send(m)
	return err
}

func RecoverAccount(e string, h string) error {
	m := email.Email{
		To: []string{
			e,
		},
		Subject: "Password Recovery",
		Text:    "Click to reset password: /auth/recover/?email=" + e + "&hash=" + h,
		HTML:    "<p>Click to reset password: /auth/recover/?email=" + e + "&hash=" + h + "</p>",
	}

	err := email.Send(m)
	return err
}

func ContactEmail(f *forms.ContactForm) error {
	m := email.Email{
		Subject: "Revel Bus Contact Form",
		Text:    "Name: " + f.Name + " \nEmail: " + f.Email + " \nPhone: " + f.Phone + " \nMessage: " + f.Message,
		HTML:    "Name: " + f.Name + " \nEmail: " + f.Email + " \nPhone: " + f.Phone + " \nMessage: " + f.Message,
	}

	err := email.Send(m)
	return err
}
