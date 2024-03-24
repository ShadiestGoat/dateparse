package dateparse_test

import (
	"testing"
	"time"

	"github.com/shadiestgoat/dateparse"
)

func testOutput(t *testing.T, inp string, err error, real, expected time.Time) {
	if err != nil || real.UnixNano() != expected.UnixNano() {
		t.Logf("Failed parsing '%v'! Wanted: '%v', Real: '%v', Err: '%v'\n", inp, expected, real, err)
		t.Fail()
	}
}

func mkCommon(list []string, d time.Time, m map[string]time.Time) {
	for _, v := range list {
		m[v] = d
	}
}

func TestParse(t *testing.T) {
	// tz2, err := time.LoadLocation("PST")
	// if err != nil {
	// 	t.Log("Couldn't load 2nd timezone:", err)
	// 	t.FailNow()
	// }

	timeStrings := map[string]time.Duration{
		"5pm":       17 * time.Hour,
		"17:00":     17 * time.Hour,
		"05:0:3":    5*time.Hour + 3*time.Second,
		"05:0:3am":  5*time.Hour + 3*time.Second,
		"05:20:3pm": 17*time.Hour + 20*time.Minute + 3*time.Second,
	}

	today := dateparse.Today(time.UTC)

	dateStrings := map[string]time.Time{
		"today":        today,
		"01/06":        time.Date(time.Now().Year(), time.June, 1, 0, 0, 0, 0, time.UTC),
		"1/6":          time.Date(time.Now().Year(), time.June, 1, 0, 0, 0, 0, time.UTC),
		"1st of June":  time.Date(time.Now().Year(), time.June, 1, 0, 0, 0, 0, time.UTC),
		"2nd of June":  time.Date(time.Now().Year(), time.June, 2, 0, 0, 0, 0, time.UTC),
		"3d of June":   time.Date(time.Now().Year(), time.June, 3, 0, 0, 0, 0, time.UTC),
		"4th of June":  time.Date(time.Now().Year(), time.June, 4, 0, 0, 0, 0, time.UTC),
		"21st of June": time.Date(time.Now().Year(), time.June, 21, 0, 0, 0, 0, time.UTC),
		"15 Dec 15":    time.Date(2015, time.December, 15, 0, 0, 0, 0, time.UTC),
	}

	mkCommon([]string{
		"1/6/05",
		"01/6/2005",
		"1/06/05",
		"1-06-05",
		"6/2005",
		"2005-06-01",
		"1st of June 2005",
		"01 Jun 05",
		"Jun 05",
	}, time.Date(2005, time.June, 1, 0, 0, 0, 0, time.UTC), dateStrings)

	combos := map[string]time.Time{
		"3d Aug 13 3pm": time.Date(2013, time.August, 3, 15, 0, 0, 0, time.UTC),
		"3pm 3d Aug 15": time.Date(2015, time.August, 3, 15, 0, 0, 0, time.UTC),
	}

	mkCommon([]string{
		"3pm 2nd of July",
		"2nd of July 3pm",
		"2nd of July 15:00",
		"2nd July at 15:00",
	}, time.Date(time.Now().Year(), time.July, 2, 15, 0, 0, 0, time.UTC), combos)

	t.Run("timeAlone", func(t *testing.T) {
		for inp, exp := range timeStrings {
			expected := today.Add(exp)
			real, err := dateparse.Parse(inp)
			testOutput(t, inp, err, real, expected)
		}
	})

	t.Run("dateAlone", func(t *testing.T) {
		for inp, expected := range dateStrings {
			real, err := dateparse.Parse(inp)
			testOutput(t, inp, err, real, expected)
		}
	})

	t.Run("combos", func(t *testing.T) {
		for inp, expected := range combos {
			real, err := dateparse.Parse(inp)
			testOutput(t, inp, err, real, expected)
		}
	})
}
