package main

import (
	"fmt"
	"time"
)

// LogUTC prints a log message formatted with a UTC ISO8601 timestamp with microsecond precision.
func LogUTC(format string, a ...interface{}) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000000Z")
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", timestamp, msg)
}
