package db

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gosimple/slug"
)

const (
	timeFormat = "2006-01-02 15:04"
)

func toTime(t string) time.Time {
	dt, err := time.Parse(timeFormat, t)
	if err != nil {
		dt = time.Now()
	}
	return dt
}

func getSlug(str string, t string) string {
	var id int
	var err error

	conn, _ := GetConnection()
	stmt := `SELECT id FROM ` + t + ` WHERE slug = ?`

	s := slug.Make(str)
	sl := s

	num := 0
	for err != sql.ErrNoRows {
		if num != 0 {
			sl = s + "-" + strconv.Itoa(num)
		}
		row := conn.QueryRow(stmt, sl)
		err = row.Scan(&id)
		num++
	}
	return sl
}
