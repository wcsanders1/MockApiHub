package config

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wcsanders1/MockApiHub/fake"
)

func TestIsAPIConfig(t *testing.T) {
	fileName := "test.toml"
	result := isAPIConfig(fileName)

	assert := assert.New(t)
	assert.True(result)

	fileName = "test.exe"
	result = isAPIConfig(fileName)

	assert.False(result)

	fileName = ""
	result = isAPIConfig(fileName)

	assert.False(result)
}

func TestIsAPI(t *testing.T) {
	dirName := "testApi"
	result := isAPI(dirName)

	assert := assert.New(t)
	assert.True(result)

	dirName = "test"
	result = isAPI(dirName)

	assert.False(result)

	dirName = ""
	result = isAPI(dirName)

	assert.False(result)
}

func TestDecodeAPIConfig(t *testing.T) {
	basicOpsPass := new(fake.BasicOps)
	basicOpsPass.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrPass := Manager{
		file: basicOpsPass,
	}

	dir := "testDir"
	apiDir := "./api/apis"
	fileName := "testfile"
	path := fmt.Sprintf("%s/%s/%s", apiDir, dir, fileName)

	result, err := mgrPass.decodeAPIConfig(dir, fileName)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType(&APIConfig{}, result)
	basicOpsPass.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))

	basicOpsFail := new(fake.BasicOps)
	basicOpsFail.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))

	mgrFail := Manager{
		file: basicOpsFail,
	}

	result2, err2 := mgrFail.decodeAPIConfig(dir, fileName)
	assert.Error(err2)
	assert.Nil(result2)
	assert.IsType(&APIConfig{}, result2)
	basicOpsFail.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))
}

func TestGetAPIConfigFromDir(t *testing.T) {
	dir := "testdir"
	expectedDir := fmt.Sprintf("%s/%s", apiDir, dir)
	fileInfo := []os.FileInfo{}
	fileInfoPass := new(fake.FileInfo)
	fileInfoPass.On("Name").Return("testconfig.toml")
	fileInfo = append(fileInfo, fileInfoPass)

	basicOpsPass := new(fake.BasicOps)
	basicOpsPass.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfo, nil)
	basicOpsPass.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrPass := Manager{
		file: basicOpsPass,
	}

	result, err := mgrPass.getAPIConfigFromDir(dir)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType(&APIConfig{}, result)
	basicOpsPass.AssertCalled(t, "ReadDir", expectedDir)

	fileInfoNotAPIConfig := []os.FileInfo{}
	fileInfoNil := new(fake.FileInfo)
	fileInfoNil.On("Name").Return("testconfig.not")
	fileInfoNotAPIConfig = append(fileInfoNotAPIConfig, fileInfoNil)

	basicOpsNil := new(fake.BasicOps)
	basicOpsNil.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoNotAPIConfig, nil)
	basicOpsNil.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrNil := Manager{
		file: basicOpsNil,
	}

	result2, err2 := mgrNil.getAPIConfigFromDir(dir)
	assert.Nil(err2)
	assert.Nil(result2)
	basicOpsNil.AssertCalled(t, "ReadDir", expectedDir)

	basicOpsErr := new(fake.BasicOps)
	basicOpsErr.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfo, nil)
	basicOpsErr.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))

	mgrErr := Manager{
		file: basicOpsErr,
	}

	result3, err3 := mgrErr.getAPIConfigFromDir(dir)
	assert.Error(err3)
	assert.Nil(result3)
	basicOpsNil.AssertCalled(t, "ReadDir", expectedDir)
}
