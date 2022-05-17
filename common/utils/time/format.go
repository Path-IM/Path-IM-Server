package timeUtils

import "time"

const (
	TimeFormat   = "2006-01-02 15:04:05"
	TimeMSFormat = "2006-01-02 15:04:05.000"
	DateFormat   = "2006-01-02"
	TimeTZFormat = `2006-01-02T15:04:05Z`
)

func StrDatetime(t time.Time) string {
	return t.Format(TimeFormat)
}

func StrDatetimeMs(t time.Time) string {
	return t.Format(TimeMSFormat)
}

func StrDate(t time.Time) string {
	return t.Format(DateFormat)
}

func ParseTime(str string) time.Time {
	t, _ := time.Parse(TimeFormat, str)
	return t
}
