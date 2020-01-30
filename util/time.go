package util

import "time"

// TimeToString convert time to string
func TimeToString(t time.Time, format ...string) string {
	_format := "2006-01-02"
	if len(format) == 1 {
		_format = format[0]
	}

	return t.Format(_format)
}

// StringToTime convert string to time
func StringToTime(t string, format ...string) (time.Time, error) {
	_format := "2006-01-02"
	if len(format) == 1 {
		_format = format[0]
	}

	return time.Parse(_format, t)
}
