package manager

import (
	"strings"
	"fmt"
	"net/http"

	"MockApiHub/api"
)

const (
	refreshAPISPath = "refreshAllAPIs"
)

func (mgr *Manager) refreshAPI (w http.ResponseWriter, r *http.Request) {
	for _, api := range mgr.apis {
		if err := api.Shutdown(); err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	mgr.apis = make(map[string]*api.API)
	err := mgr.loadMockAPIs()
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr.registerMockAPIs()
	w.Write([]byte("successfully refreshed mock apis"))
}

func (mgr *Manager) registerHubAPIHandlers() {

	mgr.hubAPIHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodPost] = make(map[string]func(http.ResponseWriter, *http.Request))

	mgr.hubAPIHandlers[http.MethodPost][strings.ToLower(refreshAPISPath)] = mgr.refreshAPI
}