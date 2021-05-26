package helpers

import (
	"log"
	"strconv"
)

func StringToInt(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Error getting integer from string")
	}
	return
}

func StringToInt64(s string) (n int64) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal("Error integer 64 from string")
	}
	return
}
