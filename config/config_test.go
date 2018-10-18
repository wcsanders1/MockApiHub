package config

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/wcsanders1/MockApiHub/fake"
	"github.com/wcsanders1/MockApiHub/file"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConfigManager(t *testing.T) {
	result := NewConfigManager()

	assert := assert.New(t)
	assert.NotNil(result)
	assert.IsType(&Manager{}, result)
	assert.NotNil(result.file)
	assert.IsType(&file.BasicOps{}, result.file)
}

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
	basicOpsPass := new(file.FakeBasicOps)
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

	basicOpsFail := new(file.FakeBasicOps)
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

func TestGetAPIConfig(t *testing.T) {
	apiDirInner := "mockApi"
	expectedDir := fmt.Sprintf("%s/%s", apiDir, apiDirInner)
	fileInfoPass := new(fake.FileInfo)
	fileInfoPass.On("Name").Return(apiDirInner)
	fileInfoPass.On("IsDir").Return(true)

	fileInfoCollection := []os.FileInfo{}
	fileInfoInner := new(fake.FileInfo)
	fileInfoInner.On("Name").Return("testconfig.toml")
	fileInfoCollection = append(fileInfoCollection, fileInfoInner)

	basicOpsIsAPI := new(file.FakeBasicOps)
	basicOpsIsAPI.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)
	basicOpsIsAPI.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrIsAPI := Manager{
		file: basicOpsIsAPI,
	}

	result, err := mgrIsAPI.GetAPIConfig(fileInfoPass)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	basicOpsIsAPI.AssertCalled(t, "ReadDir", expectedDir)

	fileInfoNotAPI := new(fake.FileInfo)
	fileInfoNotAPI.On("IsDir").Return(true)
	fileInfoNotAPI.On("Name").Return("mockApiNot")

	basicOpsNotAPI := new(file.FakeBasicOps)
	basicOpsNotAPI.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)
	basicOpsNotAPI.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrNotAPI := Manager{
		file: basicOpsNotAPI,
	}

	result2, err2 := mgrNotAPI.GetAPIConfig(fileInfoNotAPI)
	assert.Error(err2)
	assert.Nil(result2)
	basicOpsNotAPI.AssertNotCalled(t, "ReadDir", mock.AnythingOfType("string"))
	basicOpsNotAPI.AssertNotCalled(t, "DecodeFile", mock.AnythingOfType("string"), mock.Anything)

	basicOpsDecodeErr := new(file.FakeBasicOps)
	basicOpsDecodeErr.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)
	basicOpsDecodeErr.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))

	mgrDecodeErr := Manager{
		file: basicOpsDecodeErr,
	}

	result3, err3 := mgrDecodeErr.GetAPIConfig(fileInfoPass)
	assert.Error(err3)
	assert.Nil(result3)
	basicOpsIsAPI.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoCollection, nil)
	basicOpsIsAPI.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)
}

func TestGetAPIConfigFromDir(t *testing.T) {
	dir := "testdir"
	expectedDir := fmt.Sprintf("%s/%s", apiDir, dir)
	fileInfo := []os.FileInfo{}
	fileInfoPass := new(fake.FileInfo)
	fileInfoPass.On("Name").Return("testconfig.toml")
	fileInfo = append(fileInfo, fileInfoPass)

	basicOpsPass := new(file.FakeBasicOps)
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

	basicOpsNil := new(file.FakeBasicOps)
	basicOpsNil.On("ReadDir", mock.AnythingOfType("string")).Return(fileInfoNotAPIConfig, nil)
	basicOpsNil.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)

	mgrNil := Manager{
		file: basicOpsNil,
	}

	result2, err2 := mgrNil.getAPIConfigFromDir(dir)
	assert.Nil(err2)
	assert.Nil(result2)
	basicOpsNil.AssertCalled(t, "ReadDir", expectedDir)

	basicOpsErr := new(file.FakeBasicOps)
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
