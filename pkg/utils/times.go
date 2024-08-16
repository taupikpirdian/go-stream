package utils

import "time"

func EpochToTime(epoch int) time.Time {
	timeE := int64(epoch)
	t := time.Unix(timeE, 0)
	return t
}
