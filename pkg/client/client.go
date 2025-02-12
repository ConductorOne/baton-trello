package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/ratelimit"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	domain = "https://api.trello.com/1"

	getBoardById                 = "/boards/%s"
	getBoardsByOrganization      = "/organizations/%s/boards"
	getMemberById                = "/members/%s"
	getMembershipsByBoard        = "/boards/%s/memberships"
	getMembershipsByOrganization = "/organizations/%s/memberships"
	getOrganizationById          = "/organizations/%s"
	getUsersByOrganization       = "/organizations/%s/members"
)

type TrelloClient struct {
	ApiToken        string
	ApiKey          string
	BaseDomain      string
	OrganizationIDs []string
	wrapper         *uhttp.BaseHttpClient
}

func New(ctx context.Context, trelloClient *TrelloClient) (*TrelloClient, error) {
	var (
		clientKey       = trelloClient.ApiKey
		clientToken     = trelloClient.ApiToken
		clientDomain    = trelloClient.BaseDomain
		organizationIDs = trelloClient.OrganizationIDs
	)

	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, err
	}

	client := TrelloClient{
		wrapper:         cli,
		ApiKey:          clientKey,
		ApiToken:        clientToken,
		BaseDomain:      clientDomain,
		OrganizationIDs: organizationIDs,
	}

	return &client, nil
}

func NewClient(apiKey, apiToken string, organizationIDs []string, httpClient ...*uhttp.BaseHttpClient) *TrelloClient {
	var wrapper = &uhttp.BaseHttpClient{}
	if httpClient != nil || len(httpClient) != 0 {
		wrapper = httpClient[0]
	}
	return &TrelloClient{
		wrapper:         wrapper,
		BaseDomain:      domain,
		ApiKey:          apiKey,
		ApiToken:        apiToken,
		OrganizationIDs: organizationIDs,
	}
}

func (c *TrelloClient) ListUsers(ctx context.Context) ([]User, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res []User
	var annotation annotations.Annotations

	for _, id := range c.OrganizationIDs {
		queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getUsersByOrganization, id))
		if err != nil {
			l.Error(fmt.Sprintf("Error creating url: %s", err))
			return nil, nil, err
		}

		annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
		if err != nil {
			l.Error(fmt.Sprintf("Error getting resources: %s", err))
			return nil, nil, err
		}
	}

	return res, annotation, nil
}

func (c *TrelloClient) ListOrganizations(ctx context.Context) ([]Organization, annotations.Annotations, error) {
	var res []Organization
	annotation := annotations.Annotations{}

	for _, id := range c.OrganizationIDs {
		organizationDetail, incomingAnnotation, err := c.GetOrganizationDetail(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		if organizationDetail != nil {
			res = append(res, *organizationDetail)
			annotation = incomingAnnotation
		} else {
			return nil, nil, err
		}
	}

	return res, annotation, nil
}

func (c *TrelloClient) ListBoards(ctx context.Context) ([]Board, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var resources []Board
	var annotation annotations.Annotations

	for _, id := range c.OrganizationIDs {
		var res []Board
		queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getBoardsByOrganization, id))
		if err != nil {
			l.Error(fmt.Sprintf("Error creating url: %s", err))
			return nil, nil, err
		}

		annotation, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
		if err == nil {
			resources = append(resources, res...)
		} else {
			l.Error(fmt.Sprintf("Error getting resources: %s", err))
			return nil, nil, err
		}
	}

	return resources, annotation, nil
}

func (c *TrelloClient) GetBoardDetails(ctx context.Context, boardID string) (*Board, annotations.Annotations, error) {
	queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getBoardById, boardID))
	if err != nil {
		return nil, nil, err
	}
	var res *Board
	_, annotation, err := c.doRequest(ctx, http.MethodGet, queryUrl, &res)
	if err != nil {
		return nil, nil, err
	}

	return res, annotation, nil
}

func (c *TrelloClient) ListMembershipsByBoard(ctx context.Context, boardID string) ([]User, error) {
	queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getMembershipsByBoard, boardID))
	if err != nil {
		return nil, err
	}

	return c.listMembershipsByResource(ctx, queryUrl)
}

func (c *TrelloClient) GetOrganizationDetail(ctx context.Context, organizationID string) (*Organization, annotations.Annotations, error) {
	queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getOrganizationById, organizationID))
	if err != nil {
		return nil, nil, err
	}
	var res *Organization
	_, annotation, err := c.doRequest(ctx, http.MethodGet, queryUrl, &res)
	if err != nil {
		return nil, nil, err
	}

	return res, annotation, nil
}

func (c *TrelloClient) ListMembershipsByOrg(ctx context.Context, resourceID string) ([]User, error) {
	queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getMembershipsByOrganization, resourceID))
	if err != nil {
		return nil, err
	}

	return c.listMembershipsByResource(ctx, queryUrl)
}

func (c *TrelloClient) listMembershipsByResource(ctx context.Context, queryUrl string) ([]User, error) {
	var res []User
	var resources []User

	_, err := c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		return nil, err
	}

	for _, resource := range res {
		memberDetail, _, err := c.GetMemberDetails(ctx, resource.MemberID)
		if err != nil {
			return nil, err
		}
		memberDetail.MemberType = resource.MemberType

		resources = append(resources, *memberDetail)
	}

	return resources, nil
}

func (c *TrelloClient) GetMemberDetails(ctx context.Context, memberID string) (*User, annotations.Annotations, error) {
	queryUrl, err := url.JoinPath(c.BaseDomain, fmt.Sprintf(getMemberById, memberID))
	if err != nil {
		return nil, nil, err
	}
	var res *User
	_, annotation, err := c.doRequest(ctx, http.MethodGet, queryUrl, &res)
	if err != nil {
		return nil, nil, err
	}

	return res, annotation, nil
}

func (c *TrelloClient) getResourcesFromAPI(
	ctx context.Context,
	urlAddress string,
	res any,
) (annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res)

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *TrelloClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
) (http.Header, annotations.Annotations, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(authorizeEndpointUrl(c, endpointUrl))

	if err != nil {
		return nil, nil, err
	}

	req, err := c.wrapper.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
	)

	if err != nil {
		return nil, nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.wrapper.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
	case http.MethodDelete:
		resp, err = c.wrapper.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	annotation := annotations.Annotations{}
	if resp != nil {
		if desc, err := ratelimit.ExtractRateLimitData(resp.StatusCode, &resp.Header); err == nil {
			annotation.WithRateLimiting(desc)
		} else {
			return nil, annotation, err
		}

		return resp.Header, annotation, nil
	}

	return nil, nil, err
}

func authorizeEndpointUrl(c *TrelloClient, endpointUrl string) string {
	return endpointUrl + "?key=" + c.ApiKey + "&token=" + c.ApiToken
}
