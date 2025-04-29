package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.FluidTopicsClient
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	users, _, _, err := o.client.ListUsers(ctx, pToken)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userID := user.Id
		userUsage, _, err := o.client.GetUserUsage(ctx, userID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error getting user Usage %s: %w", userID, err)
		}
		userUsageCopy := userUsage
		userCopy := user
		userResource, err := parseIntoUserResource(&userCopy, &userUsageCopy, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, userResource)
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func getRoleDescription(roleName string) string {
	for _, r := range Roles {
		if r.Name == roleName {
			return r.Description
		}
	}
	for _, r := range AdminRoles {
		if r.Name == roleName {
			return r.Description
		}
	}
	return ""
}

// The Grants function in the roles resource is performed in users for a better performance,
// since in this way for each user there is, the grants are directly assigned depending on which roles he has.
func (u *userBuilder) Grants(ctx context.Context, res *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant
	var userID = res.Id.Resource

	user, _, err := u.client.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, "", nil, err
	}

	rolesTypes := []struct {
		RoleType string
		RoleList []string
	}{
		{"manual", user.ManualRoles},
		{"authentication", user.AuthenticationRoles},
		{"default", user.DefaultRoles},
	}

	for _, roleTypeData := range rolesTypes {
		for _, roleName := range roleTypeData.RoleList {
			description := getRoleDescription(roleName)

			typedRole := Role{
				Name:        roleName,
				Description: description,
				Type:        roleTypeData.RoleType,
			}

			roleResource, err := parseIntoTypedRoleResource(typedRole)
			if err != nil {
				return nil, "", nil, err
			}

			roleGrant := grant.NewGrant(roleResource, "assigned", res, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("role-grant:%s:%s:%s", roleName, userID, roleTypeData.RoleType),
			}),
			)
			grants = append(grants, roleGrant)
		}
	}
	return grants, "", nil, nil
}

func parseIntoUserResource(user *client.UserList, userUsage *client.UserUsage, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	var realm string
	if len(userUsage.AuthenticationIdentifiers) > 0 {
		realm = userUsage.AuthenticationIdentifiers[0].Realm
	}

	profile := map[string]interface{}{
		"user_id":              user.Id,
		"user_name":            user.DisplayName,
		"email_id":             user.Email,
		"creation_date":        userUsage.CreationDate.Format(time.RFC3339),
		"authentication_realm": realm,
	}

	displayName := user.DisplayName

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(displayName),
		rs.WithEmail(user.Email, true),
	}

	if !userUsage.LastLoginDate.IsZero() {
		userTraits = append(userTraits, rs.WithLastLogin(userUsage.LastLoginDate))
	}

	ret, err := rs.NewUserResource(
		displayName,
		userResourceType,
		user.Id,
		userTraits,
	)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newUserBuilder(c *client.FluidTopicsClient) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}
