package enums

import (
	"testing"
)

func TestAdd(t *testing.T) {
	expected := Genre(23)
	defaultGenre := FlaggedEnum[Genre].New(FlaggedEnum[Genre](None))

	flaggedGenre := defaultGenre.Add(Action, Adventure, Animation, Crime)

	if flaggedGenre != expected {
		t.Errorf("Expected Value %d is  got %d", expected, int64(flaggedGenre))
	}

}

func TestRemove(t *testing.T) {
	expected := Genre(67)

	d := FlaggedEnum[Genre](0)
	flaggedGenre := d.Add(Action, Adventure, Crime, Drama)
	genre := d.Remove(flaggedGenre, Crime)

	if genre != expected {
		t.Errorf("Expected Value : %d but got : %d", expected, int64(flaggedGenre))
	}
}

func TestShouldNotRemoveIfEnumValueisNotPresent(t *testing.T) {
	expected := Genre(83)

	d := FlaggedEnum[Genre](0)
	flaggedGenre := d.Add(Action, Adventure, Crime, Drama)
	genre := d.Remove(flaggedGenre, Family)

	if genre != expected {
		t.Errorf("Should expected value %d is to be the same with %d", expected, int64(genre))
	}

}

func TestHas(t *testing.T) {

	tests := map[string]struct {
		flags    Genre
		flag     Genre
		expected bool
	}{
		"True": {
			flags:    FlaggedEnum[Genre](0).Add(Action, Adventure, Crime, Drama),
			flag:     Crime,
			expected: true,
		},
		"False": {
			flags:    FlaggedEnum[Genre](0).Add(Action, Adventure, Drama),
			flag:     Crime,
			expected: false,
		},
	}

	d := FlaggedEnum[Genre](0)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			hasFlag := d.Has(test.flags, test.flag)

			if hasFlag != test.expected {
				t.Errorf("expected %v, got %v", test.expected, hasFlag)
			}
		})
	}

	flaggedGenre := d.Add(Action, Adventure, Crime, Drama)
	hasFlag := d.Has(flaggedGenre, Crime)

	if !hasFlag {
		t.Error("should have a flag")
	}
}
