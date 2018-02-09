package flash

import (
	"net/http"
	"revelforce-admin/internal/platform/session"
)

type Msg struct {
	AlertType string
	Message   string
}

func Add(w http.ResponseWriter, r *http.Request, m string, t string) error {
	sesh := session.GetSession()
	s := sesh.Load(r)
	err := s.PutObject(w, "flash", &Msg{
		AlertType: t,
		Message:   m,
	})
	if err != nil {
		return err
	}
	return nil
}

func Fetch(w http.ResponseWriter, r *http.Request) (Msg, error) {
	sesh := session.GetSession()
	s := sesh.Load(r)
	f := Msg{}
	err := s.PopObject(w, "flash", &f)
	if err != nil {
		return f, err
	}
	return f, err
}
