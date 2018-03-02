package models

import (
	"database/sql"
	"revelforce/internal/platform/domain"
	"revelforce/pkg/database"
	"time"

	"github.com/go-sql-driver/mysql"
)

type File struct {
	ID      int
	Name    string
	Thumb   string
	Created time.Time
}

type Files []*File

func (f *File) Create() error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO files (name, thumb, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, f.Name, f.Thumb)
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

func (f *File) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT name, thumb, created_at FROM files WHERE id = ?`
	err := conn.QueryRow(stmt, f.ID).Scan(&f.Name, &f.Thumb, &f.Created)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (f *File) Delete() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM files WHERE id = ?`
	_, err := conn.Exec(stmt, f.ID)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1451 {
			return domain.ErrCannotDelete
		} else if err == sql.ErrNoRows {
			return nil
		}
	}

	return err
}

func FetchFiles() (*Files, error) {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, name, thumb, created_at FROM files ORDER BY name`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := Files{}
	for rows.Next() {
		f := &File{}
		err := rows.Scan(&f.ID, &f.Name, &f.Thumb, &f.Created)
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
