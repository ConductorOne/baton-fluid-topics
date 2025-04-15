package connector

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.FluidTopicsClient
	roleCache    []*v2.Resource
	RoleUserMap  map[client.RoleKey][]string
	roleMutex    sync.RWMutex
	cacheinit    bool
}

var Roles = []string{"ADMIN", "KHUB_ADMIN", "CONTENT_PUBLISHER", "USERS_ADMIN", "PORTAL_ADMIN",
	"PERSONAL_BOOK_USER", "PERSONAL_BOOK_SHARE_USER", "HTML_EXPORT_USER", "PDF_EXPORT_USER",
	"COLLECTION_USER", "PRINT_USER", "OFFLINE_USER", "SAVED_SEARCH_USER", "BETA_USER", "DEBUG_USER",
	"ANALYTICS_USER", "FEEDBACK_USER", "RATING_USER", "GENERATIVE_AI_USER"}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	roleMap := make(map[client.RoleKey][]string)
	var resources []*v2.Resource
	users, _, _, err := r.client.ListUsers(ctx, pToken)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userID := user.Id
		userRoles, _, err := r.client.GetRolesByUserID(ctx, userID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error getting roles for user %s: %w", userID, err)
		}
		for _, role := range userRoles.ManualRoles {
			key := client.RoleKey{Name: role, Type: "manualRoles"}
			roleMap[key] = append(roleMap[key], userID)
		}
		for _, role := range userRoles.AuthenticationRoles {
			key := client.RoleKey{Name: role, Type: "authenticationRoles"}
			roleMap[key] = append(roleMap[key], userID)
		}
		for _, role := range userRoles.DefaultRoles {
			key := client.RoleKey{Name: role, Type: "defaultRoles"}
			roleMap[key] = append(roleMap[key], userID)
		}
	}

	for key := range roleMap {
		res, err := parseIntoRoleResource(key)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error parsing role %s: %w", key, err)
		}
		resources = append(resources, res)
	}

	r.roleCache = resources
	r.RoleUserMap = roleMap // <-- acá guardás la relación de usuarios por rol
	r.cacheinit = true

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	permissionName := fmt.Sprintf("assigned")

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
	const permissionName = "assigned"
	var rv []*v2.Grant

	_, err := r.GetCachedRoles(ctx, nil, nil)
	if err != nil {
		return nil, "", nil, err
	}

	// Suponiendo que resource.Id.Resource contiene el nombre del rol
	// y resource.DisplayName tiene el formato "roleName : roleType"
	parts := strings.SplitN(resource.DisplayName, ":", 2)
	if len(parts) != 2 {
		return nil, "", nil, fmt.Errorf("invalid display name format: %s", resource.DisplayName)
	}
	roleName := parts[0]
	roleType := parts[1]

	roleKey := client.RoleKey{Name: roleName, Type: roleType}

	userIDs, ok := r.RoleUserMap[roleKey]
	if !ok {
		// Nadie tiene este rol
		return nil, "", nil, nil
	}

	for _, userID := range userIDs {
		userRes := &v2.ResourceId{
			ResourceType: userResourceType.Id,
			Resource:     userID,
		}
		gr := grant.NewGrant(resource, permissionName, userRes)
		rv = append(rv, gr)
	}

	return rv, "", nil, nil
}

func (r *roleBuilder) GetCachedRoles(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, error) {
	// Intentamos leer cache
	r.roleMutex.RLock()
	if r.cacheinit {
		defer r.roleMutex.RUnlock()
		return r.roleCache, nil
	}
	r.roleMutex.RUnlock()

	// Necesitamos cargar la cache
	r.roleMutex.Lock()
	defer r.roleMutex.Unlock()

	// Doble verificación por seguridad
	if r.cacheinit {
		return r.roleCache, nil
	}

	roles, _, _, err := r.List(ctx, parentResourceID, pToken)
	if err != nil {
		return nil, err
	}

	r.roleCache = roles
	r.cacheinit = true

	return roles, nil
}

func parseIntoRoleResource(roleKey client.RoleKey) (*v2.Resource, error) {
	displayName := fmt.Sprintf("%s:%s", roleKey.Name, roleKey.Type)
	resourceID := fmt.Sprintf("%s:%s", roleKey.Name, roleKey.Type)

	ret, err := resource.NewResource(
		displayName,
		roleResourceType,
		resourceID,
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
