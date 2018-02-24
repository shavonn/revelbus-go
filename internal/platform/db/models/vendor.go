package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/forms"

	"github.com/go-sql-driver/mysql"
)

type Vendor struct {
	ID      int
	Name    string
	Address string
	City    string
	State   string
	Zip     string
	Phone   string
	Email   string
	URL     string
	Notes   string
	Active  bool
	Primary bool
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
	conn, _ := db.GetConnection()

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

func (v *Vendor) Update() error {
	conn, _ := db.GetConnection()

	stmt := `UPDATE vendors SET name = ?, address = ?, city = ?, state = ?, zip = ?, phone = ?, email = ?, url = ?, notes = ?, brand_id = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.BrandID, v.Active, v.ID)
	return err
}

func (v *Vendor) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM vendors WHERE id = ?`
	_, err := conn.Exec(stmt, v.ID)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1451 {
			return db.ErrCannotDelete
		}
	}
	return err
}

func (v *Vendor) Get() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT id, name, address, city, state, zip, phone, email, url, notes, brand_id, active FROM vendors WHERE id = ?`
	err := conn.QueryRow(stmt, v.ID).Scan(&v.ID, &v.Name, &v.Address, &v.City, &v.State, &v.Zip, &v.Phone, &v.Email, &v.URL, &v.Notes, &v.BrandID, &v.Active)
	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	err = v.GetFile()
	if err != nil {
		return err
	}

	return err
}

func (v *Vendor) GetBase() error {
	conn, _ := db.GetConnection()

	stmt := `SELECT brand_id FROM vendors WHERE id = ?`
	err := conn.QueryRow(stmt, v.ID).Scan(&v.BrandID)
	return err
}

func (v *Vendor) GetFile() error {
	conn, _ := db.GetConnection()

	f := &File{}

	stmt := `SELECT f.id, f.name, f.thumb, f.created_at FROM vendors v JOIN files f ON v.brand_id = f.id WHERE v.id = ?`

	err := conn.QueryRow(stmt, v.ID).Scan(&f.ID, &f.Name, &f.Thumb, &f.Created)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	v.Brand = f

	return nil
}

func GetVendors(oa bool) (*Vendors, error) {
	conn, _ := db.GetConnection()

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
