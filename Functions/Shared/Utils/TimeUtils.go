package Utils

import (
	"fmt"
	"time"
)

func GetDateInUtc(t time.Time) (time.Time, time.Time) {
	beginningOfTheDay := GetBeginningOfTheDayInUTC(t.UTC())
	return beginningOfTheDay, beginningOfTheDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func GetTodayInUtc() (time.Time, time.Time) {
	beginningOfTheDay := GetBeginningOfTheDayInUTC(time.Now().UTC())
	return beginningOfTheDay, beginningOfTheDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func GetADayForLocationByDayMonthYear(day int, month time.Month, year int, location string) (time.Time, time.Time) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
		// return time.Now()
	}
	beginningOfTheDay := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return beginningOfTheDay, beginningOfTheDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func GetADayForLocation(t time.Time, location string) (time.Time, time.Time) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
		// return time.Now()
	}
	year, month, day := t.In(loc).Date()
	beginningOfTheDay := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return beginningOfTheDay, beginningOfTheDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func GetTodayForLocation(location string) (time.Time, time.Time) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
		// return time.Now()
	}
	year, month, day := time.Now().In(loc).Date()
	beginningOfTheDay := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return beginningOfTheDay, beginningOfTheDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
}

func GetBeginningOfTheDayInUTC(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func GetCurrentLocalTime(location string) time.Time {
	loc, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
		// return time.Now()
	}
	return time.Now().In(loc)
}

func UTCtoLocalTime(location string, t time.Time) time.Time {
	loc, err := time.LoadLocation(location)
	if err != nil {
		fmt.Printf("errs %s\n", err)
		panic(err)
	}
	localTime := t.In(loc)
	return localTime
}

func GetUnixEpochTimeInMillis() int64 {
	// Get Unix Epoch time and subtract the number of seconds of the timeFrame which will be in hours.
	now := time.Now()
	secs := now.Unix()
	millis := secs * 1000
	return millis
}
