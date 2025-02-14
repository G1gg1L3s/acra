// Code generated by mockery v2.8.0. DO NOT EDIT.

package mocks

import (
	context "context"

	base "github.com/cossacklabs/acra/decryptor/base"

	mock "github.com/stretchr/testify/mock"

	net "net"
)

// ClientSession is an autogenerated mock type for the ClientSession type
type ClientSession struct {
	mock.Mock
}

// ClientConnection provides a mock function with given fields:
func (_m *ClientSession) ClientConnection() net.Conn {
	ret := _m.Called()

	var r0 net.Conn
	if rf, ok := ret.Get(0).(func() net.Conn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Conn)
		}
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *ClientSession) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// DatabaseConnection provides a mock function with given fields:
func (_m *ClientSession) DatabaseConnection() net.Conn {
	ret := _m.Called()

	var r0 net.Conn
	if rf, ok := ret.Get(0).(func() net.Conn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Conn)
		}
	}

	return r0
}

// DeleteData provides a mock function with given fields: _a0
func (_m *ClientSession) DeleteData(_a0 string) {
	_m.Called(_a0)
}

// GetData provides a mock function with given fields: _a0
func (_m *ClientSession) GetData(_a0 string) (interface{}, bool) {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// HasData provides a mock function with given fields: _a0
func (_m *ClientSession) HasData(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PreparedStatementRegistry provides a mock function with given fields:
func (_m *ClientSession) PreparedStatementRegistry() base.PreparedStatementRegistry {
	ret := _m.Called()

	var r0 base.PreparedStatementRegistry
	if rf, ok := ret.Get(0).(func() base.PreparedStatementRegistry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.PreparedStatementRegistry)
		}
	}

	return r0
}

// ProtocolState provides a mock function with given fields:
func (_m *ClientSession) ProtocolState() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// SetData provides a mock function with given fields: _a0, _a1
func (_m *ClientSession) SetData(_a0 string, _a1 interface{}) {
	_m.Called(_a0, _a1)
}

// SetPreparedStatementRegistry provides a mock function with given fields: registry
func (_m *ClientSession) SetPreparedStatementRegistry(registry base.PreparedStatementRegistry) {
	_m.Called(registry)
}

// SetProtocolState provides a mock function with given fields: state
func (_m *ClientSession) SetProtocolState(state interface{}) {
	_m.Called(state)
}
