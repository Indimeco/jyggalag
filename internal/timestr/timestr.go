package timestr

import (
	"time"
)

func GetCurrentDate() string {
	time := time.Now()
	return time.Format("02")
}

func GetCurrentYear() string {
	time := time.Now()
	return time.Format("2006")
}

func GetCurrentMonth() string {
	time := time.Now()
	return time.Format("01")
}

func GetCanonicalDateString() string {
	return time.Now().Format("2006-01-02")
}
