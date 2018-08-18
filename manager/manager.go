package manager

import (
	"net/http"
	"io"	
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"

	"MockApiHub/api"
	"MockApiHub/config"
	"MockApiHub/utils"

	//"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

// Manager coordinates and controls the apis
type Manager struct{
	apis map[string]*api.API
	config *config.Config
	server *http.Server
}

const apiDir = "./api/apis"

// 	thenewstack.io/building-a-web-server-in-go

// NewManager returns an instance of the Manager type
func NewManager(config *config.Config) *Manager {
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

func createManagerServer(httpConfig *config.HTTP) (*http.Server, error) {
	server := &http.Server {
		Addr: utils.GetPort(httpConfig.Port),
		Handler: http.HandlerFunc(handler),
	}

	return server, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	if (req.Method == http.MethodGet) {
		io.WriteString(w, "hi")
	}		
}

// func getPort(port int) string {
// 	return fmt.Sprintf(":%d", port)
// }

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	// mgr.initializeServer()
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
			// mgr.server.Logger.Error(err, fmt.Sprintf("Error regisering the %s API.", dir))
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
		api, err := extractAPIFromFile(file)
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

func extractAPIFromFile(file os.FileInfo) (*api.API, error) {
	if (!file.IsDir() || !isAPI(file.Name())) {
		return nil, nil
	}
	
	dir := file.Name()
	fmt.Println("Found the following mock api: ", dir)
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", apiDir, dir))
	api, err := getAPI(dir, files)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return api, nil	
}

func getAPI(dir string, files []os.FileInfo) (*api.API, error) {
	for _, file := range files {
		if (isAPIConfig(file)) {
			api, err:= decodeAPIConfig(dir, file)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return api, nil
		}
	}
	return nil, nil
}

func decodeAPIConfig(dir string, file os.FileInfo) (*api.API, error) {
	path := fmt.Sprintf("%s/%s/%s", apiDir, dir, file.Name())
	var api api.API
	if _, err := toml.DecodeFile(path, &api); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &api, nil
}

func isAPI(dir string) bool {
	return len(dir) > 3 && dir[len(dir)-3:] == "Api"
}

func isAPIConfig(file os.FileInfo) bool {
	ext := filepath.Ext(file.Name())
	return ext == ".toml"
}