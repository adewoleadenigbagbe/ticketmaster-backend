package entities

import (
	"fmt"
	"sort"
	"testing"

	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
)

const size = 2000

func TestEntitySortById(t *testing.T) {

	cities := []City{}
	tempcities := make([]City, size)

	for i := 0; i < size; i++ {
		//use city entity
		city := City{
			Id:           sequentialguid.New().String(),
			Name:         fmt.Sprintf("Dave%d", i),
			State:        fmt.Sprintf("State%d", i),
			IsDeprecated: false,
		}

		cities = append(cities, city)
	}

	copy(tempcities[:], cities)

	//sort the original list
	sort.Sort(ByID[City](cities))

	sorted := false
	for i := 0; i < size; i++ {
		if tempcities[i] != cities[i] {
			sorted = true
		}

		if sorted {
			break
		}
	}

	if !sorted {
		t.Error("Temporary and original list already ordered by Entity Id")
	}
}
