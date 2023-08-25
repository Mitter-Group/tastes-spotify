package util

import "time"

func GetTimeDiffInMillis(t time.Time) int64 {
	start := t.UnixNano() / int64(time.Millisecond)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	diff := end - start
	return diff
}
