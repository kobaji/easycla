// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
package mocks

import context "context"
import events "github.com/aws/aws-lambda-go/events"
import mock "github.com/stretchr/testify/mock"

// InterfacesInteractor is an autogenerated mock type for the InterfacesInteractor type
type InterfacesInteractor struct {
	mock.Mock
}

// Handler provides a mock function with given fields: _a0, _a1
func (_m *InterfacesInteractor) Handler(_a0 context.Context, _a1 events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 events.APIGatewayCustomAuthorizerResponse
	if rf, ok := ret.Get(0).(func(context.Context, events.APIGatewayCustomAuthorizerRequest) events.APIGatewayCustomAuthorizerResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(events.APIGatewayCustomAuthorizerResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, events.APIGatewayCustomAuthorizerRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
