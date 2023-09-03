package utilities

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	MIN_DATE     = "0001-01-01T00:00:00Z"
	MAX_DATE     = "9999-12-31T23:59:59Z"
	DEFAULT_UUID = "00000000-0000-0000-0000-000000000000"
)

type Datetime time.Time

func (d Datetime) Max(datetimes ...time.Time) time.Time {
	maxDate, _ := time.Parse(time.RFC3339, MIN_DATE)
	for i := 0; i < len(datetimes); i++ {
		if datetimes[i].After(maxDate) {
			maxDate = datetimes[i]
		}
	}

	return maxDate
}

func (d Datetime) Min(datetimes ...time.Time) time.Time {
	minDate, _ := time.Parse(time.RFC3339, MAX_DATE)
	for i := 0; i < len(datetimes); i++ {
		if datetimes[i].Before(minDate) {
			minDate = datetimes[i]
		}
	}

	return minDate
}

func (d Datetime) BeforeOrEqualTo(u time.Time) bool {

	dNano := time.Time(d).UnixNano()
	uNano := u.UnixNano()

	return dNano <= uNano
}

func (d Datetime) AfterOrEqualTo(u time.Time) bool {

	dNano := time.Time(d).UnixNano()
	uNano := u.UnixNano()

	return dNano >= uNano
}

func (d *Datetime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*d = Datetime(t)
	return nil
}

func (d Datetime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}
