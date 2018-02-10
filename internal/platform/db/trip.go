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
}

type Trips []*Trip

func (t *Trip) Create() (int, error) {
	conn, _ := GetConnection()

	slug := getSlug(t.Title, "trips")

	stmt := `INSERT INTO trips (title, slug, status, description, start, end, ticketing_url, notes, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, t.Title, slug, t.Status, t.Description, t.Start, t.End, t.TicketingURL, t.Notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (t *Trip) Update() error {
	conn, _ := GetConnection()

	if t.Slug == "" {
		t.Slug = getSlug(t.Title, "trips")
	}

	stmt := `UPDATE trips SET title = ?, slug = ?, status = ?, description = ?, start = ?, end = ?, ticketing_url = ?, notes = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, t.Title, t.Slug, t.Status, t.Description, t.Start, t.End, t.TicketingURL, t.Notes, t.ID)
	return err
}

func (t *Trip) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM trips WHERE id = ?`
	_, err := conn.Exec(stmt, t.ID)
	return err
}

func GetTripByID(id string) (*Trip, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, title, slug, status, description, start, end, ticketing_url, notes FROM trips WHERE id = ?`
	row := conn.QueryRow(stmt, id)

	t := &Trip{}
	err := row.Scan(&t.ID, &t.Title, &t.Slug, &t.Status, &t.Description, &t.Start, &t.End, &t.TicketingURL, &t.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

func GetTrips() (Trips, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, title, status, start, end FROM trips ORDER BY end, start`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := Trips{}
	for rows.Next() {
		e := &Trip{}
		err := rows.Scan(&e.ID, &e.Title, &e.Status, &e.Start, &e.End)
		if err != nil {
			return nil, err
		}
		trips = append(trips, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}
