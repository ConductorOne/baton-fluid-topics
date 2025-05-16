package client

import (
	"context"

	"github.com/conductorone/baton-sdk/pkg/annotations"
)

type FluidTopicsClientInterface interface {
	ListUsers(ctx context.Context) ([]User, string, annotations.Annotations, error)
	GetUserDetails(ctx context.Context, userID string) (User, annotations.Annotations, error)
	GetAuthenticationInfo(ctx context.Context) (AuthenticationInfo, annotations.Annotations, error)
	UpdateUserManualRoles(ctx context.Context, userID string, manualRoles []string) (annotations.Annotations, error)
	CreateUser(ctx context.Context, newUser NewUserInfo) (annotations.Annotations, error)
	GetRolesByUserID(ctx context.Context, userID string) (UserRoles, annotations.Annotations, error)
}
