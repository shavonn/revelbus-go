package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"time"
)

type File struct {
	ID      int
	Name    string
	Created time.Time
}

type Files []*File

func (f *File) Create() error {
	conn, _ := db.GetConnection()

	stmt := `INSERT INTO files (name, created_at, updated_at) VALUES(?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, f.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = int(id)
	return nil
}

func (f *File) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM files WHERE id = ?`
	_, err := conn.Exec(stmt, f.ID)
	return err
}

func (f *File) Get() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT name, created_at FROM files WHERE id = ?`
	err := conn.QueryRow(stmt, f.ID, f.Created).Scan(&f.Name)
	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	return err
}

func GetFiles() (*Files, error) {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, name, created_at FROM files ORDER BY name`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := Files{}
	for rows.Next() {
		f := &File{}
		err := rows.Scan(&f.ID, &f.Name, &f.Created)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &files, nil
}
