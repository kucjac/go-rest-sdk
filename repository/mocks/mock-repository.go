// Code generated by mockery v1.0.0
package mocks

import dberrors "github.com/kucjac/go-rest-sdk/dberrors"
import mock "github.com/stretchr/testify/mock"
import repository "github.com/kucjac/go-rest-sdk/repository"

// MockRepository is an autogenerated mock type for the MockRepository type
type MockRepository struct {
	mock.Mock
}

// Count provides a mock function with given fields: req
func (_m *MockRepository) Count(req interface{}) (int, *dberrors.Error) {
	ret := _m.Called(req)

	var r0 int
	if rf, ok := ret.Get(0).(func(interface{}) int); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 *dberrors.Error
	if rf, ok := ret.Get(1).(func(interface{}) *dberrors.Error); ok {
		r1 = rf(req)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dberrors.Error)
		}
	}

	return r0, r1
}

// Create provides a mock function with given fields: req
func (_m *MockRepository) Create(req interface{}) *dberrors.Error {
	ret := _m.Called(req)

	var r0 *dberrors.Error
	if rf, ok := ret.Get(0).(func(interface{}) *dberrors.Error); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dberrors.Error)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: req, where
func (_m *MockRepository) Delete(req interface{}, where interface{}) *dberrors.Error {
	ret := _m.Called(req, where)

	var r0 *dberrors.Error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) *dberrors.Error); ok {
		r0 = rf(req, where)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dberrors.Error)
		}
	}

	return r0
}

// Get provides a mock function with given fields: req
func (_m *MockRepository) Get(req interface{}) (interface{}, *dberrors.Error) {
	ret := _m.Called(req)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 *dberrors.Error
	if rf, ok := ret.Get(1).(func(interface{}) *dberrors.Error); ok {
		r1 = rf(req)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dberrors.Error)
		}
	}

	return r0, r1
}

// List provides a mock function with given fields: req
func (_m *MockRepository) List(req interface{}) (interface{}, *dberrors.Error) {
	ret := _m.Called(req)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}) interface{}); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 *dberrors.Error
	if rf, ok := ret.Get(1).(func(interface{}) *dberrors.Error); ok {
		r1 = rf(req)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dberrors.Error)
		}
	}

	return r0, r1
}

// ListWithParams provides a mock function with given fields: req, params
func (_m *MockRepository) ListWithParams(req interface{}, params *repository.ListParameters) (interface{}, *dberrors.Error) {
	ret := _m.Called(req, params)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(interface{}, *repository.ListParameters) interface{}); ok {
		r0 = rf(req, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 *dberrors.Error
	if rf, ok := ret.Get(1).(func(interface{}, *repository.ListParameters) *dberrors.Error); ok {
		r1 = rf(req, params)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*dberrors.Error)
		}
	}

	return r0, r1
}

// Patch provides a mock function with given fields: req, where
func (_m *MockRepository) Patch(req interface{}, where interface{}) *dberrors.Error {
	ret := _m.Called(req, where)

	var r0 *dberrors.Error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) *dberrors.Error); ok {
		r0 = rf(req, where)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dberrors.Error)
		}
	}

	return r0
}

// Update provides a mock function with given fields: req
func (_m *MockRepository) Update(req interface{}) *dberrors.Error {
	ret := _m.Called(req)

	var r0 *dberrors.Error
	if rf, ok := ret.Get(0).(func(interface{}) *dberrors.Error); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dberrors.Error)
		}
	}

	return r0
}
