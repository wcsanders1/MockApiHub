package json

import (
	"errors"
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

func (ops *fakeBasicOps) ReadDir(dir string) ([]os.FileInfo, error) {
	args := ops.Called(dir)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

func TestGetJSON(t *testing.T) {
	filePath := "testpath"
	basicOpsPass := new(fakeBasicOps)
	basicOpsPass.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	basicOpsPass.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, nil)

	result, err := GetJSON(filePath, basicOpsPass)

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Nil(err)
	basicOpsPass.AssertCalled(t, "Open", filePath)
	basicOpsPass.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))

	basicOpsOpenFail := new(fakeBasicOps)
	basicOpsOpenFail.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), errors.New(""))

	result2, err2 := GetJSON(filePath, basicOpsOpenFail)

	assert.NotNil(err2)
	assert.Nil(result2)
	basicOpsOpenFail.AssertCalled(t, "Open", filePath)
	basicOpsOpenFail.AssertNotCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))

	basicOpsReadAllFail := new(fakeBasicOps)
	basicOpsReadAllFail.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	basicOpsReadAllFail.On("ReadAll", mock.AnythingOfType("*os.File")).Return(goodJSON, errors.New(""))

	result3, err3 := GetJSON(filePath, basicOpsReadAllFail)
	assert.NotNil(err3)
	assert.Nil(result3)
	basicOpsReadAllFail.AssertCalled(t, "Open", filePath)
	basicOpsReadAllFail.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))

	basicOpsBadJSON := new(fakeBasicOps)
	basicOpsBadJSON.On("Open", mock.AnythingOfType("string")).Return(os.NewFile(1, "fakefile"), nil)
	basicOpsBadJSON.On("ReadAll", mock.AnythingOfType("*os.File")).Return(badJSON, nil)

	result4, err4 := GetJSON(filePath, basicOpsBadJSON)
	assert.NotNil(err4)
	assert.Nil(result4)
	basicOpsBadJSON.AssertCalled(t, "Open", filePath)
	basicOpsBadJSON.AssertCalled(t, "ReadAll", mock.AnythingOfType("*os.File"))
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
