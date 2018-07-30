package config

import (
	"fmt"
	"io/ioutil"
)

// Config is application configuration
type Config struct {
	HTTP http
	Apis map[string]api
}

type http struct {
	Port int
	SSL bool
}

type api struct {
	BaseURL string
	Endpoints []endpoint
}

type endpoint struct {
	Path string
	File string
}

// GetApis gets all the apis in the apis directory
func GetApis() error {
	files, err := ioutil.ReadDir("./config/apis")
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println("it's a dir")
			inFiles, _ := ioutil.ReadDir(fmt.Sprintf("./config/apis/%s", f.Name()))
			for _, innner := range inFiles {
				
				fmt.Println(innner.Name())
			}
		}
	}

	return nil
}