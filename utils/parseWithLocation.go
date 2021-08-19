package utils

import (
	"time"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

func ParseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, err
	} else {
		//转成带时区的时间
		//lt, _ := time.ParseInLocation(TIME_LAYOUT, timeStr, l)
		//直接转成相对时间
		//fmt.Println(time.Now().In(l).Format(TIME_LAYOUT))
		lt := time.Now().In(l).Format(TIME_LAYOUT)
		ltTime, _ := time.Parse(TIME_LAYOUT, lt)

		return ltTime, nil
	}
}
