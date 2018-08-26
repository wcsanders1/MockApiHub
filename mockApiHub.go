package main

import (
	"fmt"

	"MockApiHub/config"
	"MockApiHub/manager"

	"github.com/BurntSushi/toml"
)
func main() {
	var appConfig config.AppConfig
	if _, err := toml.DecodeFile("app_config.toml", &appConfig); err != nil {
		fmt.Println(err)
		return
	}

	mgr := manager.NewManager(&appConfig)
	mgr.StartMockAPIHub()
}