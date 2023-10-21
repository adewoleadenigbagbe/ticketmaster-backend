package utilities

import (
	"testing"
	"time"
)

func TestMaxDateTime(t *testing.T) {
	now := time.Now()
	expectedMaxDate := now.AddDate(0, 0, 365)
	dates := []time.Time{
		now,
		now.AddDate(0, 0, -30),
		now.AddDate(0, 0, 60),
		expectedMaxDate,
		now.AddDate(0, 0, 300),
		now.AddDate(0, 0, 150),
	}

	maxDate := Datetime(now).Max(dates...)

	if expectedMaxDate != maxDate {
		t.Errorf("Expected Max Date %s but got %s", expectedMaxDate.Format(time.RFC3339), maxDate.Format(time.RFC3339))
	}
}

func TestMinDateTime(t *testing.T) {
	now := time.Now()
	expectedMinDate := now.AddDate(0, 0, -40)
	dates := []time.Time{
		now,
		now.AddDate(0, 0, 60),
		now.AddDate(0, 0, -30),
		now.AddDate(0, 0, -20),
		now.AddDate(0, 0, -10),
		expectedMinDate,
		now.AddDate(0, 0, 5),
		now.AddDate(0, 0, 10),
	}

	minDate := Datetime(now).Min(dates...)

	if expectedMinDate != minDate {
		t.Errorf("Expected Min Date %s but got %s", expectedMinDate.Format(time.RFC3339), minDate.Format(time.RFC3339))
	}
}

func TestBeforeDateTime(t *testing.T) {
	now := time.Now()
	datetime := Datetime(now.AddDate(0, 0, -30))

	before := datetime.BeforeOrEqualTo(now)

	if !before {
		t.Errorf("Expected Date %s should be before %s", time.Time(datetime).Format(time.RFC3339), now.Format(time.RFC3339))
	}
}

func TestAfterDateTime(t *testing.T) {
	now := time.Now()

	datetime := Datetime(now.AddDate(0, 0, 30))

	after := datetime.AfterOrEqualTo(now)

	if !after {
		t.Errorf("Expected Date %s should be before %s", time.Time(datetime).Format(time.RFC3339), now.Format(time.RFC3339))
	}
}

func TestBeforeOrEqualTo(t *testing.T) {
	now := time.Now()
	datetime := Datetime(now)

	beforeOrEqualTo := datetime.BeforeOrEqualTo(now)

	if !beforeOrEqualTo {
		t.Errorf("Expected Date %s should be the same with %s", time.Time(datetime).Format(time.RFC3339), now.Format(time.RFC3339))
	}
}

func TestAfterOrEqualTo(t *testing.T) {
	now := time.Now()
	datetime := Datetime(now)

	afterOrEqualTo := datetime.AfterOrEqualTo(now)

	if !afterOrEqualTo {
		t.Errorf("Expected Date %s should be same with %s", time.Time(datetime).Format(time.RFC3339), now.Format(time.RFC3339))
	}
}

func TestUnMarshalJSON(t *testing.T) {
	parsedDateTime, _ := time.Parse(time.RFC3339, "2013-12-31T11:07:59Z")
	date := Datetime(parsedDateTime)
	s := "2013-12-31"

	err := date.UnmarshalJSON([]byte(s))

	if err != nil {
		t.Error("Could not parse the date string")
	}

	newTime := time.Time(date)
	if newTime.Hour() != 0 && newTime.Minute() != 0 && newTime.Second() != 0 {
		t.Errorf("Expected Hours: 00 and Minutes: 00 and Seconds: 00 but got Hours: %d and Minutes: %d and Seconds: %d", newTime.Hour(), newTime.Minute(), newTime.Second())
	}

}

func TestMarshalJSON(t *testing.T) {
	parsedDateTime, _ := time.Parse(time.RFC3339, "2013-12-31T00:00:00Z")
	date := Datetime(parsedDateTime)

	b, _ := date.MarshalJSON()
	jsonDate := string(b)

	expectedstring := "\"2013-12-31T00:00:00Z\""
	if jsonDate != expectedstring {
		t.Errorf("could not marshal date")
	}
}
