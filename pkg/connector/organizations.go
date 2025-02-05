package connector

import (
	"context"
	"fmt"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
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

func (o *organizationBuilder) List(ctx context.Context, _ *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	organizations, annotation, err := o.client.ListOrganizations(ctx)

	if err != nil {
		return nil, "", nil, err
	}

	for _, organization := range *organizations {
		orgCopy := organization
		orgResource, err := parseIntoOrganizationResource(ctx, &orgCopy, nil)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, orgResource)
	}

	return resources, "", annotation, nil
}

func parseIntoOrganizationResource(_ context.Context, organization *client.Organization, parentResourceID *v2.ResourceId) (*v2.Resource, error) {

	profile := map[string]interface{}{
		"organization_id": organization.ID,
		"display_name":    organization.DisplayName,
	}

	groupTraits := []resource.GroupTraitOption{
		resource.WithGroupProfile(profile),
	}

	displayName := organization.DisplayName

	ret, err := resource.NewGroupResource(
		displayName,
		organizationResourceType,
		organization.ID,
		groupTraits,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *organizationBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement
	for _, memberType := range memberTypes {
		assigmentOptions := []entitlement.EntitlementOption{
			entitlement.WithGrantableTo(userResourceType),
			entitlement.WithDescription(fmt.Sprintf("Member type %s for organization %s in Trello", memberType, resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s Organization %s", resource.DisplayName, memberType)),
		}

		entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, memberType, assigmentOptions...))
	}

	return entitlements, "", nil, nil
}

func (o *organizationBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant

	var organizationID = resource.Id.Resource

	memberships, err := o.client.ListMembershipsByOrg(ctx, organizationID)

	if err != nil {
		return nil, "", nil, err
	}

	for _, membership := range *memberships {
		userResource, _ := parseIntoUserResource(ctx, &membership, resource.Id)
		membershipGrant := grant.NewGrant(resource, membership.MemberType, userResource, grant.WithAnnotation(&v2.V1Identifier{
			Id: fmt.Sprintf("org-grant:%s:%s:%s", resource.Id.Resource, membership.MemberID, membership.MemberType),
		}))
		grants = append(grants, membershipGrant)
	}

	return grants, "", nil, nil
}

func newOrganizationBuilder(c *client.TrelloClient) *organizationBuilder {
	return &organizationBuilder{
		resourceType: organizationResourceType,
		client:       c,
	}
}
