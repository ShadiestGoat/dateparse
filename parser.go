package dateparse

import (
	"strings"
	"time"
)

//go:generate go run gen/main.go

type Parser interface {
	Parse(str string) (time.Time, error)
}

type DateTimeParser struct {
	LayoutsDate []string
	LayoutsTime []string

	// If true, does not parse unix timestamps
	DisableTimestampParsing bool
	// If true, disables the "today" date format
	DisableToday bool
}

// Parses exclusively time. Input must be a 1 word string. tz must not be nil.
func (d DateTimeParser) ParseTime(inp string, tz *time.Location) (time.Duration, error) {
	t, _, err := parse(strings.ToUpper(inp), d.LayoutsTime, tz)
	if err != nil {
		return 0, err
	}
	return t.Sub(minTime), nil
}

var cleanDateSuffix = map[string]string{
	"1st": "1", "2nd": "2", "2d": "2", "3d": "3", "3rd": "3", "th": "",
}

// Parses exclusively dates. tz must not be nil. inp should be an array of words. Cleans up the dates slightly:
// Removes "of", cleans suffixes like "1st" or "4th".
func (d DateTimeParser) ParseDate(inp []string, tz *time.Location) (time.Time, error) {
	pInp := make([]string, 0, len(inp))
	for _, v := range inp {
		v = strings.ToLower(v)
		if v == "of" {
			continue
		}
		for s, c := range cleanDateSuffix {
			if v2, ok := strings.CutSuffix(v, s); ok {
				v = v2 + c
				break
			}
		}
		pInp = append(pInp, v)
	}
	t, l, err := parse(strings.Join(pInp, " "), d.LayoutsDate, tz)
	if err != nil {
		return time.Time{}, err
	}

	if t.Year() == 0 && !strings.Contains(l, "06") {
		return time.Date(Today(tz).Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz), nil
	}
	return t, nil
}

// Parses a date/time string.
func (d DateTimeParser) Parse(str string) (time.Time, error) {
	spl := strings.Split(str, " ")
	tz := time.UTC

	if len(spl) == 1 {
		if !d.DisableTimestampParsing {
			t, err := ParseUnix(spl[0])
			if err == nil {
				return t, nil
			}
		}
	} else {
		tmpTZ, err := time.LoadLocation(spl[len(spl)-1])
		if err == nil {
			tz = tmpTZ
			spl = spl[:len(spl)-1]
		}
	}

	var fTime time.Duration

	// If an error happens while parsing time, this means that this entire string is bad.
	if isTime(spl[0]) {
		pTime, err := d.ParseTime(spl[0], tz)
		if err != nil {
			return time.Time{}, err
		}
		fTime = pTime
		spl = spl[1:]
	} else if isTime(spl[len(spl)-1]) {
		pTime, err := d.ParseTime(spl[len(spl)-1], tz)
		if err != nil {
			return time.Time{}, err
		}
		fTime = pTime
		spl = spl[:len(spl)-1]
	}

	// Only time is found
	if len(spl) == 0 {
		return Today(tz).Add(fTime), nil
	}

	if len(spl) != 0 && spl[len(spl)-1] == "at" {
		spl = spl[:len(spl)-1]
	}

	// only time is found
	if len(spl) == 0 {
		return Today(tz).Add(fTime), nil
	}

	if len(spl) == 1 && spl[0] == "today" && !d.DisableToday {
		return Today(tz), nil
	}

	parsed, err := d.ParseDate(spl, tz)

	if err != nil {
		return time.Time{}, err
	}

	if fTime != 0 {
		return parsed.Add(fTime), nil
	}

	return parsed, nil
}

// The default parser, which uses the DefaultLayoutsDate & DefaultLayoutsTime.
var DefaultParser = &DateTimeParser{
	LayoutsDate: DefaultLayoutsDate,
	LayoutsTime: DefaultLayoutsTime,
}

// Parse using default parser. A direct alias for `DefaultParser.Parse`.
func Parse(str string) (time.Time, error) {
	return DefaultParser.Parse(str)
}
