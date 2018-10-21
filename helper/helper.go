package helper

import (
	"os"

	"github.com/wcsanders1/MockApiHub/config"

	"github.com/wcsanders1/MockApiHub/fake"
)

// GetFakeFileInfoAndCollection returns fake os.FileInfo and a collection of os.FileInfo
// with that os.FileInfo in it, for testing.
func GetFakeFileInfoAndCollection(dir, file string) (os.FileInfo, []os.FileInfo) {
	fileInfo := new(fake.FileInfo)
	fileInfo.On("Name").Return(dir)
	fileInfo.On("IsDir").Return(true)

	fileInfoCollection := []os.FileInfo{}
	fileInfoInner := new(fake.FileInfo)
	fileInfoInner.On("Name").Return(file)
	fileInfoCollection = append(fileInfoCollection, fileInfoInner)

	return fileInfo, fileInfoCollection
}

// GetFakeAPIConfig returns a fake *config.APIConfig.
func GetFakeAPIConfig(port int) *config.APIConfig {
	return &config.APIConfig{
		HTTP: config.HTTP{
			Port: port,
		},
	}
}

// GetFakeAppConfig returns a fake *config.AppConfig
func GetFakeAppConfig(certFile, keyFile string) *config.AppConfig {
	return &config.AppConfig{
		HTTP: config.HTTP{
			CertFile: certFile,
			KeyFile:  keyFile,
		},
	}
}
