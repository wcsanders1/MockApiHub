/*
Package main is the main entry point for the application. It requires configuration in a file called app_config.toml.

Example configuration:

	[http]
	port = 5000
	useTLS = false
	certFile = ""
	keyFile = ""

	[log]
	loggingEnabled = true
	fileName = "testLogs/mockApiHub/default.log"
	maxFileDaysAge = 3
	formatAsJSON = true
	prettyJSON = true

*/
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/wcsanders1/MockApiHub/config"
	"github.com/wcsanders1/MockApiHub/manager"

	"github.com/BurntSushi/toml"
)

func main() {

	hubFlags := flag.NewFlagSet("hubFlags", flag.ExitOnError)
	showVersion := hubFlags.Bool("v", false, "application version")
	hubFlags.Parse(os.Args[1:])
	hubFlags.Visit(func(f *flag.Flag) {
		if f.Name == "v" {
			*showVersion = true
			fmt.Println("v0.4.0")
		}
	})

	if *showVersion {
		return
	}

	var appConfig config.AppConfig
	if _, err := toml.DecodeFile("app_config.toml", &appConfig); err != nil {
		fmt.Println(err)
		panic(err)
	}

	mgr, err := manager.NewManager(&appConfig)
	if err != nil {
		fmt.Println(err)
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
