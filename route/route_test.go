package route

import (
	"strings"
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestAddRouteToExistingBranch(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	r1 := routeTree.addRouteToExistingBranch(frags)

	assert := assert.New(t)
	assert.NotNil(r1)
	assert.Error(r1)

	r2 := routeTree.addRouteToExistingBranch(frags[:len(frags) - 1])

	assert.Nil(r2)
	assert.Equal(complete, routeTree.branches[frags[0]].routeType)

	r3 := routeTree.addRouteToExistingBranch(frags[:len(frags) - 1])

	assert.Error(r3)

	newRoute := fmt.Sprintf("%s/new", route)
	newFrags := strings.Split(newRoute, "/")

	r4 := routeTree.addRouteByFragments(newFrags)

	assert.Nil(r4)
	subTree := routeTree.branches[newFrags[0]].branches[newFrags[1]].branches[newFrags[2]]

	assert.Equal(complete, subTree.routeType)
}

func TestAddRoute(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	assert := assert.New(t)
	assert.Contains(routeTree.branches, frags[0])
	assert.NotContains(routeTree.branches, frags[1])
	assert.Contains(routeTree.branches[frags[0]].branches, frags[1])
	assert.NotContains(routeTree.branches[frags[0]].branches, frags[0])
	assert.Empty(routeTree.branches[frags[0]].branches[frags[1]].branches)
	assert.Equal(incomplete, routeTree.routeType)
	assert.Equal(incomplete, routeTree.branches[frags[0]].routeType)
	assert.Equal(complete, routeTree.branches[frags[0]].branches[frags[1]].routeType)
}