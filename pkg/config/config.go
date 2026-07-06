package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	ApiKeyField = field.StringField(
		"api-key",
		field.WithDisplayName("API Key"),
		field.WithPlaceholder("your-trello-api-key"),
		field.WithDescription("The API key for your Trello account"),
		// "API key is intended to be publicly accessible." (https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/#managing-your-api-key)
		field.WithIsSecret(false),
		field.WithRequired(true),
	)
	ApiTokenField = field.StringField(
		"api-token",
		field.WithDisplayName("API Token"),
		field.WithPlaceholder("your-trello-api-token"),
		field.WithDescription("The API token for your Trello account"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)
	OrganizationsField = field.StringSliceField(
		"organizations",
		field.WithDisplayName("Organizations"),
		field.WithDescription("Limit syncing to specific organizations by providing organization slugs."),
		field.WithRequired(true),
	)
	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Trello API URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

// Config is the configuration schema for the connector.
var Config = field.NewConfiguration(
	[]field.SchemaField{
		ApiKeyField,
		ApiTokenField,
		OrganizationsField,
		BaseURLField,
	},
	field.WithConnectorDisplayName("Trello"),
	field.WithIconUrl("/static/app-icons/trello.svg"),
	field.WithHelpUrl("/docs/baton/trello"),
)

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(c *Trello) error {
	return nil
}
