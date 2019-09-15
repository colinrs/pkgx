package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

//  const ...
const (
	// FIVEMIN1 ...
	FIVEMIN1 = "200601021504"
	// FIVEMIN2 ...
	FIVEMIN2 = "2006-01-02 15:04"
)

var (

	// FiveMinParten 2006-01-02 15:04
	FiveMinParten = regexp.MustCompile("\\d{4}-\\d{2}-\\d \\d{2}:\\d{2}")
	// TimeStr1 200601021504
	TimeStr1 = regexp.MustCompile("\\d{12}")
	// UninxTimeStr 1568304212
	UninxTimeStr = regexp.MustCompile("\\d{10}")
)

// ParseTimeStr2Time ...
func ParseTimeStr2Time(timeStr string) (t time.Time, err error) {

	switch {
	case FiveMinParten.MatchString(timeStr):
		//fmt.Print("2006-01-02 15:04\n")
		t, _ = time.Parse(FIVEMIN2, timeStr)
	case TimeStr1.MatchString(timeStr):
		//fmt.Print("200601021504\n")
		t, _ = time.Parse(FIVEMIN1, timeStr)

	case UninxTimeStr.MatchString(timeStr):
		//fmt.Print("1568304212\n")
		timestamp, _ := strconv.ParseInt(timeStr, 10, 64)
		t = time.Unix(timestamp, 0)
	default:
		//fmt.Print("xxx\n")
		t = time.Now()
		err = fmt.Errorf("Not Support this Format to time.Time")
	}
	return
}
