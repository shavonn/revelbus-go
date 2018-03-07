package models

import (
	"database/sql"
	"revelbus/internal/platform/domain"
	"revelbus/internal/platform/forms"
	"revelbus/pkg/database"
)

type FAQ struct {
	ID       int
	Question sql.NullString
	Answer   sql.NullString
	Category sql.NullString
	Order    sql.NullInt64
	Active   bool
}

type FAQs []*FAQ

type GroupedFAQs map[string][]*FAQ

type FAQForm struct {
	ID       string
	Question string
	Answer   string
	Category string
	Order    string
	Active   bool

	Errors map[string]string
}

func (f *FAQForm) Valid() bool {
	v := forms.NewValidator()

	v.Required("Question", f.Question)
	v.Required("Answer", f.Answer)
	v.Required("Category", f.Category)
	v.ValidInt("Order", f.Order)

	f.Errors = v.Errors
	return len(f.Errors) == 0
}

func (f *FAQ) Create() error {
	conn, _ := database.GetConnection()

	stmt := `INSERT INTO faqs (question, answer, category, sort_order, active, created_at, updated_at) VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	result, err := conn.Exec(stmt, f.Question, f.Answer, f.Category, f.Order, f.Active)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = int(id)
	return nil
}

func (f *FAQ) Fetch() error {
	conn, _ := database.GetConnection()

	stmt := `SELECT question, answer, category, sort_order, active FROM faqs WHERE id = ?`
	err := conn.QueryRow(stmt, f.ID).Scan(&f.Question, &f.Answer, &f.Category, &f.Order, &f.Active)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (f *FAQ) Update() error {
	conn, _ := database.GetConnection()

	stmt := `UPDATE faqs SET question = ?, answer= ?, category = ?, sort_order = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, f.Question, f.Answer, f.Category, f.Order, f.Active, f.ID)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	return err
}

func (f *FAQ) Delete() error {
	conn, _ := database.GetConnection()

	stmt := `DELETE FROM faqs WHERE id = ?`
	_, err := conn.Exec(stmt, f.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func FetchFAQs() (*FAQs, error) {
	conn, _ := database.GetConnection()

	stmt := `SELECT id, question, answer, category, sort_order, active FROM faqs ORDER BY sort_order`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	faqs := FAQs{}
	for rows.Next() {
		f := &FAQ{}
		err := rows.Scan(&f.ID, &f.Question, &f.Answer, &f.Category, &f.Order, &f.Active)
		if err != nil {
			return nil, err
		}
		faqs = append(faqs, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faqs, nil
}

func FindActiveFAQs() (*GroupedFAQs, error) {
	conn, _ := database.GetConnection()

	faqs := make(GroupedFAQs)

	stmt := `SELECT id, question, answer, category, sort_order, active FROM faqs WHERE active = 1 ORDER BY sort_order`
	rows, err := conn.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := &FAQ{}
		err := rows.Scan(&f.ID, &f.Question, &f.Answer, &f.Category, &f.Order, &f.Active)
		if err != nil {
			return nil, err
		}
		faqs[f.Category.String] = append(faqs[f.Category.String], f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faqs, nil
}
