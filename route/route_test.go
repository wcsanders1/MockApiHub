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

func TestGetRoute(t *testing.T) {
	route1 := "test/route/again/and/again"
	routeTree := NewRouteTree()
	routeTree.AddRoute(route1)
	result, err := routeTree.GetRoute(route1)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(route1, result)

	route2 := "test/route"
	result, err = routeTree.GetRoute(route2)

	assert.Error(err)
	assert.Empty(result)
	
	routeTree.AddRoute(route2)

	result, err = routeTree.GetRoute(route2)

	assert.Nil(err)
	assert.Equal(route2, result)

	route3 := "test/route/anotherroute"

	result, err = routeTree.GetRoute(route3)

	assert.Error(err)
	assert.Empty(result)

	routeTree.AddRoute(route3)

	result, err = routeTree.GetRoute(route3)

	assert.Nil(err)
	assert.Equal(route3, result)

	result, err = routeTree.GetRoute(route1)
	assert.Nil(err)
	assert.Equal(route1, result)

	result, err = routeTree.GetRoute(route2)
	assert.Nil(err)
	assert.Equal(route2, result)

	result, err = routeTree.GetRoute(route3)
	assert.Nil(err)
	assert.Equal(route3, result)

	_, err = routeTree.GetRoute("test")
	assert.Error(err)

	url := "another/43434/route"

	_, err = routeTree.GetRoute(url)
	assert.Error(err)

	route4 := "another/:param/route"
	routeTree.AddRoute(route4)

	result, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route4, result)

	url = "another/3/route"
	result, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route4, result)

	route5 := "another/param/route"
	routeTree.AddRoute(route5)

	result, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route4, result)

	result, err = routeTree.GetRoute(route5)
	assert.Nil(err)
	assert.Equal(route5, result)
}