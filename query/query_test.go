package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains_ReturnsTrue_WhenIntInArray(t *testing.T) {
	assert.True(t, ArrayContains(5, []int{4, 5, 6}))
}

func TestArrayContains_ReturnsFalse_WhenIntNotInArray(t *testing.T) {
	assert.False(t, ArrayContains(7, []int{4, 5, 6}))
}

func TestArrayContains_ReturnsTrue_WhenStringInArray(t *testing.T) {
	assert.True(t, ArrayContains(":id", []string{"hi", ":id", "432"}))
}

func TestArrayContains_ReturnsFalse_WhenStringNotInArray(t *testing.T) {
	assert.False(t, ArrayContains("id", []string{"hi", ":id", "432"}))
}
