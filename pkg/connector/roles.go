package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.FluidTopicsClient
}

type Role struct {
	Name        string
	Description string
	Type        string // "manual", "authentication", "default"
}

var Roles = []Role{
	{"PERSONAL_BOOK_USER", "Can create personal books", ""},
	{"PERSONAL_BOOK_SHARE_USER", "Can create and share personal books", ""},
	{"HTML_EXPORT_USER", "Can create personal books and download to HTML", ""},
	{"PDF_EXPORT_USER", "Can create personal books and download to PDF", ""},
	{"COLLECTION_USER", "Can create collections", ""},
	{"PRINT_USER", "Can use the print feature in the Reader page", ""},
	{"OFFLINE_USER", "Can use offline features", ""},
	{"SAVED_SEARCH_USER", "Can save searches", ""},
	{"BETA_USER", "Can use beta features", ""},
	{"DEBUG_USER", "Can access debug tools", ""},
	{"ANALYTICS_USER", "Can see analytics", ""},
	{"FEEDBACK_USER", "Can send feedback", ""},
	{"RATING_USER", "Can rate content", ""},
	{"GENERATIVE_AI_USER", "Can use generative AI features", ""},
}

var AdminRoles = []Role{
	{"ADMIN", "Administrator with full access", ""},
	{"KHUB_ADMIN", "Can administer the Knowledge Hub and publish, modify, and delete content published through any source", ""},
	{"CONTENT_PUBLISHER", "Can publish, modify, and delete content", ""},
	{"USERS_ADMIN", "Can list and manage users", ""},
	{"PORTAL_ADMIN", "Can configure portal display", ""},
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	// Iterar sobre los roles normales
	for _, role := range Roles {
		for _, roleType := range []string{"manual", "authentication", "default"} {
			typedRole := Role{Name: role.Name, Description: role.Description, Type: roleType} // Asegurarse de asignar el tipo
			roleResource, err := parseIntoTypedRoleResource(typedRole)
			if err != nil {
				return nil, "", nil, err
			}
			resources = append(resources, roleResource)
		}
	}

	// Iterar sobre los roles admin
	for _, role := range AdminRoles {
		for _, roleType := range []string{"manual", "authentication"} {
			typedRole := Role{Name: role.Name, Description: role.Description, Type: roleType}
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
		entitlement.WithDescription(resource.Description),
		entitlement.WithDisplayName(resource.DisplayName),
	}

	rv = append(rv, entitlement.NewPermissionEntitlement(resource, permissionName, assigmentOptions...))

	return rv, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoTypedRoleResource(role Role) (*v2.Resource, error) {
	resourceID := fmt.Sprintf("%s:%s", role.Type, role.Name)
	displayName := fmt.Sprintf("%s:%s", role.Type, role.Name)

	ret, err := rs.NewResource(
		resourceID,
		roleResourceType,
		displayName,
		rs.WithDescription(role.Description),
	)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func newRoleBuilder(c *client.FluidTopicsClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}
