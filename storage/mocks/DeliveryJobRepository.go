// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	data "github.com/imyousuf/webhook-broker/storage/data"
	mock "github.com/stretchr/testify/mock"
)

// DeliveryJobRepository is an autogenerated mock type for the DeliveryJobRepository type
type DeliveryJobRepository struct {
	mock.Mock
}

// DispatchMessage provides a mock function with given fields: message, deliveryJobs
func (_m *DeliveryJobRepository) DispatchMessage(message *data.Message, deliveryJobs ...*data.DeliveryJob) error {
	_va := make([]interface{}, len(deliveryJobs))
	for _i := range deliveryJobs {
		_va[_i] = deliveryJobs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*data.Message, ...*data.DeliveryJob) error); ok {
		r0 = rf(message, deliveryJobs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}