package models

import (
	"database/sql"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/forms"

	"revelbus/pkg/database"
)

type Settings struct {
	ID                int
	ContactBlurb      sql.NullString
	AboutBlurb        sql.NullString
	AboutContent      sql.NullString
	HomeGalleryID     sql.NullInt64
	HomeGalleryActive bool
}

type SettingsForm struct {
	ID                string
	ContactBlurb      string
	AboutBlurb        string
	AboutContent      string
	HomeGalleryID     int
	HomeGalleryActive bool

	Errors map[string]string
}

func (f *SettingsForm) Valid() bool {
	v := forms.NewValidator()

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (s *Settings) Create() error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO settings (contact_blurb, about_blurb, about_content, home_gallery, home_gallery_active, created_at, updated_at) VALUES(?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, s.ContactBlurb, s.AboutBlurb, s.AboutContent, s.HomeGalleryID, s.HomeGalleryActive)
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

func (s *Settings) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, contact_blurb, about_blurb, about_content, home_gallery, home_gallery_active FROM settings WHERE id = ?`
	err := conn.QueryRow(stmt, s.ID).Scan(&s.ID, &s.ContactBlurb, &s.AboutBlurb, &s.AboutContent, &s.HomeGalleryID, &s.HomeGalleryActive)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (s *Settings) Update() error {
	conn, _ := database.GetConnection()

	stmt := `UPDATE settings SET contact_blurb = ?, about_blurb = ?, about_content = ?, home_gallery = ?, home_gallery_active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, s.ContactBlurb, s.AboutBlurb, s.AboutContent, s.HomeGalleryID, s.HomeGalleryActive, s.ID)
	if err == sql.ErrNoRows {
		return s.Create()
	}
	return err
}
