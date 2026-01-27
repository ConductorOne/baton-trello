package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ApiKeyField = field.StringField(
		"api-key",
		field.WithDescription("The API key for your Trello account"),
		field.WithRequired(true),
	)
	ApiTokenField = field.StringField(
		"api-token",
		field.WithDescription("The API token for your Trello account"),
		field.WithRequired(true),
	)
	OrganizationsField = field.StringSliceField(
		"organizations",
		field.WithDescription("Limit syncing to specific organizations by providing organization slugs."),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{ApiKeyField, ApiTokenField, OrganizationsField}

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
