package apis

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// API contains information for an API
type API struct {
	BaseURL string
	Endpoints map[string]endpoint
}

type endpoint struct {
	Path string
	File string
}

// GetApis gets all the apis in the apis directory
func GetApis() (map[string]API, error) {
	files, err := ioutil.ReadDir("./apis")
	if err != nil {
		return nil, err
	}
	apis := make(map[string]API)
	for _, f := range files {
		if f.IsDir() {
			fName := f.Name()
			if fName[len(fName)-3:] == "Api" {
				fmt.Println("Found the following mock api: ", fName)
			
				inFiles, _ := ioutil.ReadDir(fmt.Sprintf("./apis/%s", f.Name()))
				for _, inner := range inFiles {
					ext := filepath.Ext(inner.Name())
					if (ext == ".toml") {
						var api API
						path := fmt.Sprintf("./apis/%s/%s", f.Name(), inner.Name())
						if _, err := toml.DecodeFile(path, &api); err != nil {
							fmt.Println(err)
							return nil, err
						}
						apis[fName] = api
						fmt.Println(inner.Name())
					}
				}
			}
		}
	}

	return apis, nil
}