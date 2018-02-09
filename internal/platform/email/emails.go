package email

func TestEmail(e string) error {
	reset := email{
		To: []string{
			e,
		},
		Subject: "This is a test!",
		Text:    "This email is a test email.",
		HTML:    "<p>This email is a test email.</p>",
	}

	err := send(reset)
	return err
}
