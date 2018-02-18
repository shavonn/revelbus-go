package db

import (
	"database/sql"
)

type FAQ struct {
	ID       int
	Question string
	Answer   string
	Category string
	Active   bool
	Order    int
}

type FAQs []*FAQ
type GroupedFAQ map[string][]*FAQ

func (f *FAQ) Create() error {
	conn, _ := GetConnection()

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

func (f *FAQ) Update() error {
	conn, _ := GetConnection()

	stmt := `UPDATE faqs SET question = ?, answer= ?, category = ?, sort_order = ?, active = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?`
	_, err := conn.Exec(stmt, f.Question, f.Answer, f.Category, f.Order, f.Active, f.ID)
	return err
}

func (s *FAQ) Delete() error {
	conn, _ := GetConnection()

	stmt := `DELETE FROM faqs WHERE id = ?`
	_, err := conn.Exec(stmt, s.ID)
	return err
}

func (f *FAQ) Get() error {
	conn, _ := GetConnection()

	stmt := `SELECT question, answer, category, sort_order, active FROM faqs WHERE id = ?`
	err := conn.QueryRow(stmt, f.ID).Scan(&f.Question, &f.Answer, &f.Category, &f.Order, &f.Active)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}

func GetFAQs() (*FAQs, error) {
	conn, _ := GetConnection()

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

func GetActiveFAQs() (*GroupedFAQ, error) {
	conn, _ := GetConnection()

	faqs := make(GroupedFAQ)

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
		faqs[f.Category] = append(faqs[f.Category], f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faqs, nil
}
