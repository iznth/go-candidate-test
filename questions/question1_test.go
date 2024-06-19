package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstructionQuestionPasses(t *testing.T) { // checks that question is valid
	candidateString := "9103105017084"
	idDetails, err := ParseSAIDNumber(candidateString)
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, idDetails, "idDetails should not be nil")
	assert.Equal(t, Male, idDetails.Gender)
	assert.Equal(t, 30, idDetails.Age)
	assert.True(t, idDetails.SACitizen)
}

func TestBrad2021Passes(t *testing.T) { // checks that question is valid
	candidateString := "9604185215084"
	idDetails, err := ParseSAIDNumber(candidateString)
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, idDetails, "idDetails should not be nil")
	assert.Equal(t, Male, idDetails.Gender)
	assert.Equal(t, 25, idDetails.Age)
	assert.True(t, idDetails.SACitizen)
}

func TestQuestion1PassesValidIDMaleSACitizen(t *testing.T) {
	candidateString := "9101015017087"
	idDetails, err := ParseSAIDNumber(candidateString)
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, idDetails, "idDetails should not be nil")
	assert.Equal(t, Male, idDetails.Gender)
	assert.Equal(t, 30, idDetails.Age)
	assert.True(t, idDetails.SACitizen)
}

func TestQuestion1ShouldFailInvalidIDLength(t *testing.T) {
	candidateString := "910101501708"
	_, err := ParseSAIDNumber(candidateString)
	if err == nil {
		t.Fatal("Should return error")
	}
}

func TestQuestion1ShouldFailCheckSum(t *testing.T) {
	candidateString := "9101015017089"
	_, err := ParseSAIDNumber(candidateString)
	if err == nil {
		t.Fatal("Should return error")
	}
}
