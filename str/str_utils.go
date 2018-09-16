package str

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// GetPort returns the port in the format that http server expects
func GetPort(port int) string {
	return fmt.Sprintf(":%d", port)
}

// GetURLFragments returns the parts of a URL in an array
func GetURLFragments(url string) ([]string, error) {
	if len(url) == 0 {
		return nil, errors.New("no url provided")
	}

	frags := strings.Split(url, "/")

	return frags, nil
}

// CleanURL returns a URL in lowercase without a trailing or preceeding slash
func CleanURL(url string) string {
	if len(url) == 0 {
		return ""
	}

	urlLower := strings.ToLower(url)
	return path.Clean(urlLower[1:])
}
