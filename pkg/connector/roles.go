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
	client       client.FluidTopicsClientInterface
}

var roles = map[string]string{
	"PERSONAL_BOOK_USER":       "Can create personal books",
	"PERSONAL_BOOK_SHARE_USER": "Can create and share personal books",
	"HTML_EXPORT_USER":         "Can create personal books and download to HTML",
	"PDF_EXPORT_USER":          "Can create personal books and download to PDF",
	"COLLECTION_USER":          "Can create collections",
	"PRINT_USER":               "Can use the print feature in the Reader page",
	"OFFLINE_USER":             "Can use offline features",
	"SAVED_SEARCH_USER":        "Can save searches",
	"BETA_USER":                "Can use beta features",
	"DEBUG_USER":               "Can access debug tools",
	"ANALYTICS_USER":           "Can see analytics",
	"FEEDBACK_USER":            "Can send feedback",
	"RATING_USER":              "Can rate content",
	"GENERATIVE_AI_USER":       "Can use generative AI features",
}

var adminRoles = map[string]string{
	"ADMIN":             "Administrator with full access",
	"KHUB_ADMIN":        "Can administer the Knowledge Hub and publish, modify, and delete content published through any source",
	"CONTENT_PUBLISHER": "Can publish, modify, and delete content",
	"USERS_ADMIN":       "Can list and manage users",
	"PORTAL_ADMIN":      "Can configure portal display",
}

const (
	manualRole         = "manual"
	authenticationRole = "authentication"
	defaultRole        = "default"
)

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType { return roleResourceType }

func (r *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	for name, desc := range roles {
		for _, roleType := range []string{manualRole, authenticationRole, defaultRole} {
			role := client.Role{
				Name:        name,
				Description: desc,
				Type:        roleType,
			}
			roleResource, err := parseIntoRoleResource(ctx, role)
			if err != nil {
				return nil, "", nil, err
			}
			resources = append(resources, roleResource)
		}
	}

	for name, desc := range adminRoles {
		for _, roleType := range []string{manualRole, authenticationRole} {
			role := client.Role{
				Name:        name,
				Description: desc,
				Type:        roleType,
			}
			roleResource, err := parseIntoRoleResource(ctx, role)
			if err != nil {
				return nil, "", nil, err
			}
			resources = append(resources, roleResource)
		}
	}

	return resources, "", nil, nil
}

func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(resource.Description),
		entitlement.WithDisplayName(resource.DisplayName),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, permissionName, assigmentOptions...))

	return entitlements, "", nil, nil
}

func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	if principal.Id.ResourceType != userResourceType.Id {
		return nil, fmt.Errorf("only users can be granted with role membership")
	}

	userID := principal.Id.Resource

	roleType, roleName, err := parseEntitlementId(entitlement.Id)
	if err != nil {
		return nil, err
	}

	if roleType != manualRole {
		return nil, fmt.Errorf("only manual roles can be granted")
	}

	userRoles, _, err := r.client.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, existingRole := range userRoles.ManualRoles {
		if existingRole == roleName {
			return annotations.New(&v2.GrantAlreadyExists{}), nil
		}
	}

	userRoles.ManualRoles = append(userRoles.ManualRoles, roleName)

	annotation, err := r.client.UpdateUserManualRoles(ctx, userID, userRoles.ManualRoles)
	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (r *roleBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	userID := grant.Principal.Id.Resource

	roleType, roleName, err := parseEntitlementId(grant.Entitlement.Id)
	if err != nil {
		return nil, err
	}

	if roleType != manualRole {
		return nil, fmt.Errorf("only manual roles can be revoked")
	}

	userRoles, _, err := r.client.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	updatedRoles := []string{}
	found := false
	for _, existingRole := range userRoles.ManualRoles {
		if existingRole == roleName {
			found = true
			continue
		}
		updatedRoles = append(updatedRoles, existingRole)
	}

	if !found {
		return annotations.New(&v2.GrantAlreadyRevoked{}), nil
	}

	userRoles.ManualRoles = updatedRoles

	annotation, err := r.client.UpdateUserManualRoles(ctx, userID, userRoles.ManualRoles)
	if err != nil {
		return nil, err
	}

	return annotation, nil
}

// The Grants function in the roles resource is performed in users for a better performance,
// since in this way for each user there is, the grants are directly assigned depending on which roles he has.
func (r *roleBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func parseIntoRoleResource(_ context.Context, role client.Role) (*v2.Resource, error) {
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

func newRoleBuilder(c client.FluidTopicsClientInterface) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}
