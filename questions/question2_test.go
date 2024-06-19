package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsGoodFibNumberPasses(t *testing.T) {
	assert.True(t, IsFibNumber(39088169, 0, 1))
}

func Test_IsBadFibNumberFailes(t *testing.T) {
	assert.False(t, IsFibNumber(701408732, 0, 1))
}
