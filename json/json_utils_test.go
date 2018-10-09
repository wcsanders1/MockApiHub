package json

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeBasicOps struct {
	mock.Mock
}

var goodJSON = []byte(`{
	"JSON": "good",
	"test": "good"
}`)
var badJSON = []byte(`{
	"JSON": "bad,
	"test": "good"
}`)

func (ops *fakeBasicOps) ReadAll(file *os.File) ([]byte, error) {
	args := ops.Called(file)
	return args.Get(0).([]byte), args.Error(1)
}

func (ops *fakeBasicOps) Open(file string) (*os.File, error) {
	args := ops.Called(file)
	return args.Get(0).(*os.File), args.Error(1)
}

func TestGetJSON(t *testing.T) {
	filePath := "testpath"
	basicOps := new(fakeBasicOps)
	basicOps.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	basicOps.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, nil)

	result, err := GetJSON(filePath, basicOps)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Nil(err)
	basicOps.AssertCalled(t, "Open", filePath)
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
