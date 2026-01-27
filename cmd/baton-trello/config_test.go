package main

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/test"
	cfg "github.com/conductorone/baton-trello/pkg/config"
)

func TestConfigs(t *testing.T) {
	configurationSchema := field.NewConfiguration(
		cfg.ConfigurationFields,
		field.WithConstraints(cfg.FieldRelationships...),
	)

	test.ExerciseTestCases(t, configurationSchema, ValidateConfig, []test.TestCase{})
}
