package models

import (
	"database/sql"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/forms"

	"revelforce/pkg/database"

	"github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
)

type Gallery struct {
	ID     int
	Name   sql.NullString
	Folder sql.NullString

	Images Files
}

type Galleries []*Gallery

type GalleryForm struct {
	ID     string
	Name   string
	Folder string

	Errors map[string]string
}

func (f *GalleryForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Name", f.Name)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (g *Gallery) Create() error {
	conn, _ := database.GetConnection()

	if g.Folder.String == "" {
		g.Folder = sql.NullString{
			String: slug.Make(g.Name.String),
			Valid:  true,
		}
	}

	stmt := `INSERT INTO galleries (name, folder, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, g.Name, g.Folder)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	g.ID = int(id)

	return nil
}

func (g *Gallery) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT name, folder FROM galleries WHERE id = ?`
	err := conn.QueryRow(stmt, g.ID).Scan(&g.Name, &g.Folder)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	err = g.GetImages()
	return err
}

func (g *Gallery) Update() error {
	conn, _ := database.GetConnection()

	if g.Folder.String == "" {
		g.Folder = sql.NullString{
			String: slug.Make(g.Name.String),
			Valid:  true,
		}
	}

	stmt := `UPDATE galleries SET name = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, g.Name, g.ID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (g *Gallery) Delete() error {
	conn, _ := database.GetConnection()

	err := g.DeleteImages()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	stmt := `DELETE FROM galleries WHERE id = ?`
	_, err = conn.Exec(stmt, g.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func FetchGalleries() (*Galleries, error) {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, name FROM galleries`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	galleries := Galleries{}
	for rows.Next() {
		g := &Gallery{}
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, err
		}
		galleries = append(galleries, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &galleries, nil
}

func (g *Gallery) GetImages() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT f.id, f.name, f.thumb FROM galleries_images gi JOIN files f ON gi.file_id = f.id WHERE gi.gallery_id = ?`
	rows, err := conn.Query(stmt, g.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	files := Files{}
	for rows.Next() {
		f := &File{}
		err := rows.Scan(&f.ID, &f.Name, &f.Thumb)
		if err != nil {
			return err
		}
		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	g.Images = files

	return nil
}

func (g *Gallery) DeleteImages() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE f, gi FROM files f JOIN galleries_images gi ON gi.file_id = f.id WHERE gi.gallery_id = ?`
	_, err := conn.Exec(stmt, g.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (g *Gallery) AttachImage(fid string) error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO galleries_images (gallery_id, file_id, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	_, err := conn.Exec(stmt, g.ID, fid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return domain.ErrDuplicate
		}
	}
	return err
}

func (g *Gallery) DetachImage(fid string) error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM galleries_images WHERE gallery_id = ? AND file_id = ?`
	_, err := conn.Exec(stmt, g.ID, fid)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}
