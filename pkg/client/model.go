package client

import "time"

type UserList struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
}

type AuthenticationIdentifiers struct {
	Identifier string `json:"identifier"`
	Realm      string `json:"realm"`
}

type UserUsage struct {
	ID                        string                      `json:"id"`
	CreationDate              time.Time                   `json:"creationDate"`
	LastLoginDate             time.Time                   `json:"lastActivityDate"`
	AuthenticationIdentifiers []AuthenticationIdentifiers `json:"authenticationIdentifiers"`
}

type UserRole struct {
	Id                  string   `json:"id"`
	ManualRoles         []string `json:"manualRoles"`
	AuthenticationRoles []string `json:"authenticationRoles"`
	DefaultRoles        []string `json:"defaultRoles"`
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
