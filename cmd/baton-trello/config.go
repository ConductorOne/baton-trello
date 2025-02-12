package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	apiKeyField = field.StringField(
		"api-key",
		field.WithDescription("The API key for your Trello account"),
		field.WithRequired(true),
	)
	apiTokenField = field.StringField(
		"api-token",
		field.WithDescription("The API token for your Trello account"),
		field.WithRequired(true),
	)
	organizations = field.StringSliceField(
		"organizations",
		field.WithDescription("Limit syncing to specific organizations by providing organization slugs."),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run. Note: these fields can be marked as optional or
	// required.
	ConfigurationFields = []field.SchemaField{apiKeyField, apiTokenField, organizations}

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
