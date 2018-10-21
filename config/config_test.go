package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/wcsanders1/MockApiHub/constants"
	"github.com/wcsanders1/MockApiHub/helper"
	"github.com/wcsanders1/MockApiHub/wrapper"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConfigManager_ReturnsNewConfigManager_WhenCalled(t *testing.T) {
	result := NewConfigManager()

	assert := assert.New(t)
	assert.NotNil(result)
	assert.IsType(&Manager{}, result)
	assert.NotNil(result.file)
	assert.IsType(&wrapper.FileOps{}, result.file)
}

func TestIsAPIConfig_ReturnsTrue_WhenFileIsConfig(t *testing.T) {
	assert.True(t, isAPIConfig("test.toml"))
}

func TestIsAPIConfig_ReturnsFalse_WhenFileIsNotConfig(t *testing.T) {
	assert.False(t, isAPIConfig("test.exe"))
}

func TestIsAPIConfig_ReturnsFalse_WhenGivenNoFile(t *testing.T) {
	assert.False(t, isAPIConfig(""))
}

func TestIsAPI_ReturnsTrue_WhenIsAPIDirectory(t *testing.T) {
	assert.True(t, isAPI("testApi"))
}

func TestIsAPI_ReturnsFalse_WhenIsNotAPIDirectory(t *testing.T) {
	assert.False(t, isAPI("test"))
}

func TestIsAPI_ReturnsFalse_WhenNoDirectoryProvided(t *testing.T) {
	assert.False(t, isAPI(""))
}

func TestDecodeAPIConfig_ReturnsAPIConfig_WhenDecodeSuccessful(t *testing.T) {
	dir := "testDir"
	file := "testfile"
	path := fmt.Sprintf("%s/%s/%s", constants.APIDir, dir, file)
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.decodeAPIConfig(dir, file)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType(&APIConfig{}, result)
	fileOps.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))
}

func TestDecodeAPIConfig_ReturnsError_WhenDecodeFails(t *testing.T) {
	dir := "testDir"
	file := "testfile"
	path := fmt.Sprintf("%s/%s/%s", constants.APIDir, dir, file)
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))

	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.decodeAPIConfig(dir, file)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	assert.IsType(&APIConfig{}, result)
	fileOps.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))
}

func TestGetAPIConfig_ReturnsAPIConfig_WhenProvidedAPI(t *testing.T) {
	dir := "mockApi"
	fileInfo, fileCollection := helper.GetFakeFileInfoAndCollection(dir, "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.GetAPIConfig(fileInfo)

	expectedDir := fmt.Sprintf("%s/%s", constants.APIDir, dir)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	fileOps.AssertCalled(t, "ReadDir", expectedDir)
}

func TestGetAPIConfig_ReturnsError_WhenNotProvidedAPI(t *testing.T) {
	fileInfo, _ := helper.GetFakeFileInfoAndCollection("mockApiNot", "")
	fileOps := new(wrapper.FakeFileOps)
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.GetAPIConfig(fileInfo)

	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	fileOps.AssertNotCalled(t, "ReadDir", mock.AnythingOfType("string"))
	fileOps.AssertNotCalled(t, "DecodeFile", mock.AnythingOfType("string"), mock.Anything)
}

func TestGetAPIConfig_ReturnsError_WhenDecodeFails(t *testing.T) {
	dir := "mockApi"
	fileInfo, fileCollection := helper.GetFakeFileInfoAndCollection(dir, "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.GetAPIConfig(fileInfo)

	expectedDir := fmt.Sprintf("%s/%s", constants.APIDir, dir)
	assert := assert.New(t)
	assert.Error(err)
	assert.Nil(result)
	fileOps.AssertCalled(t, "ReadDir", expectedDir)
}

func TestGetAPIConfigFromDir_ReturnsAPIConfig_WhenProvidedAPI(t *testing.T) {
	dir := "testdir"
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.toml")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.getAPIConfigFromDir(dir)

	expectedDir := fmt.Sprintf("%s/%s", constants.APIDir, dir)
	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)
	assert.IsType(&APIConfig{}, result)
	fileOps.AssertCalled(t, "ReadDir", expectedDir)
}

func TestGetAPIConfigFromDir_ReturnsNil_WhenNoAPIsInDirectory(t *testing.T) {
	dir := "testdir"
	_, fileCollection := helper.GetFakeFileInfoAndCollection("", "testconfig.not")
	fileOps := new(wrapper.FakeFileOps)
	fileOps.On("ReadDir", mock.AnythingOfType("string")).Return(fileCollection, nil)
	fileOps.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, nil)
	mgr := Manager{
		file: fileOps,
	}

	result, err := mgr.getAPIConfigFromDir(dir)

	expectedDir := fmt.Sprintf("%s/%s", constants.APIDir, dir)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Nil(result)
	fileOps.AssertCalled(t, "ReadDir", expectedDir)
}
