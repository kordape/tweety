// Code generated by mockery. DO NOT EDIT.

package tweets

import (
	context "context"

	predictor "github.com/kordape/tweety/internal/tweets/predictor"
	mock "github.com/stretchr/testify/mock"
)

// MockPredictor is an autogenerated mock type for the Predictor type
type MockPredictor struct {
	mock.Mock
}

// PredictFakeTweets provides a mock function with given fields: ctx, tweets
func (_m *MockPredictor) PredictFakeTweets(ctx context.Context, tweets []predictor.Tweet) (predictor.Response, error) {
	ret := _m.Called(ctx, tweets)

	var r0 predictor.Response
	if rf, ok := ret.Get(0).(func(context.Context, []predictor.Tweet) predictor.Response); ok {
		r0 = rf(ctx, tweets)
	} else {
		r0 = ret.Get(0).(predictor.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []predictor.Tweet) error); ok {
		r1 = rf(ctx, tweets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockPredictor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPredictor creates a new instance of MockPredictor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPredictor(t mockConstructorTestingTNewMockPredictor) *MockPredictor {
	mock := &MockPredictor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}