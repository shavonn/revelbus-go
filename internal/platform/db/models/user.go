package models

import (
	"database/sql"
	"revelforce/internal/platform/db"
	"revelforce/internal/platform/forms"

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

type UserForm struct {
	ID              string
	Name            string
	Email           string
	OldPassword     string
	Password        string
	ConfirmPassword string
	Role            string
	RecoveryHash    string
	Errors          map[string]string
}

func (f *UserForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Name", f.Name)
	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidSignup() bool {
	v := forms.NewValidator()

	v.Required("Name", f.Name)
	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)
	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidLogin() bool {
	v := forms.NewValidator()

	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)
	v.Required("Password", f.Password)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidForgot() bool {
	v := forms.NewValidator()

	v.Required("Email", f.Email)
	v.ValidEmail("Email", f.Email)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidPasswordUpdate() bool {
	v := forms.NewValidator()

	v.Required("OldPassword", f.OldPassword)
	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *UserForm) ValidPassword() bool {
	v := forms.NewValidator()

	v.Required("Password", f.Password)
	if f.Password != f.ConfirmPassword {
		v.Errors["Password"] = "Passwords must match."
	}

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (u *User) Create() error {
	conn, _ := db.GetConnection()

	hp, err := bcrypt.GenerateFromPassword([]byte(u.Password), viper.GetInt("cost"))
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, name, role, password, created_at, updated_at) VALUES(?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, u.Email, u.Name, u.Role, string(hp))
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)

		if ok && merr.Number == 1062 {
			return db.ErrDuplicateEmail
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
	conn, _ := db.GetConnection()

	stmt := `UPDATE users SET name = ?, email = ?, role = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, u.Name, u.Email, u.Role, u.ID)
	return err
}

func (u *User) UpdatePassword(pw string) error {
	conn, _ := db.GetConnection()
	pass, err := bcrypt.GenerateFromPassword([]byte(pw), viper.GetInt("cost"))

	stmt := `UPDATE users SET password = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err = conn.Exec(stmt, string(pass), u.ID)
	return err
}

func (u *User) Delete() error {
	conn, _ := db.GetConnection()

	stmt := `DELETE FROM users WHERE id = ?`
	_, err := conn.Exec(stmt, u.ID)
	return err
}

func (u *User) Get() error {
	conn, _ := db.GetConnection()

	snippet := `SELECT id, name, email, role FROM users WHERE`

	var err error

	if u.ID != 0 {
		stmt := snippet + ` id = ?`
		err = conn.QueryRow(stmt, u.ID).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	} else {
		stmt := snippet + ` email = ?`
		err = conn.QueryRow(stmt, u.Email).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	}

	if err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	return err
}

func GetUsers() (Users, error) {
	conn, _ := db.GetConnection()

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
	conn, _ := db.GetConnection()

	var hp []byte

	stmt := `SELECT id, name, email, role, password FROM users WHERE email = ?`

	err := conn.QueryRow(stmt, u.Email).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &hp)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hp, []byte(pw))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return db.ErrInvalidCredentials
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
	conn, _ := db.GetConnection()

	stmt := `UPDATE users SET recovery_hash = ?, updated_at = UTC_TIMESTAMP() WHERE email = ?`
	_, err := conn.Exec(stmt, h, u.Email)
	return err
}

func (u *User) CheckRecover(h string) error {
	conn, _ := db.GetConnection()

	stmt := `SELECT email FROM users WHERE email = ? AND recovery_hash = ?`
	err := conn.QueryRow(stmt, u.Email, h).Scan(&u.Email)
	if err != nil && err == sql.ErrNoRows {
		return db.ErrNotFound
	}

	return err
}

func (u *User) Recover(h string, pw string) error {
	err := u.Get()
	if err != nil {
		if err == sql.ErrNoRows {
			return db.ErrNotFound
		}
		return err
	}

	err = u.CheckRecover(h)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.ErrNotFound
		}
		return err
	}

	u.SetRecover("")
	err = u.UpdatePassword(pw)
	return err
}
