package utils

import (
	"net/http"
	"revelbus/internal/platform/domain/models"
	"revelbus/pkg/sessions"
)

func IsAuthenticated(r *http.Request) (*models.User, error) {
	sesh := sessions.GetSession()

	s := sesh.Load(r)
	loggedIn, err := s.Exists("AuthUser")
	if err != nil {
		return nil, err
	}

	if !loggedIn {
		return nil, nil
	}

	u := &models.User{}
	err = s.GetObject("AuthUser", u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func SetUserSession(w http.ResponseWriter, r *http.Request, u *models.User) error {
	sesh := sessions.GetSession()

	s := sesh.Load(r)
	err := s.PutObject(w, "AuthUser", u)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserSession(w http.ResponseWriter, r *http.Request) error {
	sesh := sessions.GetSession()

	s := sesh.Load(r)
	err := s.Remove(w, "AuthUser")
	return err
}
