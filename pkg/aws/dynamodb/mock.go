package dynamodb

import (
	"github.com/stretchr/testify/mock"
)

type DynamoMock struct {
	mock.Mock
}

func (mock *DynamoMock) Save(table string, item interface{}) error {
	args := mock.Called(item)
	return args.Error(0)
}

func (mock *DynamoMock) GetOne(table string, partitionKey string, bindTo interface{}) error {
	args := mock.Called(partitionKey, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error {
	args := mock.Called(partitionKey, sortKey, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error {
	args := mock.Called(partitionKey, limit, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) BatchGetWithSort(values map[string]interface{}) error {
	ret := mock.Called(values)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) error); ok {
		r0 = rf(values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryMultiple provides a mock function with given fields: table, partitionKey, limit, bindTo
func (_m *DynamoMock) QueryMultiple(table string, partitionKey string, limit int32, bindTo interface{}) error {
	ret := _m.Called(table, partitionKey, limit, bindTo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int32, interface{}) error); ok {
		r0 = rf(table, partitionKey, limit, bindTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
