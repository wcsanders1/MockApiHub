package apis

import (
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/labstack/echo"
)

func getJSON(c echo.Context, filePath string) error {
	fmt.Println(filePath)

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	if (!isValidJSON(bytes)) {
		return c.String(http.StatusInternalServerError, "bad json")
	}

	return c.JSONBlob(http.StatusOK, bytes)
}

func isValidJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}