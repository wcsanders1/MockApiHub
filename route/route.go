package route

import (
	"fmt"
	"errors"

	"MockApiHub/str"
)

type (
	routeType int
	
	// Tree contains routing for an API in a tree format
	Tree struct {
		routeType routeType
		branches map[string]*Tree
	}
)

const (
	incomplete routeType = 0
	complete   routeType = 1
	wildCard   routeType = 2
)

// NewRouteTree returns a new instance of RouteTree
func NewRouteTree() *Tree {
	return &Tree {
		routeType: incomplete,
		branches: make(map[string]*Tree),
	}
}

// AddRoute adds a route to the tree
func (tree *Tree) AddRoute(url string) error {
	fragments, err := str.GetURLFragments(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := tree.addRouteByFragments(fragments); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// GetRoute returns a route if it exists in the tree
func (tree *Tree) GetRoute(url string) (string, error) {
	fragments, err := str.GetURLFragments(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tree.getRouteByFragments(fragments)
}

func (tree *Tree) getRouteByFragments(fragments []string) (string, error) {
	if (len(fragments) == 0) {
		return "", nil
	}

	curFrag := fragments[0]
	remFrags := fragments[1:]
	notFoundError := errors.New("route not found")

	if branch, exists := tree.branches[curFrag]; exists {
		route, err := branch.getRouteByFragments(remFrags)
		
		if err != nil {
			return "", err
		}

		if len(route) == 0 {
			if len(remFrags) == 0 {
				switch tree.routeType {
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

	return "", nil
}

func (tree *Tree) getRouteParamsInBranch() []string {
	var params []string
	for k := range tree.branches {
		if (string(k[0]) == ":") {
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
		if (tree.routeType == complete) {
			return errors.New("route already exists")
		}
		tree.routeType = complete
		
		return nil
	}

	return tree.addRouteByFragments(remFrags)
}