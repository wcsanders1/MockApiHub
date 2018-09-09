package route

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"MockApiHub/str"
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
	if existingRoute, _ := tree.GetRoute(url); len(existingRoute) > 0 {
		return "", fmt.Errorf("route %s already registered", url)
	}

	url = strings.ToLower(url)
	fragments, err := str.GetURLFragments(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if err := tree.addRouteByFragments(fragments); err != nil {
		fmt.Println(err)
		return "", err
	}

	return tree.GetRoute(url)
}

// GetRoute returns a route if it exists in the tree
func (tree *Tree) GetRoute(url string) (string, error) {
	url = strings.ToLower(url)
	fragments, err := str.GetURLFragments(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tree.getRouteByFragments(fragments)
}

func (tree *Tree) getRouteByFragments(fragments []string) (string, error) {
	if len(fragments) == 0 {
		return "", nil
	}

	curFrag := fragments[0]
	remFrags := fragments[1:]
	notFoundError := NewHTTPError(notFoundMsg, http.StatusNotFound)

	if branch, exists := tree.branches[curFrag]; exists {
		route, err := branch.getRouteByFragments(remFrags)

		if err != nil {
			return "", err
		}

		if len(route) == 0 {
			if len(remFrags) == 0 {
				switch branch.routeType {
				case complete:
					return curFrag, nil
				case incomplete:
					return "", notFoundError
				default:
					return "", errors.New("problem finding route")
				}
			}

			return "", notFoundError
		}

		return fmt.Sprintf("%s/%s", curFrag, route), nil
	}

	params := tree.getRouteParamsInBranch()
	if len(params) == 0 {
		return "", notFoundError
	}

	for _, p := range params {
		if route, err := tree.branches[p].getRouteByFragments(remFrags); err == nil {
			return fmt.Sprintf("%s/%s", p, route), nil
		}
	}

	return "", notFoundError
}

func (tree *Tree) getRouteParamsInBranch() []string {
	var params []string
	for k := range tree.branches {
		if string(k[0]) == ":" {
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
