package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
)

const (
	// getUsersEP OAuth Scope required: scim.groups.readUsers. https://developer.ironcladapp.com/reference/retrieve-all-users.
	getUsers      = "/users"
	getUsageUsers = "/users/%s/dump"
)

// Test data.
var testUsers = []client.UserList{
	{
		Id:          "eca7c247-dd7e-4318-a7ba-6f67a3f0a4b3",
		DisplayName: "First User",
		Email:       "user1@mail.com",
	},
	{
		Id:          "e128b37d-2b92-4276-8810-a624d16fbbda",
		DisplayName: "Second User",
		Email:       "user2@mail.com",
	},
	{
		Id:          "4318-a7ba-6f67a3f0a4b3",
		DisplayName: "Third User",
		Email:       "user3@mail.com",
	},
	{
		Id:          "c247-dda-6f67",
		DisplayName: "Fourth User",
		Email:       "user4@mail.com",
	},
}

var testUserUsage = map[string]client.UserUsage{
	"eca7c247-dd7e-4318-a7ba-6f67a3f0a4b3": {
		ID:                         "eca7c247-dd7e-4318-a7ba-6f67a3f0a4b3",
		Active:                     true,
		PrivacyPolicyAgreementDate: time.Now().AddDate(-1, 0, 0),
		CreationDate:               time.Now().AddDate(-2, 0, 0),
		LastLoginDate:              time.Now().AddDate(0, -1, 0),
	},
	"e128b37d-2b92-4276-8810-a624d16fbbda": {
		ID:                         "e128b37d-2b92-4276-8810-a624d16fbbda",
		Active:                     true,
		PrivacyPolicyAgreementDate: time.Now().AddDate(-1, 0, 0),
		CreationDate:               time.Now().AddDate(-2, 0, 0),
		LastLoginDate:              time.Now().AddDate(0, -1, 0),
	},
	"4318-a7ba-6f67a3f0a4b3": {
		ID:                         "4318-a7ba-6f67a3f0a4b3",
		Active:                     true,
		PrivacyPolicyAgreementDate: time.Now().AddDate(-1, 0, 0),
		CreationDate:               time.Now().AddDate(-2, 0, 0),
		LastLoginDate:              time.Now().AddDate(0, -1, 0),
	},
	"c247-dda-6f67": {
		ID:                         "c247-dda-6f67",
		Active:                     true,
		PrivacyPolicyAgreementDate: time.Now().AddDate(-1, 0, 0),
		CreationDate:               time.Now().AddDate(-2, 0, 0),
		LastLoginDate:              time.Now().AddDate(0, -1, 0),
	},
}

var testUserRoles = map[string]client.UserRole{
	"eca7c247-dd7e-4318-a7ba-6f67a3f0a4b3": {
		Id:                  "eca7c247-dd7e-4318-a7ba-6f67a3f0a4b3",
		ManualRoles:         []string{"admin", "developer", "billing", "viewer"},
		AuthenticationRoles: []string{"sso", "viewer"},
		DefaultRoles:        []string{"reader", "guest", "viewer"},
	},
	"e128b37d-2b92-4276-8810-a624d16fbbda": {
		Id:                  "e128b37d-2b92-4276-8810-a624d16fbbda",
		ManualRoles:         []string{"editor"},
		AuthenticationRoles: []string{"ldap"},
		DefaultRoles:        []string{"viewer", "basic"},
	},
	"4318-a7ba-6f67a3f0a4b3": {
		Id:                  "4318-a7ba-6f67a3f0a4b3",
		ManualRoles:         []string{"contributor"},
		AuthenticationRoles: []string{"saml", "moderator"},
		DefaultRoles:        []string{"guest"},
	},
	"c247-dda-6f67": {
		Id:                  "c247-dda-6f67",
		ManualRoles:         []string{"moderator", "openid"},
		AuthenticationRoles: []string{"openid"},
		DefaultRoles:        []string{"basic"},
	},
}

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users" {
			getAllUsersHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/roles") {
			getUserRolesHandler(w, r)
		} else {
			getUserUsageHandler(w, r)
		}
	})

	srv := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("✅ Test server running at http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}

func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// sin verificación de token
	log.Printf("GET /users - returning %d users", len(testUsers))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testUsers)
}
func getUserUsageHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[3] != "dump" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	userID := parts[2]

	usage, ok := testUserUsage[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	log.Printf("GET /users/%s/dump - returning usage data", userID)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":                         usage.ID,
			"active":                     usage.Active,
			"privacyPolicyAgreementDate": usage.PrivacyPolicyAgreementDate.Format(time.RFC3339),
			"creationDate":               usage.CreationDate.Format(time.RFC3339),
			"lastLoginDate":              usage.LastLoginDate.Format(time.RFC3339),
		},
	})
}

func getUserRolesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[3] != "roles" {
		http.NotFound(w, r)
		return
	}
	id := parts[2]

	role, ok := testUserRoles[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json") // <--- esto es clave
	json.NewEncoder(w).Encode(role)
}

func verifyAccessToken(r *http.Request) error {
	return nil
}
