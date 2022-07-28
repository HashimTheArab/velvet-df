package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "unsafe"
)

const CoolAssArrow = "»"

var (
	valueRegex = regexp.MustCompile(`^\d+`)
	unitRegex  = regexp.MustCompile(`\w$`)
)

// ParseDuration attempts to parse a duration string. If successful, it returns a time.Duration, otherwise returning
// an error in its second value.
func ParseDuration(str string) (time.Duration, error) {
	matches := valueRegex.FindAllString(str, 1)
	if len(matches) <= 0 {
		return 0, fmt.Errorf("parse duration: invalid value")
	}
	value, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, fmt.Errorf("parse duration: value must be greater than zero")
	}
	unit := unitRegex.FindAllString(str, 1)[0]
	if duration, ok := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
	}[unit]; ok {
		return time.Duration(value) * duration, nil
	}
	return 0, fmt.Errorf("parse duration: invalid unit")
}

// PrettyParams converts the given parameters to a readable string.
func PrettyParams(params map[string]any) string {
	if len(params) == 0 {
		// Don't waste our time if there are no parameters.
		return "[]"
	}
	// Hacky but simple way to create a readable string.
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimPrefix(fmt.Sprint(params), "map"), " ", ", "), ":", "=")
}
