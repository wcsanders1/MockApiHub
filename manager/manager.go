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

	//"github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

// Manager coordinates and controls the apis
type Manager struct{
	apis map[string]*api.API
	config *config.Config
	realServer *http.Server
}

const apiDir = "./api/apis"

// NewManager returns an instance of the Manager type
func NewManager(config *config.Config) *Manager {
	server, err := createManagerServer(&config.HTTP)
	if err != nil {
		fmt.Println(err)
	}

	return &Manager{
		config: config,
		realServer: server,
		apis: make(map[string]*api.API),
	}
}

func createManagerServer(httpConfig *config.HTTP) (*http.Server, error) {
	server := &http.Server {
		Addr: getPort(httpConfig.Port),
		Handler: http.HandlerFunc(handler),
	}

	// 	thenewstack.io/building-a-web-server-in-go

	return server, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	if (req.Method == http.MethodGet) {
		io.WriteString(w, "hi")
	}		
}

func getPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	// mgr.initializeServer()
	err := mgr.loadMockAPIs()
	if err != nil {
		return err
	}

	mgr.registerMockAPIs()

	// err = mgr.startServer()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// func (mgr *Manager) initializeServer() {
// 	server := echo.New()
// 	server.Use(middleware.Logger())

// 	mgr.server = server
// }

func (mgr *Manager) startServer() error {
	// addr := fmt.Sprintf(":%d", mgr.config.HTTP.Port)
	// if mgr.config.HTTP.UseTLS {
	// 	startUsingTLS(mgr.server, &mgr.config.HTTP, addr)
	// } else {
	// 	mgr.server.Logger.Fatal(mgr.server.Start(addr))
	// }

	return nil
}

func startUsingTLS(http *config.HTTP, addr string) {
	// if _, err := os.Stat(http.CertFile); os.IsNotExist(err) {
	// 	server.Logger.Fatal(fmt.Sprintf("%s cert file does not exist", http.CertFile))
	// }

	// if _, err := os.Stat(http.KeyFile); os.IsNotExist(err) {
	// 	server.Logger.Fatal(fmt.Sprintf("%s key file does not exist", http.KeyFile))
	// }

	// server.Logger.Fatal(server.StartTLS(addr, http.CertFile, http.KeyFile))
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