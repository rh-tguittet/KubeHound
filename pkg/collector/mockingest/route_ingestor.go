// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	types "github.com/DataDog/KubeHound/pkg/globals/types"
	mock "github.com/stretchr/testify/mock"
)

// RouteIngestor is an autogenerated mock type for the RouteIngestor type
type RouteIngestor struct {
	mock.Mock
}

type RouteIngestor_Expecter struct {
	mock *mock.Mock
}

func (_m *RouteIngestor) EXPECT() *RouteIngestor_Expecter {
	return &RouteIngestor_Expecter{mock: &_m.Mock}
}

// Complete provides a mock function with given fields: _a0
func (_m *RouteIngestor) Complete(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RouteIngestor_Complete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Complete'
type RouteIngestor_Complete_Call struct {
	*mock.Call
}

// Complete is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *RouteIngestor_Expecter) Complete(_a0 interface{}) *RouteIngestor_Complete_Call {
	return &RouteIngestor_Complete_Call{Call: _e.mock.On("Complete", _a0)}
}

func (_c *RouteIngestor_Complete_Call) Run(run func(_a0 context.Context)) *RouteIngestor_Complete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RouteIngestor_Complete_Call) Return(_a0 error) *RouteIngestor_Complete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RouteIngestor_Complete_Call) RunAndReturn(run func(context.Context) error) *RouteIngestor_Complete_Call {
	_c.Call.Return(run)
	return _c
}

// IngestRoute provides a mock function with given fields: _a0, _a1
func (_m *RouteIngestor) IngestRoute(_a0 context.Context, _a1 types.RouteType) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.RouteType) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RouteIngestor_IngestRoute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IngestRoute'
type RouteIngestor_IngestRoute_Call struct {
	*mock.Call
}

// IngestRoute is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 types.RouteType
func (_e *RouteIngestor_Expecter) IngestRoute(_a0 interface{}, _a1 interface{}) *RouteIngestor_IngestRoute_Call {
	return &RouteIngestor_IngestRoute_Call{Call: _e.mock.On("IngestRoute", _a0, _a1)}
}

func (_c *RouteIngestor_IngestRoute_Call) Run(run func(_a0 context.Context, _a1 types.RouteType)) *RouteIngestor_IngestRoute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.RouteType))
	})
	return _c
}

func (_c *RouteIngestor_IngestRoute_Call) Return(_a0 error) *RouteIngestor_IngestRoute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RouteIngestor_IngestRoute_Call) RunAndReturn(run func(context.Context, types.RouteType) error) *RouteIngestor_IngestRoute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRouteIngestor interface {
	mock.TestingT
	Cleanup(func())
}

// NewRouteIngestor creates a new instance of RouteIngestor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRouteIngestor(t mockConstructorTestingTNewRouteIngestor) *RouteIngestor {
	mock := &RouteIngestor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
