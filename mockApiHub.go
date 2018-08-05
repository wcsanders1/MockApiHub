package main

import (
	"fmt"

	"MockApiHub/config"
	"MockApiHub/manager"

	"github.com/BurntSushi/toml"
)

func main() {
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	
	mgr := manager.NewManager(&config)
	mgr.StartMockAPIHub()
}