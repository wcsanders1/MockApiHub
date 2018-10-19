package json

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wcsanders1/MockApiHub/wrapper"
)

// GetJSON returns a []byte of valid JSON from a file
func GetJSON(filePath string, file wrapper.IFileOps) ([]byte, error) {
	jsonFile, err := file.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	bytes, err := file.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !isValidJSON(bytes) {
		return nil, errors.New("invalid JSON")
	}

	return bytes, nil
}

func isValidJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}
