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