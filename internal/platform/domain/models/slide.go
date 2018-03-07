package models

import (
	"database/sql"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/forms"

	"revelbus/pkg/database"
)

type Slide struct {
	ID     int
	Title  sql.NullString
	Blurb  sql.NullString
	Style  sql.NullString
	Order  sql.NullInt64
	Active bool
}

type Slides []*Slide

type SlideForm struct {
	ID     string
	Title  string
	Blurb  string
	Style  string
	Order  string
	Active bool

	Errors map[string]string
}

func (f *SlideForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Title", f.Title)
	v.Required("Blurb", f.Blurb)
	v.Required("Style", f.Style)
	v.ValidInt("Order", f.Order)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (s *Slide) Create() error {
	conn, _ := database.GetConnection()

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

func (s *Slide) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, title, blurb, style, sort_order, active FROM slides WHERE id = ?`
	err := conn.QueryRow(stmt, s.ID).Scan(&s.ID, &s.Title, &s.Blurb, &s.Style, &s.Order, &s.Active)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (s *Slide) Update() error {
	conn, _ := database.GetConnection()

	stmt := `UPDATE slides SET title = ?, blurb= ?, style = ?, sort_order = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, s.Title, s.Blurb, s.Style, s.Order, s.Active, s.ID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (s *Slide) Delete() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM slides WHERE id = ?`
	_, err := conn.Exec(stmt, s.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func FetchSlides() (*Slides, error) {
	conn, _ := database.GetConnection()

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

func FindActiveSlides() (*Slides, error) {
	conn, _ := database.GetConnection()

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
