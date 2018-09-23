package route

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/wcsanders1/MockApiHub/query"
	"github.com/wcsanders1/MockApiHub/str"
)

type (
	routeType int

	// Tree contains routing for an API in a tree format
	Tree struct {
		routeType routeType
		branches  map[string]*Tree
	}

	// HTTPError is the error type returned when an HTTP error occurs
	HTTPError struct {
		Msg    string
		Status int
	}
)

const (
	incomplete  routeType = 0
	complete    routeType = 1
	notFoundMsg string    = "route not found"
)

// NewHTTPError returns a reference to a NewHTTPError
func NewHTTPError(msg string, status int) *HTTPError {
	return &HTTPError{
		Msg:    fmt.Sprintf("msg: %s || status: %s", msg, http.StatusText(status)),
		Status: status,
	}
}

func (e *HTTPError) Error() string {
	return e.Msg
}

// NewRouteTree returns a new instance of Tree
func NewRouteTree() *Tree {
	return &Tree{
		routeType: incomplete,
		branches:  make(map[string]*Tree),
	}
}

// AddRoute adds a route to the tree
func (tree *Tree) AddRoute(url string) (string, error) {
	if len(url) == 0 {
		return "", errors.New("no url provided")
	}

	url = strings.ToLower(url)
	if existingRoute, _, _ := tree.GetRoute(url); len(existingRoute) > 0 {
		return "", fmt.Errorf("route %s already registered", url)
	}

	fragments, err := str.GetURLFragments(url)
	if err != nil {
		return "", err
	}

	if duplicateParamsExist(fragments) {
		return "", fmt.Errorf("route has duplicate parameters: %v", fragments)
	}

	if err := tree.addRouteByFragments(fragments); err != nil {
		return "", err
	}

	route, _, err := tree.GetRoute(url)
	if err != nil {
		return "", err
	}
	return path.Clean(route), nil
}

// GetRoute returns a route if it exists in the tree
func (tree *Tree) GetRoute(url string) (string, map[string]string, error) {
	url = strings.ToLower(url)
	fragments, err := str.GetURLFragments(url)
	if err != nil {
		return "", nil, err
	}

	route, params, err := tree.getRouteByFragments(fragments, make(map[string]string))
	if err != nil {
		return "", nil, err
	}
	return path.Clean(route), params, nil
}

func (tree *Tree) getRouteByFragments(fragments []string, params map[string]string) (string, map[string]string, error) {
	if len(fragments) == 0 {
		return "", params, nil
	}

	curFrag := fragments[0]
	remFrags := fragments[1:]
	notFoundError := NewHTTPError(notFoundMsg, http.StatusNotFound)

	if branch, exists := tree.branches[curFrag]; exists {
		route, params, err := branch.getRouteByFragments(remFrags, params)

		if err != nil {
			return "", params, err
		}

		if len(route) == 0 {
			if len(remFrags) == 0 {
				switch branch.routeType {
				case complete:
					return curFrag, params, nil
				case incomplete:
					return "", params, notFoundError
				default:
					return "", params, errors.New("problem finding route")
				}
			}
			return "", params, notFoundError
		}
		return fmt.Sprintf("%s/%s", curFrag, route), params, nil
	}

	newParams := tree.getRouteParamsInBranch()
	if len(newParams) == 0 {
		return "", params, notFoundError
	}

	for _, p := range newParams {
		if route, params, err := tree.branches[p].getRouteByFragments(remFrags, params); err == nil {
			param := str.RemoveColonFromParam(p)
			params[param] = curFrag
			return fmt.Sprintf("%s/%s", p, route), params, nil
		}
	}
	return "", params, notFoundError
}

func (tree *Tree) getRouteParamsInBranch() []string {
	var params []string
	for k := range tree.branches {
		if len(k) == 0 {
			return params
		}
		if str.IsParam(k) {
			params = append(params, k)
		}
	}
	return params
}

func (tree *Tree) addRouteByFragments(fragments []string) error {
	if len(fragments) == 0 {
		return nil
	}

	curFrag := fragments[0]
	remFrags := fragments[1:]

	if branch, exists := tree.branches[curFrag]; exists {
		return branch.addRouteToExistingBranch(remFrags)
	}

	tree.branches[curFrag] = NewRouteTree()
	if len(remFrags) > 0 {
		return tree.branches[curFrag].addRouteToExistingBranch(remFrags)
	}

	tree.branches[curFrag].routeType = complete

	return nil
}

func (tree *Tree) addRouteToExistingBranch(remFrags []string) error {
	if len(remFrags) == 0 {
		if tree.routeType == complete {
			return errors.New("route already exists")
		}
		tree.routeType = complete

		return nil
	}

	return tree.addRouteByFragments(remFrags)
}

func duplicateParamsExist(fragments []string) bool {
	var params []string
	for _, frag := range fragments {
		if str.IsParam(frag) {
			if query.ArrayContains(frag, params) {
				return true
			}
			params = append(params, frag)
		}
	}
	return false
}
