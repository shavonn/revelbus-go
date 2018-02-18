package db

import (
	"database/sql"
)

type Settings struct {
	ID           int
	ContactBlurb string
	AboutBlurb   string
	AboutContent string
}

func (s *Settings) Create() error {
	conn, _ := GetConnection()

	stmt := `INSERT INTO settings (contact_blurb, about_blurb, about_content, created_at, updated_at) VALUES(?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, s.ContactBlurb, s.AboutBlurb, s.AboutContent)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = int(id)

	return nil
}

func (s *Settings) Update() error {
	conn, _ := GetConnection()

	stmt := `UPDATE settings SET contact_blurb = ?, about_blurb = ?, about_content = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, s.ContactBlurb, s.AboutBlurb, s.AboutContent, s.ID)
	return err
}

func (s *Settings) Get() error {
	conn, _ := GetConnection()

	stmt := `SELECT id, contact_blurb, about_blurb, about_content FROM settings WHERE id = ?`
	row := conn.QueryRow(stmt, s.ID)

	err := row.Scan(&s.ID, &s.ContactBlurb, &s.AboutBlurb, &s.AboutContent)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}
