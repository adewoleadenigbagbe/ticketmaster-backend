package sequentialguid

import (
	"testing"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
)

const UUID_SIZE int = 36

func TestNew(t *testing.T) {

	got := sequentialguid.New().String()

	if len(got) != UUID_SIZE {
		t.Errorf("got %d, wanted %d", len(got), UUID_SIZE)
	}
}
