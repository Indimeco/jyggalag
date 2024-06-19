package timestr

import (
	"time"
)

func CurrentDate() string {
	time := time.Now()
	return time.Format("02")
}

func CurrentYear() string {
	time := time.Now()
	return time.Format("2006")
}

func CurrentMonth() string {
	time := time.Now()
	return time.Format("01")
}

func CanonicalDateString() string {
	return time.Now().Format("2006-01-02")
}
