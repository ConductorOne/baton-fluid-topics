package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-fluid-topics/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type Connector struct {
	client *client.FluidTopicsClient
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(d.client),
		newRoleBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Fluid Topics Connector",
		Description: "Connector to sync and manage users in Fluid Topics.",
		AccountCreationSchema: &v2.ConnectorAccountCreationSchema{
			FieldMap: map[string]*v2.ConnectorAccountCreationSchema_Field{
				"name": {
					DisplayName: "Name",
					Required:    true,
					Description: "The display name of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "name",
					Order:       1,
				},
				"emailAddress": {
					DisplayName: "Email Address",
					Required:    true,
					Description: "The email address of the user.",
					Field: &v2.ConnectorAccountCreationSchema_Field_StringField{
						StringField: &v2.ConnectorAccountCreationSchema_StringField{},
					},
					Placeholder: "user@mail.com",
					Order:       2,
				},
			},
		},
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	roles, annotation, err := d.client.GetAuthenticationInfo(ctx)
	if err != nil {
		return annotation, fmt.Errorf("could not fetch current user roles: %w", err)
	}

	for _, role := range roles.Profile.Roles {
		if role == "ADMIN" {
			return annotation, nil
		}
	}

	return annotation, fmt.Errorf("authentication user must have ADMIN role to use this connector")
}

// New returns a new instance of the connector.
func New(ctx context.Context, fluidTopicsBearerToken string, fluidTopicsDomain string) (*Connector, error) {
	l := ctxzap.Extract(ctx)

	fluidTopicClient, err := client.New(ctx, fluidTopicsBearerToken, fluidTopicsDomain)
	if err != nil {
		l.Error("error creating Fluid Topics client", zap.Error(err))
		return nil, err
	}

	return &Connector{
		client: fluidTopicClient,
	}, nil
}
