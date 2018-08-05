package manager

import (
	"time"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"os"
	"context"

	"MockApiHub/api"
	"MockApiHub/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/BurntSushi/toml"
)

// Manager coordinates and controls the apis
type Manager struct{
	apis map[string]*api.API
	server *echo.Echo
	config *config.Config
}

const apiDir = "./api/apis"

// NewManager returns an instance of the Manager type
func NewManager(config *config.Config) *Manager {
	return &Manager{
		config: config,
	}
}

// StartMockAPIHub registers the mock apis and serves them
func (mgr *Manager) StartMockAPIHub() error {
	mgr.clearMockAPIs()
	mgr.initializeServer()
	err := mgr.loadMockAPIs()
	if err != nil {
		return err
	}

	mgr.registerMockAPIs()
	mgr.registerUtilityEndpoints()
	err = mgr.startServer()
	if err != nil {
		return err
	}

	return nil
}

func (mgr *Manager) initializeServer() {
	server := echo.New()
	server.Use(middleware.Logger())

	mgr.server = server
}

func (mgr *Manager) clearMockAPIs() {
	mgr.apis = make(map[string]*api.API)
}

func (mgr *Manager) registerUtilityEndpoints() error {
	mgr.server.POST("refresh", mgr.refreshMockAPIRegistry)
	
	return nil
}

func (mgr *Manager) refreshMockAPIRegistry(ctx echo.Context) error {
	if err := mgr.refreshServer(); err != nil {
		mgr.server.Logger.Fatal(err)
		return err
	}

	mgr.clearMockAPIs()
	if err := mgr.loadMockAPIs(); err != nil {
		return err
	}

	mgr.registerMockAPIs()
	if err := mgr.startServer(); err != nil {
		return err
	}

	return nil
}

func (mgr *Manager) stopServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	if err := mgr.server.Shutdown(ctx); err != nil {
		fmt.Println("going in here")
		return err
	}

	return nil
}

func (mgr *Manager) refreshServer() error {
	if err := mgr.stopServer(); err != nil {
		return err
	}
	mgr.initializeServer()

	return nil
}

func (mgr *Manager) startServer() error {
	addr := fmt.Sprintf(":%d", mgr.config.HTTP.Port)
	if mgr.config.HTTP.UseTLS {
		startUsingTLS(mgr.server, &mgr.config.HTTP, addr)
	} else {
		mgr.server.Logger.Fatal(mgr.server.Start(addr))
	}

	return nil
}

func startUsingTLS(server *echo.Echo, http *config.HTTP, addr string) {
	if _, err := os.Stat(http.CertFile); os.IsNotExist(err) {
		server.Logger.Fatal(fmt.Sprintf("%s cert file does not exist", http.CertFile))
	}

	if _, err := os.Stat(http.KeyFile); os.IsNotExist(err) {
		server.Logger.Fatal(fmt.Sprintf("%s key file does not exist", http.KeyFile))
	}

	server.Logger.Fatal(server.StartTLS(addr, http.CertFile, http.KeyFile))
}

func (mgr *Manager) registerMockAPIs() {
	for dir, api := range mgr.apis {
		err := api.Register(dir, mgr.server)
		if err != nil {
			mgr.server.Logger.Error(err, fmt.Sprintf("Error regisering the %s API.", dir))
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
			return decodeAPIConfig(dir, file)
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