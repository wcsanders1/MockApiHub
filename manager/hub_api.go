package manager

import (
	"strings"
	"fmt"
	"net/http"
	"encoding/json"

	"MockApiHub/api"
	"MockApiHub/config"
)

type apiDisplay struct {
	BaseURL string
	Endpoints map[string]config.Endpoint
	Port int
}

const (
	refreshAPIsPath = "refresh-all-mock-apis"
	showAllAPIsPath = "show-all-registered-mock-apis"
)

func (mgr *Manager) refreshMockAPIs (w http.ResponseWriter, r *http.Request) {
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

func (mgr *Manager) showRegisteredMockAPIs (w http.ResponseWriter, r *http.Request) {
	apis := make(map[string]apiDisplay)
	for apiName, api := range mgr.apis {
		apis[apiName] = apiDisplay{
			BaseURL: api.GetBaseURL(),
			Port: api.GetPort(),
			Endpoints: api.GetEndpoints(),
		}
	}

	apisJSON, err := json.Marshal(apis)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(apisJSON)
}

func (mgr *Manager) registerHubAPIHandlers() {

	mgr.hubAPIHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodPost] = make(map[string]func(http.ResponseWriter, *http.Request))
	mgr.hubAPIHandlers[http.MethodGet] = make(map[string]func(http.ResponseWriter, *http.Request))

	mgr.hubAPIHandlers[http.MethodPost][strings.ToLower(refreshAPIsPath)] = mgr.refreshMockAPIs
	mgr.hubAPIHandlers[http.MethodGet][strings.ToLower(showAllAPIsPath)] = mgr.showRegisteredMockAPIs
}