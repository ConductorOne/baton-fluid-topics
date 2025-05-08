package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.FluidTopicsClient
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (u *userBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	users, _, _, err := u.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userID := user.Id
		userCopy, _, err := u.client.GetUserDetails(ctx, userID)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error getting user details %s: %w", userID, err)
		}
		userResource, err := parseIntoUserResource(ctx, &userCopy)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, userResource)
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (u *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
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

			typedRole := client.Role{
				Name:        roleName,
				Description: description,
				Type:        roleTypeData.RoleType,
			}

			roleResource, err := parseIntoRoleResource(ctx, typedRole)
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

func (u *userBuilder) CreateAccountCapabilityDetails(_ context.Context) (*v2.CredentialDetailsAccountProvisioning, annotations.Annotations, error) {
	return &v2.CredentialDetailsAccountProvisioning{
		SupportedCredentialOptions: []v2.CapabilityDetailCredentialOption{
			v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_RANDOM_PASSWORD,
		},
		PreferredCredentialOption: v2.CapabilityDetailCredentialOption_CAPABILITY_DETAIL_CREDENTIAL_OPTION_RANDOM_PASSWORD,
	}, nil, nil
}

func (u *userBuilder) CreateAccount(
	ctx context.Context,
	accountInfo *v2.AccountInfo,
	credentialOptions *v2.CredentialOptions,
) (connectorbuilder.CreateAccountResponse, []*v2.PlaintextData, annotations.Annotations, error) {
	newUser, err := createNewUserInfo(accountInfo, credentialOptions)
	if err != nil {
		return nil, nil, annotations.Annotations{}, err
	}

	_, err = u.client.CreateUser(ctx, *newUser)
	if err != nil {
		return nil, nil, annotations.Annotations{}, err
	}

	userResource, err := parseIntoUserResource(
		ctx,
		&client.User{
			DisplayName: newUser.Name,
			Email:       newUser.EmailAddress,
			Credentials: client.Credentials{
				Login:    newUser.EmailAddress,
				Password: newUser.Password,
			},
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}
	caResponse := &v2.CreateAccountResponse_SuccessResult{
		Resource: userResource,
	}

	passResult := &v2.PlaintextData{
		Name:  "password",
		Bytes: []byte(newUser.Password),
	}
	return caResponse, []*v2.PlaintextData{passResult}, nil, nil
}

func createNewUserInfo(accountInfo *v2.AccountInfo, credentialOptions *v2.CredentialOptions) (*client.NewUserInfo, error) {
	pMap := accountInfo.Profile.AsMap()

	name, ok := pMap["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required")
	}

	email, ok := pMap["emailAddress"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email is required")
	}

	generatedPassword, err := generateCredentials(credentialOptions)
	if err != nil {
		return nil, err
	}

	newUser := &client.NewUserInfo{
		Name:                   name,
		Password:               generatedPassword,
		EmailAddress:           email,
		PrivacyPolicyAgreement: true,
	}

	return newUser, nil
}

func parseIntoUserResource(ctx context.Context, user *client.User) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	var realm string
	if len(user.AuthenticationIdentifiers) > 0 {
		realm = user.AuthenticationIdentifiers[0].Realm
	}

	profile := map[string]interface{}{
		"user_id":              user.Id,
		"user_name":            user.DisplayName,
		"email_id":             user.Email,
		"creation_date":        user.CreationDate.Format(time.RFC3339),
		"authentication_realm": realm,
	}

	displayName := user.DisplayName

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(displayName),
		rs.WithEmail(user.Email, true),
	}

	if !user.LastLoginDate.IsZero() {
		userTraits = append(userTraits, rs.WithLastLogin(user.LastLoginDate))
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
