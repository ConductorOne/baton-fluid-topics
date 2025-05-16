package client

import (
	"context"

	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/stretchr/testify/mock"
)

type MockFluidTopicsClient struct {
	mock.Mock
}

func (m *MockFluidTopicsClient) ListUsers(ctx context.Context) ([]User, string, annotations.Annotations, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.String(1), args.Get(2).(annotations.Annotations), args.Error(3)
}

func (m *MockFluidTopicsClient) GetUserDetails(ctx context.Context, userID string) (User, annotations.Annotations, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Get(1).(annotations.Annotations), args.Error(2)
}

func (m *MockFluidTopicsClient) GetAuthenticationInfo(ctx context.Context) (AuthenticationInfo, annotations.Annotations, error) {
	args := m.Called(ctx)
	return args.Get(0).(AuthenticationInfo), args.Get(1).(annotations.Annotations), args.Error(2)
}

func (m *MockFluidTopicsClient) GetRolesByUserID(ctx context.Context, userID string) (UserRoles, annotations.Annotations, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(UserRoles), args.Get(1).(annotations.Annotations), args.Error(2)
}

func (m *MockFluidTopicsClient) UpdateUserManualRoles(ctx context.Context, userID string, manualRoles []string) (annotations.Annotations, error) {
	args := m.Called(ctx, userID, manualRoles)
	return args.Get(0).(annotations.Annotations), args.Error(1)
}

func (m *MockFluidTopicsClient) CreateUser(ctx context.Context, newUser NewUserInfo) (annotations.Annotations, error) {
	args := m.Called(ctx, newUser)
	return args.Get(0).(annotations.Annotations), args.Error(1)
}
