// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// DownloadArtifact provides a mock function with given fields: ctx, uri, version
func (_m *Client) DownloadArtifact(ctx context.Context, uri string, version string) (string, error) {
	ret := _m.Called(ctx, uri, version)

	if len(ret) == 0 {
		panic("no return value specified for DownloadArtifact")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, uri, version)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, uri, version)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, uri, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_DownloadArtifact_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DownloadArtifact'
type Client_DownloadArtifact_Call struct {
	*mock.Call
}

// DownloadArtifact is a helper method to define mock.On call
//   - ctx context.Context
//   - uri string
//   - version string
func (_e *Client_Expecter) DownloadArtifact(ctx interface{}, uri interface{}, version interface{}) *Client_DownloadArtifact_Call {
	return &Client_DownloadArtifact_Call{Call: _e.mock.On("DownloadArtifact", ctx, uri, version)}
}

func (_c *Client_DownloadArtifact_Call) Run(run func(ctx context.Context, uri string, version string)) *Client_DownloadArtifact_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Client_DownloadArtifact_Call) Return(_a0 string, _a1 error) *Client_DownloadArtifact_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_DownloadArtifact_Call) RunAndReturn(run func(context.Context, string, string) (string, error)) *Client_DownloadArtifact_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
