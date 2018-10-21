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

func TestGetRoute_ReturnsRoute_WhenOneRouteRegistered(t *testing.T) {
	route := "test/route/again/and/again"
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	result, params, err := routeTree.GetRoute(route)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route, result)
}

func TestGetRoute_ReturnsCorrectRoute_WhenTwoRoutesRegistered(t *testing.T) {
	route := "test/route"
	routeTree := NewRouteTree()
	routeTree.AddRoute("test/route/again/and/again")
	routeTree.AddRoute(route)

	result, params, err := routeTree.GetRoute(route)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(route, result)
	assert.Empty(params)
}

func TestGetRoute_ReturnsCorrectRoute_WhenThreeRoutesRegistered(t *testing.T) {
	route := "test/route/anotherroute"
	routeTree := NewRouteTree()
	routeTree.AddRoute("test/route/again/and/again")
	routeTree.AddRoute("test/route")
	routeTree.AddRoute(route)

	result, params, err := routeTree.GetRoute(route)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(params)
	assert.Equal(route, result)
}

func TestGetRoute_ReturnsError_WhenProvidedNonRegisteredRoute(t *testing.T) {
	routeTree := NewRouteTree()
	routeTree.AddRoute("test/route/again/and/again")
	routeTree.AddRoute("test/route")
	routeTree.AddRoute("test/route/anotherroute")

	result, params, err := routeTree.GetRoute("another/43434/route")

	assert := assert.New(t)
	assert.Error(err)
	assert.Empty(params)
	assert.Empty(result)
}

func TestGetRoute_ReturnsRouteWithParam_WhenRouteHasParam(t *testing.T) {
	paramKey := "param"
	paramVal := "43434"
	route := fmt.Sprintf("another/:%s/route", paramKey)
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	result, params, err := routeTree.GetRoute(fmt.Sprintf("another/%s/route", paramVal))

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(route, result)
	assert.Contains(params, paramKey)
	assert.Equal(paramVal, params[paramKey])
}

func TestGetRoute_ReturnsRouteWithParam_WhenRouteHasParamAtEnd(t *testing.T) {
	paramKey := "end"
	paramVal := "4325"
	route := fmt.Sprintf("param/at/:%s", paramKey)
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	result, params, err := routeTree.GetRoute(fmt.Sprintf("param/at/%s", paramVal))

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(route, result)
	assert.Contains(params, paramKey)
	assert.Equal(paramVal, params[paramKey])
}

func TestGetRoute_ReturnsRouteWithParams_WhenRouteIsOnlyParams(t *testing.T) {
	accountParamKey := "account"
	accountParamVal := "3399"
	idParamKey := "id"
	idParamVal := "654h$76"
	anotherIDParamKey := "another_id"
	anotherIDParamVal := "9**11"
	routeNoParams := "just/normal/route"
	routeOneParam := fmt.Sprintf("just/:%s/param", accountParamKey)
	routeOnlyParams := fmt.Sprintf(":%s/:%s", idParamKey, anotherIDParamKey)
	routeTree := NewRouteTree()
	routeTree.AddRoute(routeNoParams)
	routeTree.AddRoute(routeOneParam)
	routeTree.AddRoute(routeOnlyParams)

	noParamRetrievedResult, noParamRetrievedParams, noParamErr := routeTree.GetRoute(routeNoParams)
	oneParamRetrievedResult, oneParamRetrievedParams, oneParamErr := routeTree.GetRoute(fmt.Sprintf("just/%s/param", accountParamVal))
	onlyParamsRetrievedResult, onlyParamsRetrievedParams, onlyParamsErr := routeTree.GetRoute(fmt.Sprintf("%s/%s", idParamVal, anotherIDParamVal))

	assert := assert.New(t)
	assert.Nil(noParamErr)
	assert.Nil(oneParamErr)
	assert.Nil(onlyParamsErr)
	assert.Equal(routeNoParams, noParamRetrievedResult)
	assert.Equal(routeOneParam, oneParamRetrievedResult)
	assert.Equal(routeOnlyParams, onlyParamsRetrievedResult)
	assert.Empty(noParamRetrievedParams)
	assert.Contains(oneParamRetrievedParams, accountParamKey)
	assert.Equal(accountParamVal, oneParamRetrievedParams[accountParamKey])
	assert.Contains(onlyParamsRetrievedParams, idParamKey)
	assert.Contains(onlyParamsRetrievedParams, anotherIDParamKey)
	assert.Equal(idParamVal, onlyParamsRetrievedParams[idParamKey])
	assert.Equal(anotherIDParamVal, onlyParamsRetrievedParams[anotherIDParamKey])
}

func TestGetRoute_ReturnsError_WhenProvidedNothing(t *testing.T) {
	routeTree := NewRouteTree()

	result, params, err := routeTree.GetRoute("")

	assert := assert.New(t)
	assert.Empty(result)
	assert.Nil(params)
	assert.Error(err)
}

func TestAddRouteToExistingBranch_ReturnsNil_WhenRouteValid(t *testing.T) {
	frags := strings.Split("test/route", "/")
	routeTree := NewRouteTree()

	err := routeTree.addRouteToExistingBranch(frags)

	assert.Nil(t, err)
}

func TestAddRouteToExistingBranch_ReturnsError_WhenProvidedRegisteredRoute(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	err := routeTree.addRouteToExistingBranch(frags)

	assert.Error(t, err)
}

func TestAddRouteToExistingBranch_RegistersRoute_OneStepBelowExistingBranch(t *testing.T) {
	route := "test/route"
	frags := strings.Split(route, "/")
	routeTree := NewRouteTree()
	routeTree.AddRoute(route)

	err := routeTree.addRouteToExistingBranch(frags[:len(frags)-1])

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(complete, routeTree.branches[frags[0]].routeType)
}

func TestDuplicateParamsExist_ReturnsFalse_WhenNoDuplicateParams(t *testing.T) {
	frags, _ := str.GetURLFragments("no/dup/:params/here")

	result := duplicateParamsExist(frags)

	assert.False(t, result)
}

func TestDuplicateParamsExist_ReturnsTrue_WhenDuplicateParamsExist(t *testing.T) {
	frags, _ := str.GetURLFragments("dup/:params/:params/here")

	result := duplicateParamsExist(frags)

	assert.True(t, result)
}
