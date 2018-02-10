package db

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
}

type Users []*User

func (u *User) Create() error {
	conn, _ := GetConnection()

	hp, err := bcrypt.GenerateFromPassword([]byte(u.Password), viper.GetInt("cost"))
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, name, role, password, created_at, updated_at) VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, u.Email, u.Name, u.Role, string(hp))
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return ErrDuplicateEmail
		}

		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)

	return err
}

func (u *User) Update() error {
	conn, _ := GetConnection()

	stmt := `UPDATE users SET name = ?, email = ?, role = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, u.Name, u.Email, u.Role, u.ID)
	return err
}

func (u *User) UpdatePassword(pw string) error {
	conn, _ := GetConnection()
	pass, err := bcrypt.GenerateFromPassword([]byte(pw), viper.GetInt("cost"))

	stmt := `UPDATE users SET password = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err = conn.Exec(stmt, string(pass), u.ID)
	return err
}

func (u *User) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM users WHERE id = ?`
	_, err := conn.Exec(stmt, u.ID)
	return err
}

func (u *User) Get() error {
	conn, _ := GetConnection()

	snippet := `SELECT id, name, email, role FROM users WHERE`

	var row *sql.Row

	if u.ID != 0 {
		stmt := snippet + ` id = ?`
		row = conn.QueryRow(stmt, u.ID)
	} else {
		stmt := snippet + ` email = ?`
		row = conn.QueryRow(stmt, u.Email)
	}

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}

func GetUsers() (Users, error) {
	conn, _ := GetConnection()

	stmt := `SELECT id, name, role FROM users ORDER BY name`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := Users{}
	for rows.Next() {
		u := &User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) VerifyUser(pw string) error {
	conn, _ := GetConnection()

	var hp []byte
	row := conn.QueryRow("SELECT id, name, email, role, password FROM users WHERE email = ?", u.Email)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &hp)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hp, []byte(pw))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidCredentials
	} else if err != nil {
		return err
	}

	return nil
}

func (u *User) VerifyAndUpdatePassword(old string, new string) error {
	err := u.VerifyUser(old)
	if err != nil {
		return err
	}

	err = u.UpdatePassword(new)
	if err != nil {
		return err
	}
	return err
}

func (u *User) SetRecover(h string) error {
	conn, _ := GetConnection()

	stmt := `UPDATE users SET recovery_hash = ?, updated_at = UTC_TIMESTAMP() WHERE email = ?`
	_, err := conn.Exec(stmt, h, u.Email)
	return err
}

func (u *User) CheckRecover(h string) error {
	conn, _ := GetConnection()

	stmt := `SELECT email FROM users WHERE email = ? AND recovery_hash = ?`
	row := conn.QueryRow(stmt, u.Email, h)

	err := row.Scan(&u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return err
}

func (u *User) Recover(h string, pw string) error {
	err := u.Get()

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	err = u.CheckRecover(h)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	u.SetRecover("")
	err = u.UpdatePassword(pw)
	return err
}
