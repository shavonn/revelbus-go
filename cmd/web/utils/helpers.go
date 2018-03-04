package utils

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func ToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func RandomString(strlen int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(result)
}

func NewNullStr(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt(i int) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(i),
		Valid: true,
	}
}
