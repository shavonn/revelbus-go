package handlers

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/session"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func toInt(id string) int {
	i, _ := strconv.Atoi(id)
	return i
}

func randomString(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(result)
}

func uploadFile(w http.ResponseWriter, r *http.Request, fieldName string, fldr string) (string, error) {
	file, h, err := r.FormFile(fieldName)
	if err == http.ErrMissingFile {
		return "", nil
	} else if err != nil {
		return "", err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	dst, err := os.Create(filepath.Join(viper.GetString("files.static")+fldr, h.Filename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = dst.Write(bs)
	if err != nil {
		return "", err
	}
	return h.Filename, err
}

func deleteFile(fn string) error {
	err := os.Remove(viper.GetString("files.static") + fn)
	return err
}

func setUserSession(w http.ResponseWriter, r *http.Request, u *db.User) error {
	sesh := session.GetSession()

	s := sesh.Load(r)
	err := s.PutObject(w, "AuthUser", u)
	if err != nil {
		return err
	}
	return nil
}

func removeUserSession(w http.ResponseWriter, r *http.Request) error {
	sesh := session.GetSession()

	s := sesh.Load(r)
	err := s.Remove(w, "AuthUser")
	return err
}

func loggedIn(r *http.Request) (*db.User, error) {
	sesh := session.GetSession()

	s := sesh.Load(r)
	loggedIn, err := s.Exists("AuthUser")
	if err != nil {
		return nil, err
	}

	if !loggedIn {
		return nil, nil
	}

	u := &db.User{}
	err = s.GetObject("AuthUser", u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
