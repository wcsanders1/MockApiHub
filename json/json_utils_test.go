package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidJSON(t *testing.T) {
	goodJSON := []byte(`{
							"JSON": "good",
							"test": "good"
					   }`)
	badJSON := []byte(`{
						    "JSON": "bad,
						    "test": "good"
						}`)

	result := isValidJSON(goodJSON)

	assert := assert.New(t)
	assert.True(result)

	result = isValidJSON(badJSON)

	assert.False(result)

	result = isValidJSON([]byte(""))

	assert.False(result)
}
