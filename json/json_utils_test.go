package json

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeBasicOps struct{}

var goodJSON = []byte(`{
	"JSON": "good",
	"test": "good"
}`)
var badJSON = []byte(`{
	"JSON": "bad,
	"test": "good"
}`)

func (ops *fakeBasicOps) Open(file string) (*os.File, error) {
	return os.NewFile(1, "fakefile"), nil
}

func (ops *fakeBasicOps) ReadAll(file *os.File) ([]byte, error) {
	return goodJSON, nil
}

// TODO: Finish this
func TestGetJSON(t *testing.T) {

}

func TestIsValidJSON(t *testing.T) {
	result := isValidJSON(goodJSON)

	assert := assert.New(t)
	assert.True(result)

	result = isValidJSON(badJSON)

	assert.False(result)

	result = isValidJSON([]byte(""))

	assert.False(result)
}
