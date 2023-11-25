package main

import (
	"fmt"
	"time"
)

func main() {
	timeStr := "2023-10-11 23:34:45"
	comparedTime, _ := parseInLocalLocation(timeStr, TimeLayoutDefault)
	fmt.Printf("%s晚于当前时间: %v\n", timeStr, isAfterNow(comparedTime))
	originTime := comparedTime
	aWeekLaterTime := weekLater(originTime, 1)
	fmt.Printf("初始时间: %v, 一周后的时间: %v\n", originTime, aWeekLaterTime)
}

const TimeLayoutDefault = "2006-01-02 15:04:05"
const TimeLayoutDate = "20060102"
const TimeLayoutDateWithHyphen = "2006-01-02"


// parseFromString 解析字符串时间
func parseFromString(timeStr, layout string) (timeObj time.Time, err error) {
	timeObj, err = time.Parse(layout, timeStr)
	return
}

// parseInLocalLocation 用本地时区解析字符串时间
func parseInLocalLocation(timeStr, layout string) (timeObj time.Time, err error) {
	timeObj, err = time.ParseInLocation(layout, timeStr, time.Local)
	return
}

// parseInSpainLocation 用西班牙时区解析字符串时间
func parseInSpainLocation(timeStr, layout string) (timeObj time.Time, err error) {
	loc, _ := time.LoadLocation("Europe/Madrid")
	timeObj, err = time.ParseInLocation(layout, timeStr, loc)
	return
}

func isBeforeNow(comparedTime time.Time) bool {
	return comparedTime.Before(time.Now())
}

func isAfterNow(comparedTime time.Time) bool {
	return comparedTime.After(time.Now())
}

func weekLater(originTime time.Time, w int) time.Time {
	return originTime.AddDate(0, 0, w * 7)
}
