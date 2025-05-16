package client

import "time"

type User struct {
	Id                        string                      `json:"id"`
	DisplayName               string                      `json:"displayName"`
	Email                     string                      `json:"emailAddress"`
	CreationDate              time.Time                   `json:"creationDate"`
	LastLoginDate             time.Time                   `json:"lastActivityDate"`
	AuthenticationIdentifiers []AuthenticationIdentifiers `json:"authenticationIdentifiers"`
	Credentials               Credentials                 `json:"credentials"`
}

type AuthenticationIdentifiers struct {
	Identifier string `json:"identifier"`
	Realm      string `json:"realm"`
}

type AuthenticationInfo struct {
	Profile struct {
		Roles []string `json:"roles"`
	} `json:"profile"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRoles struct {
	Id                  string   `json:"id"`
	ManualRoles         []string `json:"manualRoles"`
	AuthenticationRoles []string `json:"authenticationRoles"`
	DefaultRoles        []string `json:"defaultRoles"`
}

type FluidTopicsAPIError struct {
	Timestamp  string `json:"timestamp"`
	Status     int    `json:"status"`
	ErrorText  string `json:"error"`
	MessageStr string `json:"message"`
	Path       string `json:"path"`
}

type NewUserInfo struct {
	Name                   string `json:"name"`
	EmailAddress           string `json:"emailAddress"`
	Password               string `json:"password"`
	PrivacyPolicyAgreement bool   `json:"privacyPolicyAgreement"`
}

type UserDataResponse struct {
	User User `json:"user"`
}

type Role struct {
	Name        string
	Description string
	Type        string
}
