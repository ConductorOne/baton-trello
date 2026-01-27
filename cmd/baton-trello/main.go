package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	cfg "github.com/conductorone/baton-trello/pkg/config"
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
		cfg.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connectorSchema.Connector{}),
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

	apiKey := v.GetString(cfg.ApiKeyField.FieldName)
	apiToken := v.GetString(cfg.ApiTokenField.FieldName)
	orgs := v.GetStringSlice(cfg.OrganizationsField.FieldName)

	if err := ValidateConfig(v); err != nil {
		return nil, err
	}

	connectorBuilder, err := connectorSchema.New(ctx, apiKey, apiToken, orgs)
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
