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

	getUsersByOrganization       = "/organizations/%s/members"
	getOrganizationById          = "/organizations/%s"
	getMembershipsByOrganization = "/organizations/%s/memberships"
	getMemberById                = "/members/%s"
)

type TrelloClient struct {
	apiToken        string
	apiKey          string
	baseDomain      string
	organizationIDs []string
	wrapper         *uhttp.BaseHttpClient
}

func New(ctx context.Context, trelloClient *TrelloClient) (*TrelloClient, error) {
	var (
		clientKey       = trelloClient.getApiKey()
		clientToken     = trelloClient.getApiToken()
		clientDomain    = trelloClient.getBaseDomain()
		organizationIDs = trelloClient.getOrganizationIDs()
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
		apiKey:          clientKey,
		apiToken:        clientToken,
		baseDomain:      clientDomain,
		organizationIDs: organizationIDs,
	}

	return &client, nil
}

func NewClient(httpClient ...*uhttp.BaseHttpClient) *TrelloClient {
	var wrapper = &uhttp.BaseHttpClient{}
	if httpClient != nil || len(httpClient) != 0 {
		wrapper = httpClient[0]
	}
	return &TrelloClient{
		wrapper:         wrapper,
		baseDomain:      domain,
		apiKey:          "",
		apiToken:        "",
		organizationIDs: []string{},
	}
}

func (c *TrelloClient) WithBearerToken(apiToken string) *TrelloClient {
	c.apiToken = apiToken
	return c
}

func (c *TrelloClient) WithApiKey(apiKey string) *TrelloClient {
	c.apiKey = apiKey
	return c
}

func (c *TrelloClient) WithOrganizationIDs(organizationIDs []string) *TrelloClient {
	c.organizationIDs = organizationIDs
	return c
}

func (c *TrelloClient) getApiKey() string {
	return c.apiKey
}

func (c *TrelloClient) getApiToken() string {
	return c.apiToken
}

func (c *TrelloClient) getOrganizationIDs() []string {
	return c.organizationIDs
}

func (c *TrelloClient) getBaseDomain() string {
	return c.baseDomain
}

func (c *TrelloClient) ListUsers(ctx context.Context) (*[]User, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	res := &[]User{}
	var annotation annotations.Annotations

	for _, id := range c.organizationIDs {
		queryUrl, err := url.JoinPath(c.baseDomain, fmt.Sprintf(getUsersByOrganization, id))
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

func (c *TrelloClient) ListOrganizations(ctx context.Context) (*[]Organization, annotations.Annotations, error) {
	res := &[]Organization{}
	annotation := annotations.Annotations{}

	for _, id := range c.organizationIDs {
		organizationDetail, incomingAnnotation, err := c.GetOrganizationDetail(ctx, id)
		if err != nil {
			return nil, nil, err
		}

		if organizationDetail != nil {
			*res = append(*res, *organizationDetail)
			annotation = incomingAnnotation
		} else {
			return nil, nil, err
		}
	}

	return res, annotation, nil
}

func (c *TrelloClient) GetOrganizationDetail(ctx context.Context, organizationID string) (*Organization, annotations.Annotations, error) {
	queryUrl, err := url.JoinPath(c.baseDomain, fmt.Sprintf(getOrganizationById, organizationID))
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

func (c *TrelloClient) ListMembershipsByOrg(ctx context.Context, resourceID string) (*[]User, error) {
	queryUrl, err := url.JoinPath(c.baseDomain, fmt.Sprintf(getMembershipsByOrganization, resourceID))
	if err != nil {
		return nil, err
	}

	var res *[]User
	resources := &[]User{}

	_, err = c.getResourcesFromAPI(ctx, queryUrl, &res)
	if err != nil {
		return nil, err
	}

	for _, resource := range *res {
		memberDetail, _, err := c.GetMemberDetails(ctx, resource.MemberID)
		if err != nil {
			return nil, err
		}
		memberDetail.MemberType = resource.MemberType

		*resources = append(*resources, *memberDetail)
	}

	return resources, nil
}

func (c *TrelloClient) GetMemberDetails(ctx context.Context, memberID string) (*User, annotations.Annotations, error) {
	queryUrl, err := url.JoinPath(c.baseDomain, fmt.Sprintf(getMemberById, memberID))
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
	return endpointUrl + "?key=" + c.apiKey + "&token=" + c.apiToken
}
