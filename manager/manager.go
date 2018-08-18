package manager

import (
	"net/http"
	"io"	
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"errors"

	"MockApiHub/api"
	"MockApiHub/config"
	"MockApiHub/utils"

	"github.com/BurntSushi/toml"
)

// Manager coordinates and controls the apis
type Manager struct{
	apis map[string]*api.API
	config *config.AppConfig
	server *http.Server
}

const apiDir = "./api/apis"

// 	thenewstack.io/building-a-web-server-in-go

// NewManager returns an instance of the Manager type
func NewManager(config *config.AppConfig) *Manager {
	server, err := createManagerServer(&config.HTTP)
	if err != nil {
		fmt.Println(err)
	}

	return &Manager{
		config: config,
		server: server,
		apis: make(map[string]*api.API),
	}
}

func createManagerServer(config *config.HTTP) (*http.Server, error) {
	if config.Port == 0 {
		return nil, errors.New("no port provided")
	}

	server := &http.Server {
		Addr: utils.GetPort(config.Port),
		Handler: http.HandlerFunc(handler),
	}

	return server, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	if (req.Method == http.MethodGet) {
		io.WriteString(w, "hi")
	}		
}

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	err := mgr.loadMockAPIs()
	if err != nil {
		return err
	}

	mgr.registerMockAPIs()
	mgr.startHubServer()

	return nil
}

func (mgr *Manager) startHubServer() error {
	if mgr.config.HTTP.UseTLS {
		return mgr.startHubServerUsingTLS()
	}
	mgr.server.ListenAndServe()
	
	return nil
}

func (mgr *Manager) startHubServerUsingTLS() error {
	certFile := mgr.config.HTTP.CertFile
	keyFile := mgr.config.HTTP.KeyFile
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("%s cert file does not exist", certFile))
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("%s key file does not exist", keyFile))
	}
	return mgr.server.ListenAndServeTLS(certFile, keyFile)
}

func (mgr *Manager) registerMockAPIs() {
	for dir, api := range mgr.apis {
		err := api.Register(dir)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (mgr *Manager) loadMockAPIs() error {
	files, err := ioutil.ReadDir(apiDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		apiConfig, err := getAPIConfig(file)
		if err != nil {
			fmt.Println(err)
			return err
		}

		api, err := api.NewAPI(apiConfig)
		if err != nil {
			fmt.Println(err)
			return err
		}
		
		if api != nil {
			mgr.apis[file.Name()] = api
		}
	}

	return nil
}

func getAPIConfig(file os.FileInfo) (*config.APIConfig, error) {
	if (!file.IsDir() || !isAPI(file.Name())) {
		return nil, nil
	}
	
	dir := file.Name()
	fmt.Println("Found the following mock api: ", dir)
	apiConfig, err := getAPIConfigFromDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return apiConfig, nil	
}

func getAPIConfigFromDir(dir string) (*config.APIConfig, error) {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", apiDir, dir))
	for _, file := range files {
		if (isAPIConfig(file)) {
			apiConfig, err:= decodeAPIConfig(dir, file)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return apiConfig, nil
		}
	}
	return nil, nil
}

func decodeAPIConfig(dir string, file os.FileInfo) (*config.APIConfig, error) {
	path := fmt.Sprintf("%s/%s/%s", apiDir, dir, file.Name())
	var config config.APIConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &config, nil
}

func isAPI(dir string) bool {
	return len(dir) > 3 && dir[len(dir)-3:] == "Api"
}

func isAPIConfig(file os.FileInfo) bool {
	ext := filepath.Ext(file.Name())
	return ext == ".toml"
}