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

type boardBuilder struct {
	resourceType *v2.ResourceType
	client       *client.TrelloClient
}

func (o *boardBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return boardResourceType
}

func (o *boardBuilder) List(ctx context.Context, _ *v2.ResourceId, _ resource.SyncOpAttrs) ([]*v2.Resource, *resource.SyncOpResults, error) {
	var resources []*v2.Resource

	// Note: Trello API doesn't support pagination for boards by organization queries.
	boards, annotation, err := o.client.ListBoards(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, board := range boards {
		boardCopy := board
		parentResourceId, err := resource.NewResourceID(organizationResourceType, board.IdOrganization)
		if err != nil {
			return nil, nil, err
		}
		boardResource, err := parseIntoBoardResource(ctx, &boardCopy, parentResourceId)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, boardResource)
	}

	return resources, &resource.SyncOpResults{Annotations: annotation}, nil
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

func (o *boardBuilder) Entitlements(ctx context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	var entitlements []*v2.Entitlement

	// Note: Trello API only support getting boards by ID so it is not a bulk operation. Pagination is not needed.
	board, _, err := o.client.GetBoardDetails(ctx, res.Id.Resource)
	if err != nil {
		return nil, nil, err
	}

	// Self join
	selfJoin := "self join disabled"
	if board.Preferences.SelfJoin {
		selfJoin = "self join enabled"
	}
	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Is %s for board %s in Trello", selfJoin, res.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", res.DisplayName, selfJoin)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, selfJoin, assigmentOptions...))

	// Voting
	voting := fmt.Sprintf("Voting %s", board.Preferences.Voting)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", voting, res.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", res.DisplayName, voting)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, voting, assigmentOptions...))

	// Comments
	comments := fmt.Sprintf("Comments %s", board.Preferences.Comments)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", comments, res.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", res.DisplayName, comments)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, comments, assigmentOptions...))

	// Invitations
	invitations := fmt.Sprintf("Invitations %s", board.Preferences.Invitations)
	assigmentOptions = []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("%s for board %s in Trello", invitations, res.DisplayName)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Board %s", res.DisplayName, invitations)),
	}
	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(res, invitations, assigmentOptions...))

	return entitlements, nil, nil
}

func (o *boardBuilder) Grants(ctx context.Context, res *v2.Resource, _ resource.SyncOpAttrs) ([]*v2.Grant, *resource.SyncOpResults, error) {
	var grants []*v2.Grant

	var boardID = res.Id.Resource

	// Note: Trello API only support getting boards by ID so it is not a bulk operation. Pagination is not needed.
	board, _, err := o.client.GetBoardDetails(ctx, boardID)
	if err != nil {
		return nil, nil, err
	}
	memberships, err := o.client.ListMembershipsByBoard(ctx, boardID)
	if err != nil {
		return nil, nil, err
	}

	for _, membership := range memberships {
		userResource, _ := parseIntoUserResource(ctx, &membership, res.Id)
		membershipType := membership.MemberType

		// Self join
		if board.Preferences.SelfJoin {
			selfJoin := "self join enabled"
			membershipGrant := grant.NewGrant(res, selfJoin, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", res.Id.Resource, membership.MemberID, selfJoin),
			}))
			grants = append(grants, membershipGrant)
		}

		// Voting
		if evaluateMembership(membershipType, board.Preferences.Voting) {
			voting := fmt.Sprintf("Voting %s", board.Preferences.Voting)
			membershipGrant := grant.NewGrant(res, voting, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", res.Id.Resource, membership.MemberID, voting),
			}))
			grants = append(grants, membershipGrant)
		}

		// Comments
		if evaluateMembership(membershipType, board.Preferences.Comments) {
			comments := fmt.Sprintf("Comments %s", board.Preferences.Comments)
			membershipGrant := grant.NewGrant(res, comments, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", res.Id.Resource, membership.MemberID, comments),
			}))
			grants = append(grants, membershipGrant)
		}

		// Invitations
		if evaluateMembership(membershipType, board.Preferences.Invitations) {
			invitations := fmt.Sprintf("Invitations %s", board.Preferences.Invitations)
			membershipGrant := grant.NewGrant(res, invitations, userResource, grant.WithAnnotation(&v2.V1Identifier{
				Id: fmt.Sprintf("board-grant:%s:%s:%s", res.Id.Resource, membership.MemberID, invitations),
			}))
			grants = append(grants, membershipGrant)
		}
	}

	return grants, nil, nil
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
