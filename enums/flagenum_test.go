package enums

import (
	"testing"
)

func TestAdd(t *testing.T) {
	expected := 23
	defaultGenre := FlaggedEnum[Genre].New(FlaggedEnum[Genre](None))

	flaggedGenre := defaultGenre.Add(Action, Adventure, Animation, Crime)

	if flaggedGenre != Genre(expected) {
		t.Errorf("Expected Value %d is  got %d", expected, int64(flaggedGenre))
	}

}

func TestRemove(t *testing.T) {
	expected := 67

	d := FlaggedEnum[Genre](0)
	flaggedGenre := d.Add(Action, Adventure, Crime, Drama)
	genre := d.Remove(flaggedGenre, Crime)

	if genre != Genre(expected) {
		t.Errorf("Expected Value : %d but got : %d", expected, int64(flaggedGenre))
	}
}

func TestShouldNotRemoveIfEnumValueisNotPresent(t *testing.T) {
	expected := 83

	d := FlaggedEnum[Genre](0)
	flaggedGenre := d.Add(Action, Adventure, Crime, Drama)
	genre := d.Remove(flaggedGenre, Family)

	if genre != Genre(expected) {
		t.Errorf("Should expected value %d is to be the same with %d", expected, int64(genre))
	}

}
