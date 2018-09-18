package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayContains(t *testing.T) {
	arrInt := []int{4, 5, 6}
	result := ArrayContains(5, arrInt)

	assert := assert.New(t)
	assert.True(result)

	result = ArrayContains(7, arrInt)
	assert.False(result)

	arrStr := []string{"hi", ":id", "432"}
	result = ArrayContains(":id", arrStr)
	assert.True(result)

	result = ArrayContains("id", arrStr)
	assert.False(result)
}
