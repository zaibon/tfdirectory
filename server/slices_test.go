package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSubSet(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "3"}

	assert.True(t, isSubset(a, b))
	assert.False(t, isSubset(b, a))
}
