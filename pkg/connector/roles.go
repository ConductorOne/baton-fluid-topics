package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.FluidTopicsClient
}

var Roles = []string{"PERSONAL_BOOK_USER", "PERSONAL_BOOK_SHARE_USER", "HTML_EXPORT_USER", "PDF_EXPORT_USER",
	"COLLECTION_USER", "PRINT_USER", "OFFLINE_USER", "SAVED_SEARCH_USER", "BETA_USER", "DEBUG_USER",
	"ANALYTICS_USER", "FEEDBACK_USER", "RATING_USER", "GENERATIVE_AI_USER"}

var AdminRoles = []string{"ADMIN", "KHUB_ADMIN", "CONTENT_PUBLISHER", "USERS_ADMIN", "PORTAL_ADMIN"}

type TypedRole struct {
	Name string
	Type string // "manual", "authentication", "default"
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	for _, role := range Roles {
		for _, roleType := range []string{"manual", "authentication", "default"} {
			typedRole := TypedRole{Name: role, Type: roleType}
			roleResource, err := parseIntoTypedRoleResource(typedRole)
			if err != nil {
				return nil, "", nil, err
			}
			resources = append(resources, roleResource)
		}
	}

	for _, role := range AdminRoles {
		for _, roleType := range []string{"manual", "authentication"} {
			typedRole := TypedRole{Name: role, Type: roleType}
			roleResource, err := parseIntoTypedRoleResource(typedRole)
			if err != nil {
				return nil, "", nil, err
			}
			resources = append(resources, roleResource)
		}
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	permissionName := "assigned"

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(resource.DisplayName),
		entitlement.WithDisplayName(resource.DisplayName),
	}

	rv = append(rv, entitlement.NewPermissionEntitlement(resource, permissionName, assigmentOptions...))

	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoTypedRoleResource(role TypedRole) (*v2.Resource, error) {
	resourceID := fmt.Sprintf("%s:%s", role.Type, role.Name)
	displayName := fmt.Sprintf("%s:%s", role.Type, role.Name)

	ret, err := resource.NewResource(
		resourceID,
		roleResourceType,
		displayName,
	)
	if err != nil {
		return nil, err
	}

	ret.DisplayName = displayName

	return ret, nil
}

func newRoleBuilder(c *client.FluidTopicsClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}
