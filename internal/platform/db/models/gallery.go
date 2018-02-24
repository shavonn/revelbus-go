package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/forms"

	"github.com/go-sql-driver/mysql"
)

type Gallery struct {
	ID     int
	Name   string
	Images Files
}

type Galleries []*Gallery

type GalleryForm struct {
	ID     string
	Name   string
	Errors map[string]string
}

func (f *GalleryForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Name", f.Name)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (g *Gallery) Create() error {
	conn, _ := db.GetConnection()

	stmt := `INSERT INTO galleries (name, created_at, updated_at) VALUES(?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, g.Name)
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

func (g *Gallery) Update() error {
	conn, _ := db.GetConnection()

	stmt := `UPDATE galleries SET name = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, g.Name, g.ID)
	return err
}

func (g *Gallery) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM galleries WHERE id = ?`
	_, err := conn.Exec(stmt, g.ID)
	return err
}

func (g *Gallery) Get() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT name FROM galleries WHERE id = ?`
	err := conn.QueryRow(stmt, g.ID).Scan(&g.Name)
	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	err = g.GetFiles()
	return err
}

func (g *Gallery) GetFiles() error {
	conn, _ := db.GetConnection()

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

func GetGalleries() (*Galleries, error) {
	conn, _ := db.GetConnection()

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

func (g *Gallery) AttachImage(fid string) error {
	conn, _ := db.GetConnection()

	stmt := `INSERT INTO galleries_images (gallery_id, file_id, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	_, err := conn.Exec(stmt, g.ID, fid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return db.ErrDuplicate
		}
	}
	return err
}

func (g *Gallery) DetachImage(fid string) error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM galleries_images WHERE gallery_id = ? AND file_id = ?`
	_, err := conn.Exec(stmt, g.ID, fid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return db.ErrDuplicate
		}
	}
	return err
}
