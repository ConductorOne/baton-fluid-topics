package client

import "time"

type UserList struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}

type UserUsage struct {
	ID                         string    `json:"id"`
	Active                     bool      `json:"active"`
	PrivacyPolicyAgreementDate time.Time `json:"privacyPolicyAgreementDate"`
	CreationDate               time.Time `json:"creationDate"`
	LastLoginDate              time.Time `json:"lastLoginDate"`
}

type UserRole struct {
	Id                  string   `json:"id"`
	ManualRoles         []string `json:"manualRoles"`
	AuthenticationRoles []string `json:"authenticationRoles"`
	DefaultRoles        []string `json:"defaultRoles"`
}

type RoleKey struct {
	Name string
	Type string
}

type UserGroup struct {
	Id                   string   `json:"id"`
	ManualGroups         []string `json:"manualGroups"`
	AuthenticationGroups []string `json:"authenticationGroups"`
}

type UserUsageResponse struct {
	User UserUsage `json:"user"`
}

type UserRoleResponse struct {
	User []UserRole `json:"user"`
}

type UserGroupResponse struct {
	User []UserGroup `json:"user"`
}

type RoleCache struct {
}
