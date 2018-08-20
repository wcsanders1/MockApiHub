package str

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	port := GetPort(5000)
	assert.Equal(t, ":5000", port)
}