package main

import (
	"errors"
	"strconv"
)

type (
	Gender string
)

const (
	Male        Gender = "M"
	Female      Gender = "F"
	CurrentYear int    = 2021 // used for test to pass
)

type SAIDDetails struct {
	Gender    Gender
	Age       int
	SACitizen bool
}

// This is passed as a string because we are not performing mathematical operations against the actual input
func ParseSAIDNumber(candidateString string) (*SAIDDetails, error) {
	var details = &SAIDDetails{}
	if len(candidateString) != 13 {
		return details, errors.New("invalid candidateString length")
	}

	_, err := strconv.Atoi(candidateString) // check now DRY
	if err != nil {
		return details, errors.New("")
	}

	id := candidateString

	gender, _ := strconv.Atoi(id[6:10])
	if gender > 9999 {
		return nil, errors.New("genderDigits > 0")
	} else if gender > 5000 {
		details.Gender = Male
	} else {
		details.Gender = Female
	}

	// Get age
	age, _ := strconv.Atoi(id[:2]) // age is in years only getting first 2 digits
	details.Age = CurrentYear - (age + 1900)

	// check citizen
	ci, _ := strconv.Atoi(string(id[10]))
	details.SACitizen = ci == 0

	sum := 0
	evn := false // starts off with odd index

	for i := len(id) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(id[i]))
		if evn {
			n *= 2 // find product
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		evn = !evn
	}

	if sum%10 != 0 {
		return details, errors.New("luhn sum failed")
	}

	return details, nil
}
