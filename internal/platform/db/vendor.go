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
	Logo    string
	Notes   string
}

type Vendors []*Vendor

func (v *Vendor) Create() (int, error) {
	conn, _ := GetConnection()

	stmt := `INSERT INTO vendors (name, address, city, state, zip, phone, email, url, notes, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (v *Vendor) Update() error {
	conn, _ := GetConnection()

	stmt := `UPDATE vendors SET name = ?, address = ?, city = ?, state = ?, zip = ?, phone = ?, email = ?, url = ?, notes = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, v.Name, v.Address, v.City, v.State, v.Zip, v.Phone, v.Email, v.URL, v.Notes, v.ID)
	return err
}

func (t *Vendor) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM vendors WHERE id = ?`
	_, err := conn.Exec(stmt, t.ID)
	return err
}

func GetVendorByID(id string) (*Vendor, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, name, address, city, state, zip, phone, email, url, notes FROM vendors WHERE id = ?`
	row := conn.QueryRow(stmt, id)

	v := &Vendor{}
	err := row.Scan(&v.ID, &v.Name, &v.Address, &v.City, &v.State, &v.Zip, &v.Phone, &v.Email, &v.URL, &v.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return v, nil
}

func GetVendors() (Vendors, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, name FROM vendors ORDER BY name`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vendors := Vendors{}
	for rows.Next() {
		e := &Vendor{}
		err := rows.Scan(&e.ID, &e.Name)
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
