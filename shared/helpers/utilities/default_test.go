package utilities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	//act
	type Person struct {
		FirstName string    `default:"Jones"`
		Weight    float64   `default:"98.75"`
		IsAdmin   bool      `default:"true"`
		Age       int       `default:"24"`
		DOB       time.Time `default:"2023-08-08T02:00:00Z"`
	}

	firstName := "Jones"
	age := 24
	weight := 98.75
	isAdmin := true
	datetime, _ := time.Parse(time.RFC3339, "2023-08-08T02:00:00Z")

	var person Person

	//Arrange
	SetDefaults(&person)

	//Assert
	assert.Equal(t, firstName, person.FirstName, "they should be equal")
	assert.Equal(t, weight, person.Weight, "they should be equal")
	assert.Equal(t, isAdmin, person.IsAdmin, "they should be equal")
	assert.Equal(t, age, person.Age, "they should be equal")
	assert.Equal(t, datetime, person.DOB, "they should be equal")
}

func TestShouldNotSetDefaults(t *testing.T) {
	//act
	type Person struct {
		FirstName string    `default:"Jones"`
		Weight    float64   `default:"98.75"`
		IsAdmin   bool      `default:"false"`
		Age       int       `default:"24"`
		DOB       time.Time `default:"2023-08-08T02:00:00Z"`
	}

	firstName := "Dave"
	age := 76
	weight := 109.97
	isAdmin := true
	datetime, _ := time.Parse(time.RFC3339, "2023-09-09T02:00:00Z")

	person := Person{
		FirstName: firstName,
		Age:       age,
		Weight:    weight,
		IsAdmin:   isAdmin,
		DOB:       datetime,
	}

	//Arrange
	SetDefaults(&person)

	//Assert
	assert.Equal(t, firstName, person.FirstName, "they should be equal")
	assert.Equal(t, weight, person.Weight, "they should be equal")
	assert.Equal(t, isAdmin, person.IsAdmin, "they should be equal")
	assert.Equal(t, age, person.Age, "they should be equal")
	assert.Equal(t, datetime, person.DOB, "they should be equal")
}

func TestShouldNotDefaultBoolValueToFalse(t *testing.T) {
	//act
	type Person struct {
		FirstName string    `default:"Jones"`
		Weight    float64   `default:"98.75"`
		IsAdmin   bool      `default:"false"`
		Age       int       `default:"24"`
		DOB       time.Time `default:"2023-08-08T02:00:00Z"`
	}

	firstName := "Dave"
	age := 76
	weight := 109.97
	isAdmin := false
	datetime, _ := time.Parse(time.RFC3339, "2023-09-09T02:00:00Z")

	person := Person{
		FirstName: firstName,
		Age:       age,
		Weight:    weight,
		IsAdmin:   isAdmin,
		DOB:       datetime,
	}

	//Arrange
	err := SetDefaults(&person)
	errorString := "cannot default bool value to be false"

	//Assert
	assert.ErrorContains(t, err, errorString)
}

func TestIfDefaultNotStruct(t *testing.T) {
	//act
	word := "Jones"

	//Arrange
	err := SetDefaults(&word)
	errorString := "default allowed only allowed on struct"

	//Assert
	assert.ErrorContains(t, err, errorString)
}
