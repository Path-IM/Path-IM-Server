package randUtils

import "time"

func RandTime(min string, max string) string {
	minTime, _ := time.Parse("2006-01-02 15:04:05", min)
	maxTime, _ := time.Parse("2006-01-02 15:04:05", max)
	randTime := minTime.Add(time.Duration(RandInt64(int64(minTime.Unix()), int64(maxTime.Unix()))) * time.Second)
	return randTime.Format("2006-01-02 15:04:05")
}

func RandDate(min string, max string) string {
	minTime, _ := time.Parse("2006-01-02", min)
	maxTime, _ := time.Parse("2006-01-02", max)
	randTime := minTime.Add(time.Duration(RandInt64(int64(minTime.Unix()), int64(maxTime.Unix()))) * time.Second)
	return randTime.Format("2006-01-02")
}
