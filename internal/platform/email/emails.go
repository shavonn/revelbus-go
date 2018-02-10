package email

func NewPassword(e string, pw string) error {
	m := email{
		To: []string{
			e,
		},
		Subject: "Your New Password",
		Text:    "Your new password is: " + pw,
		HTML:    "<p>Your new password is: " + pw + "</p>",
	}

	err := send(m)
	return err
}

func RecoverAccount(e string, h string) error {
	m := email{
		To: []string{
			e,
		},
		Subject: "Password Recovery",
		Text:    "Click to reset password: /auth/recover/?email=" + e + "&hash=" + h,
		HTML:    "<p>Click to reset password: /auth/recover/?email=" + e + "&hash=" + h + "</p>",
	}

	err := send(m)
	return err
}
