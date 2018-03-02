package models

import (
	"database/sql"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/forms"
	"strconv"
	"time"

	"revelforce/pkg/database"

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
	ImageID      int
	GalleryID    int

	Image    *File
	Partners Vendors
	Venues   Vendors

	CalendarLinks map[string]string
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
	ImageID      int
	GalleryID    int
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
	conn, _ := database.GetConnection()

	slug := domain.GetSlug(t.Title, "trips")

	stmt := `INSERT INTO trips (title, slug, status, blurb, description, start, end, price, ticketing_url, notes, gallery_id, image_id, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, t.Title, slug, t.Status, t.Blurb, t.Description, t.Start, t.End, t.Price, t.TicketingURL, t.Notes, t.GalleryID, t.ImageID)
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

func (t *Trip) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT title, slug, status, blurb, description, start, end, price, ticketing_url, notes, image_id, gallery_id FROM trips WHERE id = ?`
	err := conn.QueryRow(stmt, t.ID).Scan(&t.Title, &t.Slug, &t.Status, &t.Blurb, &t.Description, &t.Start, &t.End, &t.Price, &t.TicketingURL, &t.Notes, &t.ImageID, &t.GalleryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return err
	}

	err = t.GetImage()
	if err != nil {
		return err
	}

	err = t.GetTripVendors()
	return err
}

func FindBySlug(s string) (*Trip, error) {
	conn, _ := database.GetConnection()
	t := &Trip{}

	stmt := `SELECT id, title, slug, status, blurb, description, start, end, price, ticketing_url, image_id, gallery_id FROM trips WHERE slug = ?`
	err := conn.QueryRow(stmt, s).Scan(&t.ID, &t.Title, &t.Slug, &t.Status, &t.Blurb, &t.Description, &t.Start, &t.End, &t.Price, &t.TicketingURL, &t.ImageID, &t.GalleryID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	err = t.GetImage()
	if err != nil {
		return nil, err
	}

	err = t.GetTripVendors()
	return t, err
}

func (t *Trip) Update() error {
	conn, _ := database.GetConnection()

	if t.Slug == "" {
		t.Slug = domain.GetSlug(t.Title, "trips")
	}

	stmt := `UPDATE trips SET title = ?, slug = ?, status = ?, blurb = ?, description = ?, start = ?, end = ?, price = ?, ticketing_url = ?, notes = ?, image_id = ?, gallery_id = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, t.Title, t.Slug, t.Status, t.Blurb, t.Description, t.Start, t.End, t.Price, t.TicketingURL, t.Notes, t.ImageID, t.GalleryID, t.ID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (t *Trip) Delete() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM trips WHERE id = ?`
	_, err := conn.Exec(stmt, t.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (t *Trip) GetBase() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT image_id FROM trips WHERE id = ?`
	err := conn.QueryRow(stmt, t.ID).Scan(&t.ImageID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func FetchTrips() (*Trips, error) {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, title, status, start, end FROM trips ORDER BY start, end`
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

func FindUpcomingTrips(limit int) (*Trips, error) {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, title, slug, start, end, image_id, blurb FROM trips WHERE (start > NOW() - INTERVAL 1 DAY) AND status = 'published' ORDER BY start, end`

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
		err := rows.Scan(&t.ID, &t.Title, &t.Slug, &t.Start, &t.End, &t.ImageID, &t.Blurb)
		if err != nil {
			return nil, err
		}

		err = t.GetImage()
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

func FindUpcomingTripsByMonth() (*GroupedTrips, error) {
	conn, _ := database.GetConnection()

	trips := make(GroupedTrips)

	stmt := `SELECT id, title, slug, start, end, image_id, blurb FROM trips WHERE (start > NOW() - INTERVAL 1 DAY) AND status = 'published' ORDER BY start, end`

	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &Trip{}
		err := rows.Scan(&t.ID, &t.Title, &t.Slug, &t.Start, &t.End, &t.ImageID, &t.Blurb)
		if err != nil {
			return nil, err
		}
		month := t.Start.Format("01")

		err = t.GetImage()
		if err != nil {
			return nil, err
		}

		trips[month] = append(trips[month], t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &trips, nil
}

func (t *Trip) GetTripPartners() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT v.id, v.name, v.brand_id, v.url FROM trips_partners tp JOIN vendors v ON tp.partner_id = v.id WHERE tp.trip_id = ? AND v.active = 1 ORDER BY name`
	rows, err := conn.Query(stmt, t.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	partners := Vendors{}
	for rows.Next() {
		p := &Vendor{}
		err := rows.Scan(&p.ID, &p.Name, &p.BrandID, &p.URL)
		if err != nil {
			return err
		}

		err = p.GetImage()
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

func (t *Trip) GetTripVenues() error {
	conn, _ := database.GetConnection()

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

func (t *Trip) GetImage() error {
	conn, _ := database.GetConnection()

	f := &File{}

	stmt := `SELECT f.id, f.name, f.thumb FROM trips t JOIN files f ON t.image_id = f.id WHERE t.id = ?`

	err := conn.QueryRow(stmt, t.ID).Scan(&f.ID, &f.Name, &f.Thumb)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			return err
		}
	}

	t.Image = f

	return nil
}

func (t *Trip) GetTripVendors() error {
	err := t.GetTripPartners()
	if err != nil {
		return err
	}

	err = t.GetTripVenues()
	return err
}

func (t *Trip) AttachVendor(r string, vid string) error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO trips_` + r + `s (trip_id, ` + r + `_id, created_at, updated_at) VALUES(?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	_, err := conn.Exec(stmt, t.ID, vid)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return domain.ErrDuplicate
		}
	}
	return err
}

func (t *Trip) DetachVendor(r string, vid string) error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM trips_` + r + `s WHERE trip_id = ? AND ` + r + `_id = ?`
	_, err := conn.Exec(stmt, t.ID, vid)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (t *Trip) SetVenueStatus(vid string, isPrimary bool) error {
	conn, _ := database.GetConnection()

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
