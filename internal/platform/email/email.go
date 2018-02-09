package email

import (
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type email struct {
	To      []string
	Subject string
	Text    string
	HTML    string
}

func send(e email) error {
	m := gomail.NewMessage()
	d := gomail.NewPlainDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.user"), viper.GetString("smtp.password"))

	m.SetBody("text/html", e.HTML)
	m.AddAlternative("text/plain", e.Text)
	m.SetHeaders(map[string][]string{
		"From":    []string{viper.GetString("from")},
		"To":      e.To,
		"Subject": []string{e.Subject},
	})

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}
