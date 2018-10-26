package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wcsanders1/MockApiHub/constants"
	"github.com/wcsanders1/MockApiHub/wrapper"
)

type (
	// AppConfig is application configuration
	AppConfig struct {
		HTTP HTTP
		Log  Log
	}

	// APIConfig is configuration for an individual mock api
	APIConfig struct {
		HTTP      HTTP
		BaseURL   string
		Endpoints map[string]Endpoint
		Log       Log
	}

	// Log is configuration for logging
	Log struct {
		LoggingEnabled bool
		Filename       string
		MaxFileSize    int
		MaxFileBackups int
		MaxFileDaysAge int
		FormatAsJSON   bool
		Level          string
		PrettyJSON     bool
	}

	// HTTP contains information regarding server setup
	HTTP struct {
		Port     int
		UseTLS   bool
		CertFile string
		KeyFile  string
	}

	// Endpoint contains information regarding an endpoint
	Endpoint struct {
		Path             string
		File             string
		Method           string
		EnforceValidJSON bool
	}

	// IManager provides functionality to manage configurations, such as getting
	// a mock API configuration from the disk
	IManager interface {
		GetAPIConfig(file os.FileInfo) (*APIConfig, error)
	}

	// Manager is a concrete implementation of IManager
	Manager struct {
		file wrapper.IFileOps
	}
)

// NewConfigManager returns a reference to a new Manager
func NewConfigManager() *Manager {
	return &Manager{
		file: &wrapper.FileOps{},
	}
}

// GetAPIConfig gets a mock API configuration from the disk
func (mgr *Manager) GetAPIConfig(fileInfo os.FileInfo) (*APIConfig, error) {
	dir := fileInfo.Name()
	if !fileInfo.IsDir() || !isAPI(dir) {
		return nil, errors.New("not a mock API directory")
	}

	apiConfig, err := mgr.getAPIConfigFromDir(dir)
	if err != nil {
		return nil, err
	}
	return apiConfig, nil
}

func (mgr *Manager) getAPIConfigFromDir(dir string) (*APIConfig, error) {
	files, _ := mgr.file.ReadDir(fmt.Sprintf("%s/%s", constants.APIDir, dir))
	for _, file := range files {
		if isAPIConfig(file.Name()) {
			apiConfig, err := mgr.decodeAPIConfig(dir, file.Name())
			if err != nil {
				return nil, err
			}
			return apiConfig, nil
		}
	}
	return nil, nil
}

func (mgr *Manager) decodeAPIConfig(dir string, fileName string) (*APIConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", constants.APIDir, dir, fileName)
	var Manager APIConfig
	if _, err := mgr.file.DecodeFile(path, &Manager); err != nil {
		return nil, err
	}
	return &Manager, nil
}

func isAPIConfig(fileName string) bool {
	return filepath.Ext(fileName) == ".toml"
}

func isAPI(dir string) bool {
	return len(dir) > 3 && dir[len(dir)-3:] == constants.APIDirExt
}
