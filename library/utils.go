package library

import (
	"fmt"
	"time"
)

func Log(format string, args ...interface{}) {
	now := time.Now()
	dt := fmt.Sprintf("%d:%d:%d+%d", now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000000)
	other := fmt.Sprintf(format, args...)
	fmt.Printf("%s: %s\n", dt, other)
}
