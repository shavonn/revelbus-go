package utils

import (
	"image"
	"image/jpeg"
	"image/png"
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

		if makeThumb {
			rn := "thumb_" + fn
			ext := filepath.Ext(fn)

			file, err := os.Open(filepath.Join(uploadDir, fn))
			if err != nil {
				log.Fatal(err)
			}

			var img image.Image

			if ext == ".png" {
				img, err = png.Decode(file)
			} else {
				img, err = jpeg.Decode(file)
			}
			if err != nil {
				log.Fatal(err)
			}
			file.Close()

			m := resize.Resize(200, 0, img, resize.NearestNeighbor)

			out, err := os.Create(filepath.Join(uploadDir, rn))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			if ext == ".png" {
				png.Encode(out, m)
			} else {
				jpeg.Encode(out, m, nil)
			}

			f.Thumb = filepath.Join(folder, rn)
		}

		err = f.Create()
		if err != nil {
			return uploaded, err
		}

		uploaded = append(uploaded, f)
	}

	return uploaded, err
}

func DeleteFile(f *models.File) error {
	err := f.Get()
	if err != nil {
		return err
	}

	err = f.Delete()
	if err != nil {
		return err
	}

	err = os.Remove(viper.GetString("files.static") + f.Name)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if len(f.Thumb) > 0 {
		err := os.Remove(viper.GetString("files.static") + f.Thumb)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
