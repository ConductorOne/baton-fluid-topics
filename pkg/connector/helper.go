package connector

import (
	"errors"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/crypto"
)

const permissionName = "assigned"

// getRoleDescription looks for the role description by its name.
func getRoleDescription(roleName string) string {
	if desc, ok := roles[roleName]; ok {
		return desc
	}
	if desc, ok := adminRoles[roleName]; ok {
		return desc
	}
	return ""
}

// parseEntitlementId is responsible for cutting the id of the entitlement to know its name and type.
func parseEntitlementId(entitlementId string) (string, string, error) {
	parts := strings.Split(entitlementId, ":")
	if len(parts) != 4 {
		return "", "", fmt.Errorf("unexpected entitlement id format: %q", entitlementId)
	}
	return parts[1], parts[2], nil
}

// generateCredentials if the credential option is "Random Password", it returns a randomly generated password.
func generateCredentials(credentialOptions *v2.CredentialOptions) (string, error) {
	if credentialOptions.GetRandomPassword() == nil {
		return "", errors.New("unsupported credential option")
	}

	return crypto.GenerateRandomPassword(
		&v2.CredentialOptions_RandomPassword{
			Length: min(13, credentialOptions.GetRandomPassword().GetLength()),
		},
	)
}
