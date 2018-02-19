package utils

import (
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"revelforce/internal/platform/db/models"

	"github.com/kennygrant/sanitize"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
)

func UploadFile(w http.ResponseWriter, r *http.Request, fieldName string, folder string, makeThumb bool) ([]*models.File, error) {
	uploaded := []*models.File{}
	uploadDir := filepath.Join(viper.GetString("files.static") + folder)

	err := r.ParseMultipartForm(100000)
	if err != nil {
		return uploaded, err
	}

	m := r.MultipartForm
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	files := m.File[fieldName]
	for i := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return uploaded, err
		}

		// new file
		f := &models.File{}

		// clean up file name
		fn := sanitize.Name(files[i].Filename)

		dst, err := os.Create(filepath.Join(uploadDir, fn))
		if err != nil {
			return uploaded, err
		}

		if _, err := io.Copy(dst, file); err != nil {
			return uploaded, err
		}

		f.Name = filepath.Join(folder, fn)
		err = f.Create()
		if err != nil {
			return uploaded, err
		}

		if makeThumb {
			rn := "thumb_" + fn

			file, err := os.Open(filepath.Join(uploadDir, fn))
			if err != nil {
				log.Fatal(err)
			}

			img, err := jpeg.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			m := resize.Resize(400, 0, img, resize.Lanczos3)

			out, err := os.Create(filepath.Join(uploadDir, rn))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			jpeg.Encode(out, m, nil)
		}

		uploaded = append(uploaded, f)
	}

	return uploaded, err
}

func DeleteFile(f *models.File) error {
	err := os.Remove(viper.GetString("files.static") + f.Name)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if f.Thumb != "" {
		err := os.Remove(viper.GetString("files.static") + f.Thumb)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
