package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

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
