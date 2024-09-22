package sequentialguid

import (
	"testing"
)

const UUID_SIZE int = 36

func TestNew(t *testing.T) {

	got := New().String()

	if len(got) != UUID_SIZE {
		t.Errorf("got %d, wanted %d", len(got), UUID_SIZE)
	}
}
