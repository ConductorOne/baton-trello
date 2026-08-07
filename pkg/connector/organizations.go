package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-trello/pkg/client"
)

type organizationBuilder struct {
	resourceType *v2.ResourceType
	client       *client.TrelloClient
}

var memberTypes = []string{"admin", "normal", "observer"}

func (o *organizationBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return organizationResourceType
}

func (o *organizationBuilder) List(ctx context.Context, _ *v2.ResourceId, _ resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	var resources []*v2.Resource

	// Note: Trello API only support getting organizations by ID so it is not a bulk operation. Pagination is not needed.
	organizations, annotation, err := o.client.ListOrganizations(ctx)

	if err != nil {
		return nil, nil, err
	}

	for _, organization := range organizations {
		orgCopy := organization
		orgResource, err := parseIntoOrganizationResource(ctx, &orgCopy, nil)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, orgResource)
	}

	return resources, &resource.SyncOpResults{Annotations: annotation}, nil
}

func parseIntoOrganizationResource(_ context.Context, organization *client.Organization, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"organization_id": organization.ID,
		"display_name":    organization.DisplayName,
	}

	groupTraits := []resource.GroupTraitOption{}

	displayName := organization.DisplayName

	ret, err := resource.NewGroupResource(
		displayName,
		organizationResourceType,
		organization.ID,
		groupTraits,
		resource.WithResourceProfile(profile),
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *organizationBuilder) Entitlements(_ context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	var entitlements []*v2.Entitlement
	for _, memberType := range memberTypes {
		assigmentOptions := []entitlement.EntitlementOption{
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Member type %s for organization %s in Trello", memberType, res.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s Organization %s", res.DisplayName, memberType)),
		}

		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, memberType, assigmentOptions...))
	}

	return entitlements, nil, nil
}

func (o *organizationBuilder) Grants(ctx context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	var grants []*v2.Grant

	var organizationID = res.Id.Resource

	memberships, err := o.client.ListMembershipsByOrg(ctx, organizationID)
	if err != nil {
		return nil, nil, err
	}

	for _, membership := range memberships {
		userResource, _ := parseIntoUserResource(ctx, &membership, res.Id)
		membershipGrant := grant.NewGrant(res, membership.MemberType, userResource, grant.WithAnnotation(&v2.V1Identifier{
			Id: fmt.Sprintf("org-grant:%s:%s:%s", res.Id.Resource, membership.MemberID, membership.MemberType),
		}))
		grants = append(grants, membershipGrant)
	}

	return grants, nil, nil
}

func newOrganizationBuilder(c *client.TrelloClient) *organizationBuilder {
	return &organizationBuilder{
		resourceType: organizationResourceType,
		client:       c,
	}
}
