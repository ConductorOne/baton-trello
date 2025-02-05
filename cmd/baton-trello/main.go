package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-trello/pkg/client"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/conductorone/baton-sdk/pkg/types"
	connectorSchema "github.com/conductorone/baton-trello/pkg/connector"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-trello",
		getConnector,
		field.Configuration{
			Fields: ConfigurationFields,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)
	trelloClient := client.NewClient()

	if err := ValidateConfig(v); err != nil {
		return nil, err
	}

	apiKey := v.GetString(apiKeyField.FieldName)
	apiToken := v.GetString(apiTokenField.FieldName)
	orgs := v.GetStringSlice(organizations.FieldName)

	trelloClient = trelloClient.WithApiKey(apiKey).WithBearerToken(apiToken).WithOrganizationIDs(orgs)

	connectorBuilder, err := connectorSchema.New(ctx, trelloClient)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	opts := make([]connectorbuilder.Opt, 0)

	connector, err := connectorbuilder.NewConnector(ctx, connectorBuilder, opts...)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}
	return connector, nil
}
