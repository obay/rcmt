package helpers

import (
	"log"
	"strconv"
	"time"
)

func GetEpochTime(epochString string) (t time.Time) {
	var epoch int64
	epoch, err := strconv.ParseInt(epochString, 10, 64)
	if err != nil {
		log.Fatal("Error getting Epoch time")
	}
	t = time.Unix(epoch, 0)
	return
}
