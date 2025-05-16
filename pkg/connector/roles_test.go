package connector

import (
	"context"
	"errors"
	"testing"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/stretchr/testify/require"
)

func TestRoleBuilder_GrantAndRevoke(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	entitlementID := "Role:manual:KHUB_ADMIN:assigned"
	roleName := "KHUB_ADMIN"

	principal := &v2.Resource{
		Id: &v2.ResourceId{
			Resource:     userID,
			ResourceType: userResourceType.Id,
		},
	}
	entitlement := &v2.Entitlement{Id: entitlementID}
	grant := &v2.Grant{Principal: principal, Entitlement: entitlement}

	t.Run("Grant role to user", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{ManualRoles: []string{}}, annotations.New(nil), nil).Once()
		mockClient.On("UpdateUserManualRoles", ctx, userID, []string{roleName}).
			Return(annotations.New(nil), nil).Once()

		annotationsTest, err := rb.Grant(ctx, principal, entitlement)
		require.NoError(t, err)
		require.Empty(t, annotationsTest)

		mockClient.AssertExpectations(t)
	})

	t.Run("Grant role that is already assigned", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{ManualRoles: []string{roleName}}, annotations.New(nil), nil).Once()

		annotationsTest, err := rb.Grant(ctx, principal, entitlement)
		require.NoError(t, err)
		require.IsType(t, annotations.New(&v2.GrantAlreadyExists{}), annotationsTest)

		mockClient.AssertExpectations(t)
	})

	t.Run("Revoke existing role", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{ManualRoles: []string{roleName}}, annotations.New(nil), nil).Once()
		mockClient.On("UpdateUserManualRoles", ctx, userID, []string{}).
			Return(annotations.New(nil), nil).Once()

		annotationsTest, err := rb.Revoke(ctx, grant)
		require.NoError(t, err)
		require.IsType(t, annotations.New(&v2.GrantAlreadyRevoked{}), annotationsTest)

		mockClient.AssertExpectations(t)
	})

	t.Run("Revoke non-assigned role", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{ManualRoles: []string{}}, annotations.New(nil), nil).Once()

		annotationsTest, err := rb.Revoke(ctx, grant)
		require.NoError(t, err)
		require.IsType(t, annotations.New(&v2.GrantAlreadyRevoked{}), annotationsTest)

		mockClient.AssertExpectations(t)
	})

	t.Run("Grant fails if GetRolesByUserID returns error", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{}, annotations.New(nil), errors.New("API failure")).Once()

		annotationsTest, err := rb.Grant(ctx, principal, entitlement)
		require.Error(t, err)
		require.Nil(t, annotationsTest)

		mockClient.AssertExpectations(t)
	})

	t.Run("Grant fails if UpdateUserManualRoles returns error", func(t *testing.T) {
		mockClient := &client.MockFluidTopicsClient{}
		rb := newRoleBuilder(mockClient)

		mockClient.On("GetRolesByUserID", ctx, userID).
			Return(client.UserRoles{ManualRoles: []string{}}, annotations.New(nil), nil).Once()
		mockClient.On("UpdateUserManualRoles", ctx, userID, []string{roleName}).
			Return(annotations.New(nil), errors.New("update failed")).Once()

		annotationsTest, err := rb.Grant(ctx, principal, entitlement)
		require.Error(t, err)
		require.Nil(t, annotationsTest)

		mockClient.AssertExpectations(t)
	})
}
