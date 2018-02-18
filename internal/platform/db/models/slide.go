package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/forms"
)

type Slide struct {
	ID     int
	Title  string
	Blurb  string
	Style  string
	Active bool
	Order  int
}

type Slides []*Slide

type SlideForm struct {
	ID     string
	Title  string
	Blurb  string
	Style  string
	Active bool
	Order  string
	Errors map[string]string
}

func (f *SlideForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Title", f.Title)
	v.Required("Blurb", f.Blurb)
	v.Required("Style", f.Style)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (s *Slide) Create() error {
	conn, _ := db.GetConnection()

	stmt := `INSERT INTO slides (title, blurb, style, sort_order, active, created_at, updated_at) VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, s.Title, s.Blurb, s.Style, s.Order, s.Active)
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

func (s *Slide) Update() error {
	conn, _ := db.GetConnection()

	stmt := `UPDATE slides SET title = ?, blurb= ?, style = ?, sort_order = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, s.Title, s.Blurb, s.Style, s.Order, s.Active, s.ID)
	return err
}

func (s *Slide) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM slides WHERE id = ?`
	_, err := conn.Exec(stmt, s.ID)
	return err
}

func (s *Slide) Get() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, title, blurb, style, sort_order, active FROM slides WHERE id = ?`
	err := conn.QueryRow(stmt, s.ID).Scan(&s.ID, &s.Title, &s.Blurb, &s.Style, &s.Order, &s.Active)
	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	return err
}

func GetSlides() (*Slides, error) {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, title, blurb, style, sort_order, active FROM slides ORDER BY sort_order`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slides := Slides{}
	for rows.Next() {
		s := &Slide{}
		err := rows.Scan(&s.ID, &s.Title, &s.Blurb, &s.Style, &s.Order, &s.Active)
		if err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &slides, nil
}

func GetActiveSlides() (*Slides, error) {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, title, blurb, style, sort_order FROM slides WHERE active = 1 ORDER BY sort_order`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slides := Slides{}
	for rows.Next() {
		s := &Slide{}
		err := rows.Scan(&s.ID, &s.Title, &s.Blurb, &s.Style, &s.Order)
		if err != nil {
			return nil, err
		}
		slides = append(slides, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &slides, nil
}
