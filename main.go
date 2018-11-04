package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/manager"

	"github.com/BurntSushi/toml"
)

func main() {
	var appConfig config.AppConfig
	if _, err := toml.DecodeFile("app_config.toml", &appConfig); err != nil {
		fmt.Println(err)
		return
	}

	mgr, err := manager.NewManager(&appConfig)
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		if err := mgr.StartMockAPIHub(); err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()

	<-shutdown
	mgr.StopMockAPIHub()
}
