package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/forms"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Trip struct {
	ID           int
	Status       string
	Slug         string
	Title        string
	Blurb        string
	Description  string
	Start        time.Time
	End          time.Time
	Price        string
	TicketingURL string
	Notes        string
	Image        string

	Partners Vendors
	Venues   Vendors
}

type Trips []*Trip

type GroupedTrips map[string][]*Trip

type TripForm struct {
	ID           string
	Title        string
	Slug         string
	Status       string
	Blurb        string
	Description  string
	Start        string
	End          string
	Price        string
	TicketingURL string
	Notes        string
	Image        string
	Errors       map[string]string
}

func (f *TripForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Title", f.Title)
	v.ValidSlug("Slug", f.Slug)
	v.ValidDateTime("Start", f.Start)
	v.ValidDateTime("End", f.End)
	v.ValidDateTimeRange("End", f.Start, f.End)
	v.ValidURL("TicketingURL", f.TicketingURL)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (t *Trip) Create() error {
	conn, _ := db.GetConnection()

	slug := db.GetSlug(t.Title, "trips")

	stmt := `INSERT INTO trips (title, slug, status, blurb, description, start, end, price, ticketing_url, notes, image, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, t.Title, slug, t.Status, t.Blurb, t.Description, t.Start, t.End, t.Price, t.TicketingURL, t.Notes, t.Image)
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
	conn, _ := db.GetConnection()

	if t.Slug == "" {
		t.Slug = db.GetSlug(t.Title, "trips")
	}

	stmt := `UPDATE trips SET title = ?, slug = ?, status = ?, blurb = ?, description = ?, start = ?, end = ?, price = ?, ticketing_url = ?, notes = ?, image = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, t.Title, t.Slug, t.Status, t.Blurb, t.Description, t.Start, t.End, t.Price, t.TicketingURL, t.Notes, t.Image, t.ID)
	return err
}

func (t *Trip) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM trips WHERE id = ?`
	_, err := conn.Exec(stmt, t.ID)
	return err
}

func (t *Trip) Get() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT title, slug, status, blurb, description, start, end, price, ticketing_url, notes, image FROM trips WHERE id = ?`
	err := conn.QueryRow(stmt, t.ID).Scan(&t.Title, &t.Slug, &t.Status, &t.Blurb, &t.Description, &t.Start, &t.End, &t.Price, &t.TicketingURL, &t.Notes, &t.Image)
	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	err = t.GetVendors()

	return err
}

func GetBySlug(s string) (*Trip, error) {
	conn, _ := db.GetConnection()
	t := &Trip{}

	stmt := `SELECT id, title, slug, status, blurb, description, start, end, price, ticketing_url, image FROM trips WHERE slug = ?`
	err := conn.QueryRow(stmt, s).Scan(&t.ID, &t.Title, &t.Slug, &t.Status, &t.Blurb, &t.Description, &t.Start, &t.End, &t.Price, &t.TicketingURL, &t.Image)
	if err == sql.ErrNoRows {
		return nil, db.ErrNotFound
	}

	err = t.GetVendors()
	if err != nil {
		return nil, err
	}

	return t, err
}

func GetTrips() (*Trips, error) {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, title, status, start, end, blurb FROM trips ORDER BY end, start`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := Trips{}
	for rows.Next() {
		t := &Trip{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.Start, &t.End, &t.Blurb)
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

func GetUpcomingTrips(limit int) (*Trips, error) {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, title, slug, start, end, image, blurb FROM trips WHERE (start > NOW() - INTERVAL 1 DAY) AND status = 'published' ORDER BY start, end`

	if limit > 0 {
		stmt = stmt + ` LIMIT ` + strconv.Itoa(limit)
	}

	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := Trips{}
	for rows.Next() {
		t := &Trip{}
		err := rows.Scan(&t.ID, &t.Title, &t.Slug, &t.Start, &t.End, &t.Image, &t.Blurb)
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

func GetUpcomingTripsByMonth() (*GroupedTrips, error) {
	conn, _ := db.GetConnection()

	trips := make(GroupedTrips)

	stmt := `SELECT id, title, slug, start, end, image, blurb FROM trips WHERE (start > NOW() - INTERVAL 1 DAY) AND status = 'published' ORDER BY start, end`

	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &Trip{}
		err := rows.Scan(&t.ID, &t.Title, &t.Slug, &t.Start, &t.End, &t.Image, &t.Blurb)
		if err != nil {
			return nil, err
		}
		month := t.Start.Format("01")
		trips[month] = append(trips[month], t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &trips, nil
}

func (t *Trip) GetPartners() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT v.id, v.name FROM trips_partners tp JOIN vendors v ON tp.partner_id = v.id WHERE tp.trip_id = ? AND v.active = 1 ORDER BY name`
	rows, err := conn.Query(stmt, t.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	partners := Vendors{}
	for rows.Next() {
		p := &Vendor{}
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			return err
		}
		partners = append(partners, p)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	t.Partners = partners

	return nil
}

func (t *Trip) GetVenues() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT v.id, v.name, v.address, v.city, v.state, v.zip, v.phone, tv.is_primary FROM trips_venues tv JOIN vendors v ON tv.venue_id = v.id WHERE tv.trip_id = ? AND v.active = 1 ORDER BY name`
	rows, err := conn.Query(stmt, t.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	venues := Vendors{}
	for rows.Next() {
		v := &Vendor{}
		err := rows.Scan(&v.ID, &v.Name, &v.Address, &v.City, &v.State, &v.Zip, &v.Phone, &v.Primary)
		if err != nil {
			return err
		}
		venues = append(venues, v)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	t.Venues = venues

	return nil
}

func (t *Trip) GetVendors() error {
	err := t.GetPartners()
	if err != nil {
		return err
	}

	err = t.GetVenues()
	return err
}

func (t *Trip) AddVendor(r string, vid string) error {
	conn, _ := db.GetConnection()

	stmt := `INSERT INTO trips_` + r + `s (trip_id, ` + r + `_id, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	_, err := conn.Exec(stmt, t.ID, vid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return db.ErrDuplicate
		}
	}
	return err
}

func (t *Trip) RemoveVendor(r string, vid string) error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM trips_` + r + `s WHERE trip_id = ? AND ` + r + `_id = ?`
	_, err := conn.Exec(stmt, t.ID, vid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return db.ErrDuplicate
		}
	}
	return err
}

func (t *Trip) SetVenueStatus(vid string, isPrimary bool) error {
	conn, _ := db.GetConnection()

	if isPrimary {
		stmt := `UPDATE trips_venues SET is_primary = false, updated_at = UTC_TIMESTAMP() WHERE trip_id = ? AND is_primary = true`
		_, err := conn.Exec(stmt, t.ID)
		if err != nil {
			return err
		}
	}

	stmt := `UPDATE trips_venues SET is_primary = ?, updated_at = UTC_TIMESTAMP() WHERE venue_id = ? AND trip_id = ?`
	_, err := conn.Exec(stmt, isPrimary, vid, t.ID)
	return err
}
