package common

import (
	"fmt"
	"math"
	"time"
)

// GetCurrentTime get the current time and check if it's +7 timezone
func GetCurrentTime() (err error, now time.Time) {
	now = time.Now()
	_, offset := now.Zone()
	if offset != 25200 { // 7 hours * 60 minutes * 60 seconds
		err = fmt.Errorf("Timezone is not +7")
	}
	return err, now
}

// GetStartEndOfDay returns 0h0m0s of today and 0h0m0s of tomorrow
func GetStartEndOfDay(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()
	year, month, day := tm.Date()
	startTime = time.Date(year, month, day, 0, 0, 0, 0, loc)
	endTime = time.Date(year, month, day+1, 0, 0, 0, 0, loc)
	return startTime, endTime
}

// GetStartEndOfWeek returns Monday 0h0m0s of current week and Monday 0h0m0s of next week
func GetStartEndOfWeek(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()

	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	year, month, day := tm.Date()
	startTime = time.Date(year, month, day, 0, 0, 0, 0, loc).Add(-1 * (weekday - 1) * 24 * time.Hour)

	endTime = startTime.AddDate(0, 0, 7)
	return startTime, endTime
}

// GetStartEndOfMonth returns 1st day 0h0m0s of current month and 1st day 0h0m0s of next month
func GetStartEndOfMonth(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()
	year, month, _ := tm.Date()
	startTime = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endTime = time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
	return startTime, endTime
}

func CheckValidHour(from int, to int) (error, bool) {
	err, now := GetCurrentTime()
	if err != nil {
		return err, false
	}

	hour, _, _ := now.Clock()
	return nil, to >= hour || hour >= from
}

func ConvertUnixToTime(un float64) (error, time.Time) {
	sec, dec := math.Modf(un)
	return nil, time.Unix(int64(sec), int64(dec*(1e9)))
}

// GetCurrentTime get the current time and check if it's +7 timezone
func GetCurrentUnixTime() (err error, unixT int64) {
	now := time.Now()
	_, offset := now.Zone()
	if offset != 25200 { // 7 hours * 60 minutes * 60 seconds
		err = fmt.Errorf("Timezone is not +7")
	}
	return err, now.Unix()
}
