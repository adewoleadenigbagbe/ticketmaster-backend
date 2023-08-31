package utilities

import (
	"fmt"
	"sort"
	"testing"

	"github.com/Wolechacho/ticketmaster-backend/database/entities"
	sequentialguid "github.com/Wolechacho/ticketmaster-backend/helpers"
)

const size = 2000

func TestEntitySortById(t *testing.T) {

	cities := []entities.City{}
	tempcities := make([]entities.City, size)

	for i := 0; i < size; i++ {
		//use city entity
		city := entities.City{
			Id:           sequentialguid.New().String(),
			Name:         fmt.Sprintf("Dave%d", i),
			State:        fmt.Sprintf("State%d", i),
			IsDeprecated: false,
		}

		cities = append(cities, city)
	}

	copy(tempcities[:], cities)

	//sort the original list
	sort.Sort(ByID[entities.City](cities))

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
