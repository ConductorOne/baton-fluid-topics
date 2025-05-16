package connector

import (
	"context"
	"testing"
	"time"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserBuilder_WithMockClient(t *testing.T) {
	ctx := context.Background()
	mockClient := &client.MockFluidTopicsClient{}
	ub := newUserBuilder(mockClient)

	testUser := client.User{
		Id:           "a061ccd9-3b8d-4f73-8d21-d045b3680a9d",
		DisplayName:  "Test User",
		Email:        "test@x.com",
		CreationDate: time.Now(),
		AuthenticationIdentifiers: []client.AuthenticationIdentifiers{
			{
				Identifier: "test@x.com",
				Realm:      "exampleRealm",
			},
		},
		Credentials: client.Credentials{
			Login:    "test@x.com",
			Password: "secretexamplepass",
		},
	}

	mockClient.On("ListUsers", mock.Anything).Return([]client.User{testUser}, "", annotations.Annotations{}, nil)
	mockClient.On("GetUserDetails", ctx, "a061ccd9-3b8d-4f73-8d21-d045b3680a9d").Return(testUser, annotations.Annotations{}, nil)

	t.Run("List should fetch users and details", func(t *testing.T) {
		users, _, _, err := ub.List(ctx, nil, nil)

		require.NoError(t, err)
		require.Len(t, users, 1)
		require.Equal(t, "Test User", users[0].DisplayName)

		mockClient.AssertExpectations(t)
	})

	t.Run("Grants should return role grants", func(t *testing.T) {
		mockClient.On("GetRolesByUserID", ctx, "u123").Return(client.UserRoles{
			ManualRoles:         []string{"COLLECTION_USER"},
			AuthenticationRoles: []string{"PRINT_USER", "ADMIN"},
			DefaultRoles:        []string{"PRINT_USER"},
		}, annotations.Annotations{}, nil)

		resource := &v2.Resource{Id: &v2.ResourceId{Resource: "u123"}}
		grants, _, _, err := ub.Grants(ctx, resource, nil)
		require.NoError(t, err)
		require.Len(t, grants, 4)

		var actualEntitlementIDs []string
		for _, g := range grants {
			actualEntitlementIDs = append(actualEntitlementIDs, g.Entitlement.Id)
		}

		expectedEntitlementIDs := []string{
			"Role:manual:COLLECTION_USER:assigned",
			"Role:authentication:PRINT_USER:assigned",
			"Role:authentication:ADMIN:assigned",
			"Role:default:PRINT_USER:assigned",
		}

		require.ElementsMatch(t, expectedEntitlementIDs, actualEntitlementIDs)
	})
}
