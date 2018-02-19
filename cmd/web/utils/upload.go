package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/viper"
)

func UploadFile(w http.ResponseWriter, r *http.Request, fieldName string, fldr string) ([]string, error) {
	uploaded := []string{}
	dir := filepath.Join(viper.GetString("files.static") + fldr)

	err := r.ParseMultipartForm(100000)
	if err != nil {
		return uploaded, err
	}

	m := r.MultipartForm
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	files := m.File[fieldName]
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return uploaded, err
		}

		fn := sanitize.Name(files[i].Filename)

		dst, err := os.Create(filepath.Join(dir, fn))
		defer dst.Close()
		if err != nil {
			return uploaded, err
		}

		if _, err := io.Copy(dst, file); err != nil {
			return uploaded, err
		}

		if fldr != "" {
			fn = fldr + "/" + fn
		}

		uploaded = append(uploaded, fn)
	}

	return uploaded, err
}

func DeleteFile(fn string) error {
	err := os.Remove(viper.GetString("files.static") + fn)

	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
