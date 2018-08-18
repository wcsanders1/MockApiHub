package api

import (
	"os"
	"fmt"
	// "net/http"
	"io/ioutil"
	"encoding/json"
	"errors"

	// "github.com/labstack/echo"
)

func getJSON(filePath string) ([]byte, error) {
	fmt.Println(filePath)

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
		// return c.String(http.StatusInternalServerError, "bad json")
	}

	// return c.JSONBlob(http.StatusOK, bytes)
	return bytes, nil
}

func isValidJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}