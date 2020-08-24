package extensions

import (
	"time"

	"github.com/jaceklubzinski/pd-tools-bot/pkg/base"
)

func BeginningOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func EndOfMonth(date time.Time) time.Time {
	firstDay := BeginningOfMonth(date)
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return lastDay
}

func StringToDate(value string) time.Time {
	layoutISO := "2006-01-02 15:04:05"
	converted, err := time.Parse(layoutISO, value)
	base.CheckErr(err)
	return converted
}

//AddDurationToDate add duration to start date
func AddDurationToDate(start string, timer string) string {
	layoutISO := "2006-01-02 15:04:05"
	startDate := StringToDate(start)
	timerDuration, err := time.ParseDuration(timer)
	base.CheckErr(err)
	return startDate.Add(timerDuration).Format(layoutISO)
}
