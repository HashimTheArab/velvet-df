package utils

import (
	"time"
	_ "unsafe"
)

//go:linkname unitMap time.unitMap
var unitMap map[string]int64

func init() {
	unitMap["d"] = int64(time.Hour * 24) // todo
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

/*
public static function stringToTime(string $string) : ?DateTime {
		if(trim($string) === '') return null;

		$t = new DateTime();

		preg_match_all('/[0-9]+(y|mo|w|d|h|m|s)|[0-9]+/', $string, $found);

		if(count($found[0]) < 1 || count($found[1]) < 1) return null;
		$m = match($found[1][0]){
			'y' => 'year',
			'mo' => 'month',
			'w' => 'week',
			'd' => 'day',
			'h' => 'hour',
			'm' => 'minute',
			's' => 'second',
			default => null
		};
		if(is_null($m)) return null;
		$t->modify('+' . preg_replace('/[^0-9]/', '', $found[0][0]) . " $m");
		return $t;
	}
*/

// DurationToString formats a time duration into a string
// eg: a duration of 14 days will be returned as "14 Days"
func DurationToString(t time.Duration) string {
	return t.String()
}
