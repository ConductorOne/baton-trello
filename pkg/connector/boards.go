package connector

import (
	"context"
	"fmt"
	"sync"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-trello/pkg/client"
)

type boardBuilder struct {
	resourceType     *v2.ResourceType
	client           *client.TrelloClient
	memberships      []client.User
	membershipsMutex sync.RWMutex
}

func (o *boardBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return boardResourceType
}

func (o *boardBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	boards, annotation, err := o.client.ListBoards(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, board := range *boards {
		boardCopy := board
		parentResourceId, err := resource.NewResourceID(organizationResourceType, board.IdOrganization)
		if err != nil {
			return nil, "", nil, err
		}
		boardResource, err := parseIntoBoardResource(ctx, &boardCopy, parentResourceId)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, boardResource)
	}

	return resources, "", annotation, nil
}

func parseIntoBoardResource(_ context.Context, board *client.Board, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"board_id":         board.ID,
		"display_name":     board.Name,
		"description":      board.Description,
		"permission_level": board.Preferences.PermissionLevel,
		"hide_votes":       board.Preferences.HideVotes,
		"voting":           board.Preferences.Voting,
		"comments":         board.Preferences.Comments,
		"invitations":      board.Preferences.Invitations,
		"self_join":        board.Preferences.SelfJoin,
	}

	groupTraits := []resource.GroupTraitOption{
		resource.WithGroupProfile(profile),
	}

	displayName := board.Name

	ret, err := resource.NewGroupResource(
		displayName,
		boardResourceType,
		board.ID,
		groupTraits,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *boardBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement

	board, _, err := o.client.GetBoardDetails(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, err
	}

	// Self join
	selfJoin := "self join disabled"
	if board.Preferences.SelfJoin {
		selfJoin = "self join enabled"
	}
	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Is %s for board %s in Trello", selfJoin, resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", resource.DisplayName, selfJoin)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, selfJoin, assigmentOptions...))

	// Voting
	voting := fmt.Sprintf("Voting %s", board.Preferences.Voting)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", voting, resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", resource.DisplayName, voting)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, voting, assigmentOptions...))

	// Comments
	comments := fmt.Sprintf("Comments %s", board.Preferences.Comments)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", comments, resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", resource.DisplayName, comments)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, comments, assigmentOptions...))

	// Invitations
	invitations := fmt.Sprintf("Invitations %s", board.Preferences.Invitations)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", invitations, resource.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", resource.DisplayName, invitations)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, invitations, assigmentOptions...))

	return entitlements, "", nil, nil
}

func (o *boardBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant

	var boardID = resource.Id.Resource

	board, _, err := o.client.GetBoardDetails(ctx, boardID)
	if err != nil {
		return nil, "", nil, err
	}
	err = o.GetMemberships(ctx, boardID)
	if err != nil {
		return nil, "", nil, err
	}

	for _, membership := range o.memberships {
		userResource, _ := parseIntoUserResource(ctx, &membership, resource.Id)
		membershipType := membership.MemberType

		// Self join
		if board.Preferences.SelfJoin {
			selfJoin := "self join enabled"
			membershipGrant := grant.NewGrant(resource, selfJoin, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", resource.Id.Resource, membership.MemberID, selfJoin),
			}))
			grants = append(grants, membershipGrant)
		}

		// Voting
		if evaluateMembership(membershipType, board.Preferences.Voting) {
			voting := fmt.Sprintf("Voting %s", board.Preferences.Voting)
			membershipGrant := grant.NewGrant(resource, voting, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", resource.Id.Resource, membership.MemberID, voting),
			}))
			grants = append(grants, membershipGrant)
		}

		// Comments
		if evaluateMembership(membershipType, board.Preferences.Comments) {
			comments := fmt.Sprintf("Comments %s", board.Preferences.Comments)
			membershipGrant := grant.NewGrant(resource, comments, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", resource.Id.Resource, membership.MemberID, comments),
			}))
			grants = append(grants, membershipGrant)
		}

		// Invitations
		if evaluateMembership(membershipType, board.Preferences.Invitations) {
			invitations := fmt.Sprintf("Invitations %s", board.Preferences.Invitations)
			membershipGrant := grant.NewGrant(resource, invitations, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", resource.Id.Resource, membership.MemberID, invitations),
			}))
			grants = append(grants, membershipGrant)
		}
	}

	return grants, "", nil, nil
}

func newBoardBuilder(c *client.TrelloClient) *boardBuilder {
	return &boardBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}

func evaluateMembership(membershipType, permission string) bool {
	return (membershipType == "admin" && permission == "admins") || permission == "members"
}

func (o *boardBuilder) GetMemberships(ctx context.Context, boardID string) error {
	o.membershipsMutex.RLock()
	defer o.membershipsMutex.RUnlock()

	if o.memberships != nil || len(o.memberships) > 0 {
		return nil
	}

	memberships, err := o.client.ListMembershipsByBoard(ctx, boardID)

	if err != nil {
		return err
	}

	o.memberships = append(o.memberships, *memberships...)

	return nil
}
