// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	grpc "github.com/Haba1234/sysmon/internal/grpc/api"
	mock "github.com/stretchr/testify/mock"
)

// StatisticsServer is an autogenerated mock type for the StatisticsServer type
type StatisticsServer struct {
	mock.Mock
}

// ListStatistics provides a mock function with given fields: _a0, _a1
func (_m *StatisticsServer) ListStatistics(_a0 *grpc.SubscriptionRequest, _a1 grpc.Statistics_ListStatisticsServer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*grpc.SubscriptionRequest, grpc.Statistics_ListStatisticsServer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mustEmbedUnimplementedStatisticsServer provides a mock function with given fields:
func (_m *StatisticsServer) mustEmbedUnimplementedStatisticsServer() {
	_m.Called()
}