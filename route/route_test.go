package route

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wcsanders1/MockApiHub/str"

	"github.com/stretchr/testify/assert"
)

func TestNewRouteTree_ReturnsNewRouteTree_WhenCalled(t *testing.T) {
	result := NewRouteTree()

	assert := assert.New(t)
	assert.Equal(incomplete, result.routeType)
	assert.NotNil(result.branches)
}

func TestAddRoute_AddsRoute_WhenProvidedValidRoute(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	result, err := routeTree.AddRoute(route)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(route, result)
	assert.Contains(routeTree.branches, frags[0])
	assert.NotContains(routeTree.branches, frags[1])
	assert.Contains(routeTree.branches[frags[0]].branches, frags[1])
	assert.NotContains(routeTree.branches[frags[0]].branches, frags[0])
	assert.Empty(routeTree.branches[frags[0]].branches[frags[1]].branches)
	assert.Equal(incomplete, routeTree.routeType)
	assert.Equal(incomplete, routeTree.branches[frags[0]].routeType)
	assert.Equal(complete, routeTree.branches[frags[0]].branches[frags[1]].routeType)
}

func TestAddRoute_ReturnsError_WhenProvidedRegisteredRoute(t *testing.T) {
	route := "test/route"
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	result, err := routeTree.AddRoute(route)

	assert := assert.New(t)
	assert.Error(err)
	assert.Empty(result)
}

func TestAddRoute_ReturnsError_WhenProvidedRouteWithDuplicateParams(t *testing.T) {
	routeTree := NewRouteTree()

	result, err := routeTree.AddRoute("test/:param/customers/:param")

	assert := assert.New(t)
	assert.Error(err)
	assert.Empty(result)
}

func TestAddRoute_ReturnsError_WhenProvidedNothing(t *testing.T) {
	routeTree := NewRouteTree()

	result, err := routeTree.AddRoute("")

	assert := assert.New(t)
	assert.Error(err)
	assert.Empty(result)
}

func TestGetRoute(t *testing.T) {
	route1 := "test/route/again/and/again"
	routeTree := NewRouteTree()
	routeTree.AddRoute(route1)
	result, params, err := routeTree.GetRoute(route1)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route1, result)

	route2 := "test/route"
	result, params, err = routeTree.GetRoute(route2)

	assert.Error(err)
	assert.Empty(result)
	assert.Empty(params)

	routeTree.AddRoute(route2)

	result, params, err = routeTree.GetRoute(route2)

	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route2, result)

	route3 := "test/route/anotherroute"

	result, params, err = routeTree.GetRoute(route3)

	assert.Error(err)
	assert.Empty(params)
	assert.Empty(result)

	routeTree.AddRoute(route3)

	result, params, err = routeTree.GetRoute(route3)

	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route3, result)

	result, params, err = routeTree.GetRoute(route1)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route1, result)

	result, params, err = routeTree.GetRoute(route2)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route2, result)

	result, params, err = routeTree.GetRoute(route3)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route3, result)

	_, _, err = routeTree.GetRoute("test")
	assert.Error(err)

	url := "another/43434/route"

	_, params, err = routeTree.GetRoute(url)
	assert.Error(err)
	assert.Empty(params)

	route4 := "another/:param/route"
	routeTree.AddRoute(route4)

	result, params, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route4, result)
	assert.Contains(params, "param")
	assert.Equal("43434", params["param"])

	url = "another/3/route"
	result, params, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route4, result)
	assert.Contains(params, "param")
	assert.Equal("3", params["param"])

	route5 := "param/at/:end"
	routeTree.AddRoute(route5)

	url = "param/at/4325"
	result, params, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route5, result)
	assert.Contains(params, "end")
	assert.Equal("4325", params["end"])

	route6 := ":id/:another_id"
	routeTree.AddRoute(route6)

	url = "blah/blah"
	result, params, err = routeTree.GetRoute(url)
	assert.Nil(err)
	assert.Equal(route6, result)
	assert.Contains(params, "another_id")
	assert.Equal("blah", params["another_id"])

	emptyURL := ""
	emptyResult, params, err := routeTree.GetRoute(emptyURL)

	assert.Empty(emptyResult)
	assert.Nil(params)
	assert.NotNil(err)
	assert.Error(err)
}

func TestAddRouteToExistingBranch(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	r1 := routeTree.addRouteToExistingBranch(frags)

	assert := assert.New(t)
	assert.NotNil(r1)
	assert.Error(r1)

	r2 := routeTree.addRouteToExistingBranch(frags[:len(frags)-1])

	assert.Nil(r2)
	assert.Equal(complete, routeTree.branches[frags[0]].routeType)

	r3 := routeTree.addRouteToExistingBranch(frags[:len(frags)-1])

	assert.Error(r3)

	newRoute := fmt.Sprintf("%s/new", route)
	newFrags := strings.Split(newRoute, "/")

	r4 := routeTree.addRouteByFragments(newFrags)

	assert.Nil(r4)
	subTree := routeTree.branches[newFrags[0]].branches[newFrags[1]].branches[newFrags[2]]

	assert.Equal(complete, subTree.routeType)
}

func TestAddRouteByFragments(t *testing.T) {
	var frags []string
	routeTree := NewRouteTree()
	err := routeTree.addRouteByFragments(frags)

	assert := assert.New(t)
	assert.Nil(err)
}

func TestDuplicateParamsExist(t *testing.T) {
	noDups := "no/dup/:params/here"
	noDupFrags, _ := str.GetURLFragments(noDups)
	result := duplicateParamsExist(noDupFrags)

	assert := assert.New(t)
	assert.False(result)

	dups := "dup/:params/:params/here"
	dupFrags, _ := str.GetURLFragments(dups)
	result = duplicateParamsExist(dupFrags)

	assert.True(result)
}
