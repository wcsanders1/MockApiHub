package utils

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
)

// GetJSON returns a []byte of valid JSON from a file
func GetJSON(filePath string) ([]byte, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()
	
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if (!isValidJSON(bytes)) {
		return nil, errors.New("invalid JSON")
	}

	return bytes, nil
}

func isValidJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}