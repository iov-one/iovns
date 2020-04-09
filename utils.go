package iovnsd

import (
	"fmt"
	"time"
)

// SecondsToTime converts unix seconds to time
func SecondsToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}

func GetAccountKey(domain, name string) string {
	return fmt.Sprintf("%s*%s", domain, name)
}
