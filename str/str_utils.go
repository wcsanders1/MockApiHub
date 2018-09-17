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

// RemoveColonFromParam removes the colon from a route parameter so it looks nice when logged
func RemoveColonFromParam(param string) string {
	if len(param) == 0 {
		return ""
	}

	return param[1:]
}

// IsParam returns true if the string passed to it is a route parameter
func IsParam(routeFrag string) bool {
	if len(routeFrag) == 0 {
		return false
	}

	if string(routeFrag[0]) == ":" {
		return true
	}

	return false
}
