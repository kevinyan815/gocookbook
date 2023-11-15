package main

import (
	"fmt"
	"time"
)

func main() {
	timeStr := "2023-10-11 23:34:45"
	comparedTime, _ := parseFromString(timeStr, TimeLayoutDefault)
	fmt.Printf("%s晚于当前时间: %v", timeStr, isAfterNow(comparedTime))
}

const TimeLayoutDefault = "2006-01-02 15:04:05"
const TimeLayoutDate = "20060102"
const TimeLayoutDateWithHyphen = "2006-01-02"


func parseFromString(timeStr, layout string) (timeObj time.Time, err error) {
	timeObj, err = time.Parse(layout, timeStr)
	return
}

func isBeforeNow(comparedTime time.Time) bool {
	return comparedTime.Before(time.Now())
}

func isAfterNow(comparedTime time.Time) bool {
	return comparedTime.After(time.Now())
}
