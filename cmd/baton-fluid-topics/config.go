package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	bearerTokenField = field.StringField(
		"bearer-token",
		field.WithDescription("Bearer Token for authentication"),
		field.WithRequired(true),
	)
	domainField = field.StringField(
		"domain",
		field.WithDescription("Base domain for API, e.g. https://example.fluidtopics.net"),
		field.WithRequired(true),
	)
	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{
		bearerTokenField,
		domainField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated. For example, a
	// username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(v *viper.Viper) error {
	return nil
}
