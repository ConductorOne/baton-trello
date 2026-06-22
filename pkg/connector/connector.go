package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	cfg "github.com/conductorone/baton-trello/pkg/config"
	"github.com/conductorone/baton-trello/pkg/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type Connector struct {
	client *client.TrelloClient
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client),
		newOrganizationBuilder(d.client),
		newBoardBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Trello Connector",
		Description: "Connector to sync users, organizations and boards from Trello",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(_ context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, apiKey string, apiToken string, orgs []string, baseURL string) (*Connector, error) {
	l := ctxzap.Extract(ctx)

	trelloClient, err := client.New(ctx, client.NewClient(apiKey, apiToken, orgs, baseURL))
	if err != nil {
		l.Error("error creating Trello client", zap.Error(err))
		return nil, err
	}

	return &Connector{
		client: trelloClient,
	}, nil
}

// NewLambdaConnector returns a new ConnectorBuilderV2 suitable for use with RunConnector.
func NewLambdaConnector(ctx context.Context, ac *cfg.Trello, _ *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	c, err := New(ctx, ac.ApiKey, ac.ApiToken, ac.Organizations, ac.BaseUrl)
	if err != nil {
		return nil, nil, err
	}
	return c, nil, nil
}
