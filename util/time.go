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

// Diff difference between two time
// oldtime = time.Date(2015, 5, 1, 0, 0, 0, 0, time.UTC)
// newtime = time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC)
// fmt.Println(Diff(oldtime, newtime))  // Expected: 1 1 1 1 1 1
func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}

	if min < 0 {
		min += 60
		hour--
	}

	if hour < 0 {
		hour += 24
		day--
	}

	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}

	if month < 0 {
		month += 12
		year--
	}

	return
}
