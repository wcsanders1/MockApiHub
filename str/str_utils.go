package str

import (
	"fmt"
	"strings"
	"errors"
)

// GetPort returns the port in the format that http server expects
func GetPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

// GetURLFragments returns the parts of a URL in an array
func GetURLFragments(url string) ([]string, error) {
	if (len(url) == 0) {
		return nil, errors.New("no url provided")
	}

	frags := strings.Split(url, "/")
	
	return frags, nil
}