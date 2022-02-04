package utils

import (
	"time"
	_ "unsafe"
)

const CoolAssArrow = "»"

//go:linkname unitMap time.unitMap
var unitMap map[string]int64

func init() {
	unitMap["d"] = int64(time.Hour * 24)
}

// DurationFromString returns a time duration from a string such as "14d" or "10h".
// If an invalid string is passed, this function will return -1.
func DurationFromString(t string) time.Duration {
	if t == "" {
		return -1
	}
	parsed, err := time.ParseDuration(t)
	if err != nil {
		return -1
	}
	return parsed
}

// DurationToString formats a time duration into a string
// eg: a duration of 14 days will be returned as "14 Days"
func DurationToString(t time.Duration) string {
	return t.String()
}
