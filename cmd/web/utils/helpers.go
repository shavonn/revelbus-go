package utils

import (
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
