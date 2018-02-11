package db

import (
	"database/sql"
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
	Brand   string
}

type Vendors []*Vendor

func (v *Vendor) Create() error {
	conn, _ := GetConnection()

	stmt := `INSERT INTO vendors (name, address, city, state, zip, phone, email, url, notes, brand, active, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.Brand, v.Active)
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
	conn, _ := GetConnection()

	stmt := `UPDATE vendors SET name = ?, address = ?, city = ?, state = ?, zip = ?, phone = ?, email = ?, url = ?, notes = ?, brand = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.Brand, v.Active, v.ID)
	return err
}

func (v *Vendor) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM vendors WHERE id = ?`
	_, err := conn.Exec(stmt, v.ID)
	return err
}

func (v *Vendor) Get() error {
	conn, _ := GetConnection()

	stmt := `SELECT id, name, address, city, state, zip, phone, email, url, notes, brand, active FROM vendors WHERE id = ?`
	row := conn.QueryRow(stmt, v.ID)

	err := row.Scan(&v.ID, &v.Name, &v.Address, &v.City, &v.State, &v.Zip, &v.Phone, &v.Email, &v.URL, &v.Notes, &v.Brand, &v.Active)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}

func GetVendors(oa bool) (Vendors, error) {
	conn, _ := GetConnection()

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

	return vendors, nil
}
