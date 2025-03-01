// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	model "trustwallet/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// EthereumClient is an autogenerated mock type for the EthereumClient type
type EthereumClient struct {
	mock.Mock
}

// GetLatestBlockNumber provides a mock function with given fields:
func (_m *EthereumClient) GetLatestBlockNumber() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockNumber")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsByBlockNumber provides a mock function with given fields: blockNumber
func (_m *EthereumClient) GetTransactionsByBlockNumber(blockNumber int64) ([]model.Transaction, error) {
	ret := _m.Called(blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionsByBlockNumber")
	}

	var r0 []model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) ([]model.Transaction, error)); ok {
		return rf(blockNumber)
	}
	if rf, ok := ret.Get(0).(func(int64) []model.Transaction); ok {
		r0 = rf(blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEthereumClient creates a new instance of EthereumClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEthereumClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *EthereumClient {
	mock := &EthereumClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
