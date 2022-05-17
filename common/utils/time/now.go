package timeUtils

import (
	"time"
)

func NowStrDatetime() string {
	return Now().Format("2006-01-02 15:04:05")
}

func Now() time.Time {
	return time.Now()
}

//Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return Now().UnixNano() / 1e6
}
