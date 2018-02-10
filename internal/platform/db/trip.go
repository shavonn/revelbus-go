package db

import (
	"database/sql"
	"time"
)

type Trip struct {
	ID           int
	Title        string
	Slug         string
	Status       string
	Description  string
	Start        time.Time
	End          time.Time
	TicketingURL string
	Notes        string
	Image        string
}

type Trips []*Trip

func (t *Trip) Create() error {
	conn, _ := GetConnection()

	slug := getSlug(t.Title, "trips")

	stmt := `INSERT INTO trips (title, slug, status, description, start, end, ticketing_url, notes, image, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, t.Title, slug, t.Status, t.Description, t.Start, t.End, t.TicketingURL, t.Notes, t.Image)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = int(id)

	return nil
}

func (t *Trip) Update() error {
	conn, _ := GetConnection()

	if t.Slug == "" {
		t.Slug = getSlug(t.Title, "trips")
	}

	stmt := `UPDATE trips SET title = ?, slug = ?, status = ?, description = ?, start = ?, end = ?, ticketing_url = ?, notes = ?, image = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, t.Title, t.Slug, t.Status, t.Description, t.Start, t.End, t.TicketingURL, t.Notes, t.Image, t.ID)
	return err
}

func (t *Trip) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM trips WHERE id = ?`
	_, err := conn.Exec(stmt, t.ID)
	return err
}

func (t *Trip) Get() error {
	conn, _ := GetConnection()

	stmt := `SELECT title, slug, status, description, start, end, ticketing_url, notes, image FROM trips WHERE id = ?`
	row := conn.QueryRow(stmt, t.ID)

	err := row.Scan(&t.Title, &t.Slug, &t.Status, &t.Description, &t.Start, &t.End, &t.TicketingURL, &t.Notes, &t.Image)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}

func GetTrips() (*Trips, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, title, status, start, end FROM trips ORDER BY end, start`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := Trips{}
	for rows.Next() {
		t := &Trip{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.Start, &t.End)
		if err != nil {
			return nil, err
		}
		trips = append(trips, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &trips, nil
}
