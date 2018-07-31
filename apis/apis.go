package apis

import (
	"fmt"
	"io/ioutil"
)

// API contains the information needed to register the mock api
type API struct {
	BaseURL string
	Endpoints []endpoint
}

type endpoint struct {
	Path string
	File string
}

// GetApis gets all the apis in the apis directory
func GetApis() error {
	files, err := ioutil.ReadDir("./apis")
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			fName := f.Name()
			if fName[len(fName)-3:] == "Api" {
				fmt.Println("Found the following mock api: ", fName)
			
				inFiles, _ := ioutil.ReadDir(fmt.Sprintf("./apis/%s", f.Name()))
				for _, innner := range inFiles {
				
				fmt.Println(innner.Name())
			}
			}
			
		}
	}

	return nil
}