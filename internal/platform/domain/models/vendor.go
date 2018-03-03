package models

import (
	"database/sql"
	"revelforce/internal/platform/domain"
	"revelforce/internal/platform/forms"
	"revelforce/pkg/database"

	"github.com/go-sql-driver/mysql"
)

type Vendor struct {
	ID      int
	Name    sql.NullString
	Address sql.NullString
	City    sql.NullString
	State   sql.NullString
	Zip     sql.NullString
	Phone   sql.NullString
	Email   sql.NullString
	URL     sql.NullString
	Notes   sql.NullString
	Primary bool
	Active  bool
	BrandID int

	Brand *File
}

type Vendors []*Vendor

type VendorForm struct {
	ID      string
	Name    string
	Address string
	City    string
	State   string
	Zip     string
	Phone   string
	Email   string
	URL     string
	Notes   string
	BrandID int
	Active  bool
	Brand   string
	Errors  map[string]string
}

func (f *VendorForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Name", f.Name)
	v.ValidEmail("Email", f.Email)
	v.ValidURL("URL", f.URL)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (v *Vendor) Create() error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO vendors (name, address, city, state, zip, phone, email, url, notes, brand_id, active, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.BrandID, v.Active)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	v.ID = int(id)

	return nil
}

func (v *Vendor) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, name, address, city, state, zip, phone, email, url, notes, brand_id, active FROM vendors WHERE id = ?`
	err := conn.QueryRow(stmt, v.ID).Scan(&v.ID, &v.Name, &v.Address, &v.City, &v.State, &v.Zip, &v.Phone, &v.Email, &v.URL, &v.Notes, &v.BrandID, &v.Active)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	err = v.GetImage()
	if err != nil {
		return err
	}

	return err
}

func (v *Vendor) Update() error {
	conn, _ := database.GetConnection()

	stmt := `UPDATE vendors SET name = ?, address = ?, city = ?, state = ?, zip = ?, phone = ?, email = ?, url = ?, notes = ?, brand_id = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.BrandID, v.Active, v.ID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (v *Vendor) Delete() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM vendors WHERE id = ?`
	_, err := conn.Exec(stmt, v.ID)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1451 {
			return domain.ErrCannotDelete
		} else if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
	}
	return err
}

func (v *Vendor) GetBase() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT brand_id FROM vendors WHERE id = ?`
	err := conn.QueryRow(stmt, v.ID).Scan(&v.BrandID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (v *Vendor) GetImage() error {
	conn, _ := database.GetConnection()

	f := &File{}

	stmt := `SELECT f.id, f.name, f.thumb FROM vendors v JOIN files f ON v.brand_id = f.id WHERE v.id = ?`

	err := conn.QueryRow(stmt, v.ID).Scan(&f.ID, &f.Name, &f.Thumb)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	v.Brand = f

	return nil
}

func FetchVendors(oa bool) (*Vendors, error) {
	conn, _ := database.GetConnection()

	var stmt string

	if oa {
		stmt = `SELECT id, name, active FROM vendors WHERE active = 1 ORDER BY active DESC, name`
	} else {
		stmt = `SELECT id, name, active FROM vendors ORDER BY active DESC, name`
	}

	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vendors := Vendors{}
	for rows.Next() {
		e := &Vendor{}
		err := rows.Scan(&e.ID, &e.Name, &e.Active)
		if err != nil {
			return nil, err
		}
		vendors = append(vendors, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &vendors, nil
}
