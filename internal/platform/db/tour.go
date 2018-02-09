package db

import (
	"database/sql"
	"revelforce-admin/internal/platform/forms"
	"time"
)

type Tour struct {
	ID           int
	Name         string
	Slug         sql.NullString
	Status       string
	Description  sql.NullString
	Start        time.Time
	End          time.Time
	TicketingURL sql.NullString
	Notes        sql.NullString
}

func Create(f forms.TourForm) (int, error) {
	conn, _ := GetConnection()

	start := toTime(f.Start)
	end := toTime(f.End)

	slug := getSlug(f.Name, "tours")

	stmt := `INSERT INTO tours (name, slug, status, description, start, end, ticketing_url, notes, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := conn.Exec(stmt, f.Name, slug, f.Status, f.Description, start, end, f.TicketingURL, f.Notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
