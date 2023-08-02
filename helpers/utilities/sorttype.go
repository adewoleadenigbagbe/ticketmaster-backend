package utilities

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
)

type ByMovieID []entities.Movie

func (s ByMovieID) Len() int {
	return len(s)
}

func (s ByMovieID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByMovieID) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

type ByCityID []entities.City

func (s ByCityID) Len() int {
	return len(s)
}

func (s ByCityID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCityID) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

type ByCinemaID []entities.Cinema

func (s ByCinemaID) Len() int {
	return len(s)
}

func (s ByCinemaID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCinemaID) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

type ByCinemaHallID []entities.CinemaHall

func (s ByCinemaHallID) Len() int {
	return len(s)
}

func (s ByCinemaHallID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCinemaHallID) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

type ByCinemaSeatID []entities.CinemaSeat

func (s ByCinemaSeatID) Len() int {
	return len(s)
}

func (s ByCinemaSeatID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCinemaSeatID) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

// create a time alias
type JsonReleaseDate time.Time

// Implement Marshaler and Unmarshaler interface
func (j *JsonReleaseDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonReleaseDate(t)
	return nil
}

func (j JsonReleaseDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}
