package manager

import (
	"errors"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"

	"github.com/wcsanders1/MockApiHub/api"
	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewManager(t *testing.T) {
	port := 4000
	cfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: port,
		},
	}

	result, err := NewManager(cfg)

	assert := assert.New(t)
	assert.Nil(err)
	assert.NotNil(result)

	badCfg := &config.AppConfig{
		HTTP: config.HTTP{
			Port: 0,
		},
	}

	result, err = NewManager(badCfg)

	assert.Nil(result)
	assert.Error(err)
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

func TestApiByPortExists(t *testing.T) {
	port := 4000
	mgr := Manager{
		apis: make(map[string]*api.API),
	}

	cfg := &config.APIConfig{
		HTTP: config.HTTP{
			Port: port,
		},
	}

	mgr.apis["test"], _ = api.NewAPI(cfg)
	result := mgr.apiByPortExists(port)

	assert := assert.New(t)
	assert.True(result)

	result = mgr.apiByPortExists(port + 1)
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
	assert.IsType(&config.APIConfig{}, result)
	basicOpsPass.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))

	basicOpsFail := new(fake.BasicOps)
	basicOpsFail.On("DecodeFile", mock.AnythingOfType("string"), mock.Anything).Return(toml.MetaData{}, errors.New(""))

	mgrFail := Manager{
		file: basicOpsFail,
	}

	result2, err2 := mgrFail.decodeAPIConfig(dir, fileName)
	assert.Error(err2)
	assert.Nil(result2)
	assert.IsType(&config.APIConfig{}, result2)
	basicOpsFail.AssertCalled(t, "DecodeFile", path, mock.AnythingOfType("*config.APIConfig"))
}
