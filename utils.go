package iovnsd

import "time"

// SecondsToTime converts unix seconds to time
func SecondsToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}
